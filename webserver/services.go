/*
* Copyright 2022-2024 Thorsten A. Knieling
*
* Licensed under the Apache License, Version 2.0 (the "License");
* you may not use this file except in compliance with the License.
* You may obtain a copy of the License at
*
*    http://www.apache.org/licenses/LICENSE-2.0
*
 */

package webserver

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/go-openapi/runtime/flagext"
	"github.com/go-openapi/swag"
	"github.com/tknie/clu/server"
	"github.com/tknie/log"
	"github.com/tknie/services"
	"golang.org/x/net/netutil"
)

// MaxHeaderSize maximum header size
var MaxHeaderSize = flagext.ByteSize(1)

// CleanupTimeout cleanup timeout
var CleanupTimeout = 10 * time.Second
var mainHandler http.Handler
var wg = new(sync.WaitGroup)
var servers = make([]*http.Server, 0)

var interrupt chan os.Signal
var interrupted = false
var once = new(sync.Once)

func init() {
	interrupt = make(chan os.Signal, 1)
	MaxHeaderSize.Set("1MiB")

}

// InitServices init services signal handler
func InitServices() {
	signalNotify(interrupt)
	go handleInterrupt(once)

}

// StartServices start services HTTP and HTTPS
func StartServices(mainHandler http.Handler) error {
	startSocket()
	err := startHTTP()
	if err != nil {
		fmt.Println("Error starting server", err)
		return err
	}
	err = startHTTPS()
	if err != nil {
		fmt.Println("Error starting server", err)
		return err
	}
	wg.Wait()
	return nil
}

// startSocket start local socket
func startSocket() {
	for _, s := range server.Viewer.Server.Service {
		if strings.ToLower(s.Type) == "socket" {
			domainSocket := new(http.Server)
			domainSocket.MaxHeaderBytes = int(MaxHeaderSize)
			domainSocket.Handler = mainHandler
			if int64(CleanupTimeout) > 0 {
				domainSocket.IdleTimeout = CleanupTimeout
			}
			home, err := os.UserHomeDir()
			if err != nil {
				log.Log.Fatal("User HOME directory evaluation error:", err)
			}
			socketPath := home + "/.cluapi/run/cluapi.sock"

			domSockListener, err := net.Listen("unix", string(socketPath))
			if err != nil {
				log.Log.Fatal(err)
			}
			servers = append(servers, domainSocket)
			go func(l net.Listener) {
				defer wg.Done()
				if err := domainSocket.Serve(l); err != nil && err != http.ErrServerClosed {
					log.Log.Fatal("Error starting server on socket path", err)
				}
				log.Log.Debugf("Stopped serving clutron at unix://%s", socketPath)
			}(domSockListener)
			wg.Add(1)
		}
	}
}

// startHTTP start HTTP service
func startHTTP() error {
	for _, s := range server.Viewer.Server.Service {
		if strings.ToLower(s.Type) == "http" {

			listener, err := net.Listen("tcp", net.JoinHostPort(s.Host, strconv.Itoa(s.Port)))
			if err != nil {
				return err
			}

			h, p, err := swag.SplitHostPort(listener.Addr().String())
			if err != nil {
				return err
			}
			s.Host = h
			s.Port = p
			httpServer := new(http.Server)
			httpServer.MaxHeaderBytes = int(s.MaxHeaderSize)
			httpServer.ReadTimeout = s.ReadTimeout
			httpServer.WriteTimeout = s.WriteTimeout
			httpServer.SetKeepAlivesEnabled(int64(s.KeepAlive) > 0)
			if s.ListenLimit > 0 {
				listener = netutil.LimitListener(listener, s.ListenLimit)
			}

			if int64(s.CleanupTimeout) > 0 {
				httpServer.IdleTimeout = s.CleanupTimeout
			}

			httpServer.Handler = mainHandler

			servers = append(servers, httpServer)
			wg.Add(1)
			log.Log.Debugf("Serving clutron at http://%s", listener.Addr())
			services.ServerMessage("Listen HTTP on address %s", listener.Addr())
			go func(l net.Listener) {
				defer wg.Done()
				if err := httpServer.Serve(l); err != nil && err != http.ErrServerClosed {
					log.Log.Fatal(err)
				}
				log.Log.Debugf("Stopped serving clutron at http://%s", l.Addr())
			}(listener)
		}
	}
	return nil
}

// startHTTPS start HTTPS service
func startHTTPS() error {
	for _, s := range server.Viewer.Server.Service {
		if strings.ToLower(s.Type) == "https" {
			tlsListener, err := net.Listen("tcp", net.JoinHostPort(s.Host, strconv.Itoa(s.Port)))
			if err != nil {
				return err
			}

			sh, sp, err := swag.SplitHostPort(tlsListener.Addr().String())
			if err != nil {
				return err
			}
			s.Host = sh
			s.Port = sp

			httpsServer := new(http.Server)
			httpsServer.MaxHeaderBytes = int(s.MaxHeaderSize)
			httpsServer.ReadTimeout = s.ReadTimeout
			httpsServer.WriteTimeout = s.WriteTimeout
			httpsServer.SetKeepAlivesEnabled(int64(s.KeepAlive) > 0)
			if s.ListenLimit > 0 {
				tlsListener = netutil.LimitListener(tlsListener, s.ListenLimit)
			}
			if int64(s.CleanupTimeout) > 0 {
				httpsServer.IdleTimeout = s.CleanupTimeout
			}
			httpsServer.Handler = mainHandler

			// Inspired by https://blog.bracebin.com/achieving-perfect-ssl-labs-score-with-go
			httpsServer.TLSConfig = &tls.Config{
				// Causes servers to use Go's default ciphersuite preferences,
				// which are tuned to avoid attacks. Does nothing on clients.
				PreferServerCipherSuites: true,
				// Only use curves which have assembly implementations
				// https://github.com/golang/go/tree/master/src/crypto/elliptic
				CurvePreferences: []tls.CurveID{tls.CurveP256},
				// Use modern tls mode https://wiki.mozilla.org/Security/Server_Side_TLS#Modern_compatibility
				NextProtos: []string{"h2", "http/1.1"},
				// https://www.owasp.org/index.php/Transport_Layer_Protection_Cheat_Sheet#Rule_-_Only_Support_Strong_Protocols
				MinVersion: tls.VersionTLS12,
				// These ciphersuites support Forward Secrecy: https://en.wikipedia.org/wiki/Forward_secrecy
				CipherSuites: []uint16{
					tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
					tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
					tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
					tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
					tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,
					tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,
				},
			}

			// build standard config from server options
			if s.TLSCertificate != "" && s.TLSCertificateKey != "" {
				/*httpsServer.TLSConfig.Certificates = make([]tls.Certificate, 1)
				httpsServer.TLSConfig.Certificates[0], err = tls.LoadX509KeyPair(string(s.TLSCertificate), string(s.TLSCertificateKey))
				if err != nil {
					return err
				}*/
				reloader, err := NewCertificateReloader(string(s.TLSCertificate), string(s.TLSCertificateKey))
				if err != nil {
					return err
				}
				httpsServer.TLSConfig.GetCertificate = reloader.GetCertificateFunc()
			}

			if s.TLSCACertificate != "" {
				// include specified CA certificate
				caCert, caCertErr := os.ReadFile(string(s.TLSCACertificate))
				if caCertErr != nil {
					return caCertErr
				}
				caCertPool := x509.NewCertPool()
				ok := caCertPool.AppendCertsFromPEM(caCert)
				if !ok {
					return fmt.Errorf("cannot parse CA certificate")
				}
				httpsServer.TLSConfig.ClientCAs = caCertPool
				httpsServer.TLSConfig.ClientAuth = tls.RequireAndVerifyClientCert
			}

			// call custom TLS configurator
			configureTLS(httpsServer.TLSConfig)

			if len(httpsServer.TLSConfig.Certificates) == 0 && httpsServer.TLSConfig.GetCertificate == nil {
				// after standard and custom config are passed, this ends up with no certificate
				if s.TLSCertificate == "" {
					if s.TLSCertificateKey == "" {
						log.Log.Fatal("the required flags `--tls-certificate` and `--tls-key` were not specified")
					}
					log.Log.Fatal("the required flag `--tls-certificate` was not specified")
				}
				if s.TLSCertificateKey == "" {
					log.Log.Fatal("the required flag `--tls-key` was not specified")
				}
				// this happens with a wrong custom TLS configurator
				log.Log.Fatal("no certificate was configured for TLS")
			}

			servers = append(servers, httpsServer)
			wg.Add(1)
			log.Log.Debugf("Serving clutron at https://%s", tlsListener.Addr())
			services.ServerMessage("Listen HTTPS on address %s", tlsListener.Addr())
			go func(l net.Listener) {
				defer wg.Done()
				if err := httpsServer.Serve(l); err != nil && err != http.ErrServerClosed {
					log.Log.Fatal(err)
				}
				log.Log.Debugf("Stopped serving clutron at https://%s", l.Addr())
			}(tls.NewListener(tlsListener, httpsServer.TLSConfig))
		}
	}
	return nil
}

// configureTLS The TLS configuration before HTTPS server starts.
func configureTLS(tlsConfig *tls.Config) {
	// Make all necessary changes to the TLS configuration here.
	tlsConfig.Certificates = make([]tls.Certificate, 1)
	for _, s := range server.Viewer.Server.Service {
		if strings.ToLower(s.Type) == "https" {
			if s.Certificate != "" && s.Key != "" {
				var err error
				certPath := os.ExpandEnv(s.Certificate)
				keyPath := os.ExpandEnv(s.Key)
				reloader, err := NewCertificateReloader(certPath, keyPath)
				if err != nil {
					services.ServerMessage("TLS configuration error: %v", err)
					os.Exit(1)
				}
				tlsConfig.GetCertificate = reloader.GetCertificateFunc()

			} else {
				services.ServerMessage("TLS default configuration used")
			}
			break
		}
	}
}

// Shutdown shutdown
func Shutdown() error {
	for _, s := range servers {
		s.Close()
	}
	return nil
}

func handleInterrupt(once *sync.Once) {
	once.Do(func() {
		for range interrupt {
			if interrupted {
				log.Log.Debugf("Server already shutting down")
				continue
			}
			interrupted = true
			log.Log.Debugf("Shutting down... ")
			if err := Shutdown(); err != nil {
				log.Log.Debugf("HTTP server Shutdown: %v", err)
			}
		}
	})
}

func signalNotify(interrupt chan<- os.Signal) {
	signal.Notify(interrupt, syscall.SIGINT, syscall.SIGTERM)
}

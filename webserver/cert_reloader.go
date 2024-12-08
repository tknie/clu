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
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/fsnotify/fsnotify"
	"github.com/tknie/log"
	"github.com/tknie/services"
)

type keypairReloader struct {
	certLock sync.RWMutex
	cert     *tls.Certificate
	certPath string
	keyPath  string
	done     chan bool
	watcher  *fsnotify.Watcher
}

// newCertificateReloader certificate reloader
func newCertificateReloader(certPath, keyPath string) (*keypairReloader, error) {
	result := &keypairReloader{
		certPath: certPath,
		keyPath:  keyPath,
		done:     make(chan bool),
	}
	cert, err := tls.LoadX509KeyPair(certPath, keyPath)
	if err != nil {
		return nil, err
	}
	result.cert = &cert
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGHUP)
		for range c {
			log.Log.Debugf("Received SIGHUP, reloading TLS certificate and key from %q and %q", certPath, keyPath)
			if err := result.checkReload(); err != nil {
				log.Log.Debugf("Keeping old TLS certificate because the new one could not be loaded: %v", err)
			}
		}
	}()
	watcher, err := fsnotify.NewWatcher()
	if err == nil {

		go func() {
			for {
				select {
				// watch for events
				case event := <-watcher.Events:
					services.ServerMessage("Noticed configuration changed in %s (%v)", event.Name, event.Op)
					if err := result.checkReload(); err != nil {
						log.Log.Debugf("Keeping old TLS certificate because the new one could not be loaded: %v", err)
					}

				case err := <-watcher.Errors:
					services.ServerMessage("Watcher ERROR received: %v", err)
				case <-result.done:
					watcher.Close()
					return
				}
			}
		}()
		// out of the box fsnotify can watch a single file, or a single directory
		if err = watcher.Add(certPath); err != nil {
			services.ServerMessage("ERROR add watcher %s: %v", certPath, err)
		} else {
			log.Log.Infof("Certificate watcher enabled for %s", certPath)
		}
		if err = watcher.Add(keyPath); err != nil {
			services.ServerMessage("ERROR add watcher %s: %v", keyPath, err)
		} else {
			log.Log.Infof("Key watcher enabled for %s", keyPath)
		}
		result.watcher = watcher
	} else {
		services.ServerMessage("ERROR creating watcher", err)
	}
	if err != nil {
		services.ServerMessage("Certificate error init watching: %v", err)
		return nil, err
	}
	services.ServerMessage("SSL/TLS Server certificate watcher enabled")
	return result, nil
}

func (reloader *keypairReloader) checkReload() error {
	newCert, err := tls.LoadX509KeyPair(reloader.certPath, reloader.keyPath)
	if err != nil {
		return err
	}
	reloader.certLock.Lock()
	defer reloader.certLock.Unlock()
	reloader.cert = &newCert
	return nil
}

// GetCertificateFunc implement reloader function
func (reloader *keypairReloader) GetCertificateFunc() func(*tls.ClientHelloInfo) (*tls.Certificate, error) {
	return func(clientHello *tls.ClientHelloInfo) (*tls.Certificate, error) {
		reloader.certLock.RLock()
		defer reloader.certLock.RUnlock()
		return reloader.cert, nil
	}
}

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

	"github.com/tknie/log"
)

type keypairReloader struct {
	certLock sync.RWMutex
	cert     *tls.Certificate
	certPath string
	keyPath  string
}

func NewKeypairReloader(certPath, keyPath string) (*keypairReloader, error) {
	result := &keypairReloader{
		certPath: certPath,
		keyPath:  keyPath,
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
			if err := result.maybeReload(); err != nil {
				log.Log.Debugf("Keeping old TLS certificate because the new one could not be loaded: %v", err)
			}
		}
	}()
	return result, nil
}

func (reloader *keypairReloader) maybeReload() error {
	newCert, err := tls.LoadX509KeyPair(reloader.certPath, reloader.keyPath)
	if err != nil {
		return err
	}
	reloader.certLock.Lock()
	defer reloader.certLock.Unlock()
	reloader.cert = &newCert
	return nil
}

func (reloader *keypairReloader) GetCertificateFunc() func(*tls.ClientHelloInfo) (*tls.Certificate, error) {
	return func(clientHello *tls.ClientHelloInfo) (*tls.Certificate, error) {
		reloader.certLock.RLock()
		defer reloader.certLock.RUnlock()
		return reloader.cert, nil
	}
}

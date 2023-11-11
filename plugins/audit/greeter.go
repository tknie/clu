/*
* Copyright 2022-2023 Thorsten A. Knieling
*
* Licensed under the Apache License, Version 2.0 (the "License");
* you may not use this file except in compliance with the License.
* You may obtain a copy of the License at
*
*    http://www.apache.org/licenses/LICENSE-2.0
*
 */

package main

import (
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/tknie/clu/server"
	"github.com/tknie/log"
	"github.com/tknie/services"
)

type greeting string

type session struct {
	user  string
	uuid  string
	start time.Time
}

var sessionMap sync.Map

// Types type of plugin working with
func (g greeting) Types() []int {
	return []int{1}
}

// Name name of the plugin
func (g greeting) Name() string {
	return "Audit Access"
}

// Version version of the number
func (g greeting) Version() string {
	return "1.1"
}

// Stop stop plugin
func (g greeting) Stop() {
}

// ReceiveAudit receive audit info incoming request
func (g greeting) ReceiveAudit(user string, uuid string, r *http.Request) {
	sessionMap.Store(fmt.Sprintf("%p", r), &session{start: time.Now()})
	log.Log.Infof("Incoming Token %s User: %s Method: %s %s %s Host: %s",
		uuid, user, r.Method, r.RequestURI, r.RemoteAddr, r.Host)
}

// SendAudit audit of http trigger
func (g greeting) SendAudit(elapsed time.Duration, user string, uuid string, w *http.Request) {
	reqURI := strings.ReplaceAll(w.RequestURI, "%", "%%")
	if e, ok := sessionMap.Load(fmt.Sprintf("%p", w)); ok {
		x := e.(*session)

		if strings.HasPrefix(strings.ToLower(w.RequestURI), "/login") {
			services.ServerMessage("Used: %v User: %s -> %s %s -> %s from %s)",
				time.Since(x.start), user, w.Method, reqURI, uuid, server.RemoteHost(w))
			log.Log.Infof("Used: %v User: %s -> %s %s -> %s from %s)",
				time.Since(x.start), user, w.Method, reqURI, uuid, server.RemoteHost(w))
			return
		}
		log.Log.Infof("Used: %v User: %s -> %s %s from %s)",
			time.Since(x.start), user, w.Method, reqURI, server.RemoteHost(w))
		return
	}
	if u, _, ok := w.BasicAuth(); ok {
		services.ServerMessage("Failed: %v Token %s User: %s %s %s %s Host: %s",
			elapsed, uuid, u, w.Method, server.RemoteHost(w), reqURI, w.Host)
		log.Log.Errorf("Failed: %v Token %s User: %s %s %s %s Host: %s",
			elapsed, uuid, u, w.Method, server.RemoteHost(w), reqURI, w.Host)
		return
	}
	services.ServerMessage("Failed: %v Token %s Unknown user %s %s %s Host: %s",
		elapsed, uuid, w.Method, server.RemoteHost(w), reqURI, w.Host)
	log.Log.Errorf("Failed: %v Token %s Unknown user %s %s %s Host: %s",
		elapsed, uuid, w.Method, server.RemoteHost(w), reqURI, w.Host)
}

// SendAuditError audit of http trigger
func (g greeting) SendAuditError(elapsed time.Duration, user string, uuid string, w *http.Request, err error) {
	reqURI := strings.ReplaceAll(w.RequestURI, "%", "%%")
	u, _, b := w.BasicAuth()
	services.ServerMessage("Error access: %s %s %s <%s> BasicAuth: %v Host: %s Error: %v",
		w.Method, server.RemoteHost(w), reqURI, u, b, w.Host, err)
}

// exported

// Loader loader for initialize plugin
var Loader greeting

// Audit audit specific entry methods
var Audit greeting

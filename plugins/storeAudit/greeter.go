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

package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/tknie/clu/plugins"
	"github.com/tknie/clu/server"
	"github.com/tknie/flynn"
	"github.com/tknie/flynn/common"
	"github.com/tknie/log"
	"github.com/tknie/services"
	"github.com/tknie/services/auth"
)

const startEventMethod = "Start"
const stopEventMethod = "Stop"
const adminUser = "Bob@Admin"

var wg sync.WaitGroup

type greeting string

type session struct {
	user  string
	uuid  string
	start time.Time
}

const startSessionInfo = "Init"
const endSessionInfo = "End"
const triggerSessionInfo = "Trigger"

var sessionMap sync.Map
var fieldList = []string{"Triggered", "Elapsed",
	"RequestUser", "UUID", "ServerHost",
	"Method", "RemoteAddr", "RemoteHost",
	"Service", "URI", "Host", "Status",
	"TableName", "AlbumID", "Fields"}

// SessionInfo session info storage
type SessionInfo struct {
	ID         uint64 `flynn:"ID::SERIAL"`
	Triggered  time.Time
	Elapsed    int64
	ServerHost string
	Method     string
	RemoteAddr string
	RemoteHost string
	UUID       string
	RequestURI string `flynn:"URI"`
	Service    string `flynn:"::2024"`
	User       string `flynn:"RequestUser"`
	Host       string
	TableName  string
	AlbumID    uint64
	Fields     string
	Status     string
}

const pluginName = "Audit Log Store"

// var id common.RegDbID
var storeChan = make(chan *SessionInfo, 1000)
var tableName = ""
var disableStore = false
var url string
var dbRef *common.Reference
var password string

func pluginMessage(msg string, argv ...interface{}) {
	services.ServerMessage(pluginName+": "+msg, argv...)
}

// var auditStoreID common.RegDbID

func init() {
	url = os.Getenv("REST_AUDIT_LOG_URL")
	tableName = os.Getenv("REST_AUDIT_LOG_TABLENAME")
	if url == "" || tableName == "" {
		pluginMessage("Log parameter storage disabled...")
		log.Log.Debugf("STORE_AUDIT: Disable due to URL error")
		disableStore = true
		return
	}
	var err error
	log.Log.Debugf("Datbase target %s", url)
	dbRef, password, err = common.NewReference(url)
	if err != nil || dbRef == nil {
		log.Log.Fatal("REST audit URL incorrect: " + url)
	}
	if password == "" {
		password = os.Getenv("REST_AUDIT_LOG_PASS")
	}
	dbRef.User = "admin"

	pluginMessage("Storing audit data to table '%s'", tableName)
	auditStoreID, err := flynn.Handler(dbRef, password)
	if err != nil {
		pluginMessage("Register error log: %v", err)
		return
	}
	log.Log.Debugf("Receive handler %s", auditStoreID)
	defer auditStoreID.FreeHandler()
	defer auditStoreID.Close()

	v := services.BuildVersion
	si := NewSessionInfo(startSessionInfo, v, adminUser, startEventMethod)

	go startStore()

	dbTables := flynn.Maps()
	for _, d := range dbTables {
		if d == tableName {
			// pluginMessage("Database log found")
			storeChan <- si
			return
		}
	}
	err = auditStoreID.CreateTable(tableName, si)
	if err != nil {
		pluginMessage("Databaase log creating failed: %v", err)
		return
	}
	pluginMessage("Databaase log creating succeed")

	storeChan <- si
}

// NewSessionInfo new session info creation
func NewSessionInfo(status, uuid, user, method string) *SessionInfo {
	hostname, _ := os.Hostname()
	return &SessionInfo{Status: status, UUID: uuid,
		Method: method, User: user,
		Triggered:  time.Now(),
		ServerHost: hostname}

}

func (si *SessionInfo) adapt(req *http.Request) {
	si.RemoteHost = server.RemoteHost(req)
	si.extractURI(req)
	si.Host = req.Host
}

func startStore() {
	insert := &common.Entries{Fields: fieldList}
	wg.Add(1)
	defer wg.Done()
	lock := sync.Mutex{}
	defer pluginMessage("Ending store audit log")

	for {
		// log.Log.Debugf("STORE_AUDIT: Waiting store channel (%v)", disableStore)
		select {
		case si := <-storeChan:
			log.Log.Debugf("STORE_AUDIT: Receive store channel (%v)", disableStore)
			if !disableStore {
				lock.Lock()
				auditStoreID, err := flynn.Handler(dbRef, password)
				if err != nil {
					pluginMessage("Register error log: %v", err)
					return
				}
				log.Log.Debugf("STORE_AUDIT: Store channel (%v)", disableStore)
				x := strings.Index(si.RemoteHost, ",")
				addr := si.RemoteAddr
				host := si.RemoteHost
				if x > 0 {
					addr = strings.Trim(si.RemoteHost[:x], " ")
					host = strings.Trim(si.RemoteHost[x+1:], " ")
				}
				insert.Values = [][]any{{si.Triggered, si.Elapsed,
					si.User, si.UUID, si.ServerHost,
					si.Method, addr, host,
					si.Service, si.RequestURI, si.Host, si.Status,
					si.TableName, si.AlbumID, si.Fields}}
				log.Log.Debugf("STORE_AUDIT: Insert store channel (%v)", disableStore)
				_, err = auditStoreID.Insert(tableName, insert)
				if err != nil {
					pluginMessage("Error store to session %s/%s(%s) : %v",
						host, addr, tableName, err)
					log.Log.Debugf("STORE_AUDIT: Disable due to INSERT error")
					disableStore = true
				}
				log.Log.Debugf("STORE_AUDIT: Commit store channel (%v)", disableStore)
				err = auditStoreID.Commit()
				if err != nil {
					pluginMessage("Error commiting to session %s/%s(%s) : %v",
						host, addr, tableName, err)
					log.Log.Debugf("STORE_AUDIT: Disable due to COMMIT error")
					disableStore = true
				}
				if si.Method == stopEventMethod {
					log.Log.Debugf("STORE_AUDIT: Disable due to STOP")
					disableStore = true
				}
				log.Log.Debugf("STORE_AUDIT: free handler (%v) %s", disableStore, auditStoreID)
				auditStoreID.FreeHandler()
				// auditStoreID.Close()
				lock.Unlock()
			}
			log.Log.Debugf("STORE_AUDIT: End store channel (%v)", disableStore)
		case <-time.After(30 * time.Second):
			if len(storeChan) > cap(storeChan)-100 {
				for si := range storeChan {
					fmt.Println("Skip session store of", si.ID)
				}
				log.Log.Debugf("STORE_AUDIT: Disable due to store capacity")
				disableStore = true
			}
		}
		// log.Log.Debugf("STORE_AUDIT: Check disable (%v)", disableStore)
		if disableStore {
			log.Log.Debugf("STORE_AUDIT: store disabled, exiting ...")
			return
		}
		// log.Log.Debugf("STORE_AUDIT: loop (%v)", disableStore)
	}
	// log.Log.Debugf("STORE_AUDIT: exit (%v)", disableStore)
}

// Types type of plugin working with
func (g greeting) Types() []plugins.PluginTypes {
	return []plugins.PluginTypes{plugins.AuditPlugin}
}

// Name name of the plugin
func (g greeting) Name() string {
	return pluginName
}

// Version version of the number
func (g greeting) Version() string {
	return "1.2"
}

// Stop stop plugin
func (g greeting) Stop() {
	if disableStore {
		return
	}
	si := NewSessionInfo(endSessionInfo, services.BuildVersion, adminUser, stopEventMethod)
	storeChan <- si
	wg.Wait()
}

func key(uuid string, r *http.Request) string {
	return uuid + "#" + server.RemoteHost(r) + r.RequestURI
}

// LoginAudit login audit info incoming request
func (g greeting) LoginAudit(method string, status string, session *auth.SessionInfo, user *auth.UserInfo) {
	if disableStore {
		return
	}
	log.Log.Debugf("STORE_AUDIT: login audit %s -> %s", user, status)
	log.Log.Debugf("STORE_AUDIT: login session %#v", session)
	if session == nil {
		si := NewSessionInfo(status, "PREUUID", user.User, method)

		log.Log.Debugf("STORE_AUDIT: Send audit to store channel (%v/%s)", disableStore, si.User)
		storeChan <- si
		return
	}
	si := NewSessionInfo(status, session.UUID, user.User, method)

	log.Log.Debugf("STORE_AUDIT: Send audit to store channel (%v/%s)", disableStore, si.User)
	storeChan <- si
}

// ReceiveAudit receive audit info incoming request
func (g greeting) ReceiveAudit(user string, uuid string, r *http.Request) {
	log.Log.Debugf("STORE_AUDIT: Receive audit to map for user %s", user)
	sessionMap.Store(key(uuid, r), &session{user: user, uuid: uuid, start: time.Now()})
}

// SendAudit audit of http trigger
func (g greeting) SendAudit(elapsed time.Duration, user string, uuid string, w *http.Request) {
	if disableStore {
		return
	}
	log.Log.Debugf("STORE_AUDIT: send audit %s/%s", user, uuid)
	si := NewSessionInfo(triggerSessionInfo, uuid, user, w.Method)
	si.adapt(w)
	si.extractURI(w)
	if e, ok := sessionMap.Load(key(uuid, w)); ok {
		x := e.(*session)
		si.Elapsed = time.Since(x.start).Milliseconds()
		si.User = x.user
	} else if u, _, ok := w.BasicAuth(); ok {
		si.User = u
	} else {
		si.User = user
	}
	log.Log.Debugf("STORE_AUDIT: Send audit to store channel (%v/%s)", disableStore, si.User)
	storeChan <- si
}

func (si *SessionInfo) extractURI(w *http.Request) {
	x := strings.IndexRune(w.RequestURI, '?')
	uri := w.RequestURI
	if x > 0 {
		uri = w.RequestURI[:x]
	}
	uri = removePrefix(uri, "/image/")
	uri = removePrefix(uri, "/rest/view/")
	uriPart := strings.Split(uri, "/")
	fields := []string{}
	if len(uriPart) > 1 {
		fields = strings.Split(uriPart[1], ",")
	}
	aid := uint64(0)
	if len(uriPart) > 2 && strings.Contains(uriPart[2], "albumid=") {
		x := strings.Index(uriPart[2], "albumid=")
		aidValue := uriPart[2][x+8:]
		if aidValue != "" {
			aid, _ = strconv.ParseUint(aidValue, 10, 0)
		}
	}
	si.TableName = uriPart[0]
	si.AlbumID = aid
	si.Fields = fmt.Sprintf("%v", fields)
	si.Service = uriPart[len(uriPart)-1]
}

func removePrefix(uri, prefix string) string {
	x := strings.Index(uri, prefix)
	if x == -1 {
		return uri
	}
	return uri[x+len(prefix):]
}

// SendAuditError audit of http trigger
func (g greeting) SendAuditError(elapsed time.Duration, user string, uuid string, w *http.Request, err error) {
	pluginMessage("Failed: %v User: %s Error: %v",
		elapsed, user, err)
	si := NewSessionInfo(err.Error(), uuid, user, w.Method)
	si.adapt(w)
	si.extractURI(w)
	storeChan <- si
}

// exported

// Loader loader for initialize plugin
var Loader greeting

// Audit audit specific entry methods
var Audit greeting

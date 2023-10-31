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
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/tknie/clu/server"
	"github.com/tknie/flynn"
	"github.com/tknie/flynn/common"
	"github.com/tknie/services"
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

var sessionMap sync.Map
var fieldList = []string{"Triggered", "Elapsed",
	"RequestUser", "UUID", "ServerHost",
	"Method", "RemoteAddr", "RemoteHost",
	"Service", "URI", "Host", "Status",
	"TableName", "AlbumID", "Fields"}

// SessionInfo session info storage
type SessionInfo struct {
	ID         uint64 `dbsql:"ID::SERIAL"`
	Triggered  time.Time
	Elapsed    int64
	ServerHost string
	Method     string
	RemoteAddr string
	RemoteHost string
	UUID       string
	RequestURI string `dbsql:"URI"`
	Service    string
	User       string `dbsql:"RequestUser"`
	Host       string
	TableName  string
	AlbumID    uint64
	Fields     string
	Status     string
}

// var id common.RegDbID
var storeChan = make(chan *SessionInfo, 1000)
var tableName = ""
var disableStore = false
var url string
var dbRef *common.Reference
var password string

func init() {
	go startStore()

	url = os.Getenv("REST_AUDIT_LOG_URL")
	tableName = os.Getenv("REST_AUDIT_LOG_TABLENAME")
	if url == "" || tableName == "" {
		services.ServerMessage("Log parameter storage disabled...")
		disableStore = true
		return
	}
	var err error
	dbRef, password, err = common.NewReference(url)
	if err != nil {
		log.Fatal("REST audit URL incorrect: " + url)
	}
	if password == "" {
		password = os.Getenv("REST_AUDIT_LOG_PASS")
	}
	dbRef.User = "admin"

	services.ServerMessage("Storing audit data to table '%s'", tableName)
	id, err := flynn.RegisterDatabase(dbRef, password)
	if err != nil {
		services.ServerMessage("Register error log: %v", err)
		return
	}

	defer flynn.Unregister(id)

	si := NewSessionInfo("Init", "0000-0000", adminUser, startEventMethod)

	dbTables := flynn.Maps()
	for _, d := range dbTables {
		if d == tableName {
			// services.ServerMessage("Database log found")
			storeChan <- si
			return
		}
	}
	err = id.CreateTable(tableName, si)
	if err != nil {
		services.ServerMessage("Databaase log creating failed: %v", err)
		return
	}
	services.ServerMessage("Databaase log creating succeed")

	storeChan <- si
}

// NewSessionInfo new session info
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
	for {
		select {
		case si := <-storeChan:
			if !disableStore {
				id, err := flynn.RegisterDatabase(dbRef, password)
				if err != nil {
					services.ServerMessage("Register store audit log fails: %v(%s)", err, url)
					return
				}
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
				err = id.Insert(tableName, insert)
				if err != nil {
					services.ServerMessage("Error store to session %s/%s(%s) : %v",
						host, addr, tableName, err)
					disableStore = true
				}
				if si.Method == stopEventMethod {
					disableStore = true
				}
				flynn.Unregister(id)
			}
		case <-time.After(1 * time.Second):
			if len(storeChan) > cap(storeChan)-100 {
				for si := range storeChan {
					fmt.Println("Skip session store of", si.ID)
				}
				disableStore = true
			}
		}
		if disableStore {
			return
		}
	}
}

// Types type of plugin working with
func (g greeting) Types() []int {
	return []int{1}
}

// Name name of the plugin
func (g greeting) Name() string {
	return "Audit Log Store"
}

// Version version of the number
func (g greeting) Version() string {
	return "1.1"
}

// Stop stop plugin
func (g greeting) Stop() {
	if disableStore {
		return
	}
	si := NewSessionInfo("End", "0000-0000", adminUser, stopEventMethod)
	storeChan <- si
	wg.Wait()
}

func key(uuid string, r *http.Request) string {
	return uuid + "#" + server.RemoteHost(r) + r.RequestURI
}

// ReceiveAudit receive audit info incoming request
func (g greeting) ReceiveAudit(user string, uuid string, r *http.Request) {
	sessionMap.Store(key(uuid, r), &session{user: user, uuid: uuid, start: time.Now()})
}

// SendAudit audit of http trigger
func (g greeting) SendAudit(elapsed time.Duration, user string, uuid string, w *http.Request) {
	if disableStore {
		return
	}
	si := NewSessionInfo("Ok", uuid, "Unknown", w.Method)
	si.adapt(w)
	si.extractURI(w)
	if e, ok := sessionMap.Load(key(uuid, w)); ok {
		x := e.(*session)
		si.Elapsed = time.Since(x.start).Milliseconds()
		si.User = x.user
	} else if u, _, ok := w.BasicAuth(); ok {
		si.User = u
	} else {
		si.User = "Unknown"
	}
	storeChan <- si
}

func (si *SessionInfo) extractURI(w *http.Request) {
	x := strings.IndexRune(w.RequestURI, '?')
	uri := w.RequestURI
	if x > 0 {
		uri = w.RequestURI[:x]
	}
	uri = removePrefix(uri, "/image")
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
	services.ServerMessage("Failed: %v User: %s Error: %v",
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

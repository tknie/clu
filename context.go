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

package clu

import (
	"net/http"
	"time"
)

// Audit callback function to be enabled
var Audit func(time.Time, *http.Request, error)

// Context server context
type Context struct {
	User string
	Pass string
	X    struct {
		UUID    string
		Name    string
		Roles   []string
		Remote  string
		Session interface{}
	}
	Token          string
	started        time.Time
	CurrentRequest *http.Request
	dataMap        map[string]any
}

// NewContext new server context with user and password
func NewContext(user, pass string) *Context {
	return &Context{
		User: user, dataMap: make(map[string]any),
		Pass: pass, started: time.Now()}
}

// Deadline dead line
func (sc *Context) Deadline() (deadline time.Time, ok bool) { return time.Now(), false }

// Done context done
func (sc *Context) Done() <-chan struct{} { return make(<-chan struct{}) }

// Err error return
func (sc *Context) Err() error { return nil }

// Value value of key
func (sc *Context) Value(key any) any {
	return nil
}

// UUID UUID interface function
func (sc *Context) UUID() string {
	return sc.X.UUID
}

// UserName user interface function
func (sc *Context) UserName() string {
	return sc.User
}

// Name Name interface function
func (sc *Context) Name() string {
	return sc.X.Name
}

// AddRoles add roles interface function
func (sc *Context) AddRoles(r []string) {
	sc.X.Roles = append(sc.X.Roles, r...)
}

// Remote remote info interface function
func (sc *Context) Remote() string {
	return sc.X.Remote
}

// SetRemote set remote info interface function
func (sc *Context) SetRemote(r string) {
	sc.X.Remote = r
}

// Roles roles info interface function
func (sc *Context) Roles() []string {
	return sc.X.Roles
}

// Session session info interface function
func (sc *Context) Session() interface{} {
	return sc.X.Session
}

// SetSession set session info interface function
func (sc *Context) SetSession(s interface{}) {
	sc.X.Session = s
}

// SendAuditError send audit error in context
func (sc *Context) SendAuditError(started time.Time, err error) {
	if Audit != nil {
		Audit(started, sc.CurrentRequest, err)
	}
}

// StoreData entry specific storage of data
func (sc *Context) StoreData(key string, value any) {
	sc.dataMap[key] = value
}

// GetData entry specific storage of data is requested
func (sc *Context) GetData(key string) any {
	return sc.dataMap[key]
}

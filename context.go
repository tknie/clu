/*
* Copyright 2022-2025 Thorsten A. Knieling
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

	"github.com/tknie/services/auth"
)

// Audit callback function to be enabled
var Audit func(time.Time, *http.Request, error)

// Context server context
type Context struct {
	user *auth.UserInfo
	// User string
	Pass string
	Auth struct {
		Roles   []string
		Remote  string
		Session interface{}
	}
	Token          string
	started        time.Time
	CurrentRequest *http.Request
	dataMap        map[string]any
}

// NewContextUserInfo new server context with user information and password
func NewContextUserInfo(userInfo *auth.UserInfo, pass string) *Context {
	return &Context{
		user: userInfo,
		Auth: struct {
			Roles   []string
			Remote  string
			Session interface{}
		}{Session: auth.NewSessionInfo(userInfo.User)},
		dataMap: make(map[string]any),
		Pass:    pass, started: time.Now()}
}

// NewContext new server context with user and password
//
// Deprecated: Use NewContextUserInfo instead
func NewContext(user, pass string) *Context {
	created := time.Now()
	return NewContextUserInfo(&auth.UserInfo{User: user, Created: created}, pass)
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
	return sc.Auth.Session.(*auth.SessionInfo).UUID
}

// UserName user interface function
func (sc *Context) UserName() string {
	return sc.user.User
}

// Name Name interface function
func (sc *Context) Name() string {
	return sc.user.LongName
}

// AddRoles add roles interface function
func (sc *Context) AddRoles(r []string) {
	sc.Auth.Roles = append(sc.Auth.Roles, r...)
}

// Remote remote info interface function
func (sc *Context) Remote() string {
	return sc.Auth.Remote
}

// SetRemote set remote info interface function
func (sc *Context) SetRemote(r string) {
	sc.Auth.Remote = r
}

// Roles roles info interface function
func (sc *Context) Roles() []string {
	return sc.Auth.Roles
}

// Session session info interface function
func (sc *Context) Session() interface{} {
	return sc.Auth.Session
}

// SetSession set session info interface function
func (sc *Context) SetSession(s interface{}) {
	sc.Auth.Session = s
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

// Permission return user permission
func (sc *Context) Permission() *auth.User {
	return sc.user.Permission
}

// Created context creation time
func (sc *Context) Created() time.Time {
	return sc.user.Created
}

// EMail email for this user
func (sc *Context) EMail() string {
	return sc.user.EMail
}

// Picture thumbnail picture store [optional]
func (sc *Context) Picture() string {
	return string(sc.user.Picture)
}

// LastLogin last login of user
func (sc *Context) LastLogin() time.Time {
	return sc.user.LastLogin
}

// UpdateLastLogin set last login of user
func (sc *Context) UpdateLastLogin() {
	sc.user.LastLogin = time.Now()
}

// LongName long name of user
func (sc *Context) LongName() string {
	return sc.user.LongName
}

// SetLongName set long name of user
func (sc *Context) SetLongName(longName string) {
	sc.user.LongName = longName
}

// User return user information
func (sc *Context) User() *auth.UserInfo {
	return sc.user
}

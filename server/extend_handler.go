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

package server

import (
	"context"
	"net/http"
	"path/filepath"
	"strings"
	"sync"

	ht "github.com/ogen-go/ogen/http"
	"github.com/tknie/clu"
	"github.com/tknie/clu/api"
	"github.com/tknie/log"
	"github.com/tknie/services/auth"
)

var entryMap = sync.Map{}

// RestExtend Extend method to send to plugin
type RestExtend interface {
	EntryPoint() []string
	CallExtendGet(path string, req *http.Request) (r api.CallExtendRes, _ error)
	CallExtendPut(path string, req *http.Request) (r api.TriggerExtendRes, _ error)
	CallExtendPost(path string, req *http.Request) (r api.CallPostExtendRes, _ error)
}

// RegisterExtend register the extend handler
func RegisterExtend(extend RestExtend) {
	for _, e := range extend.EntryPoint() {
		entryMap.Store(e, extend)
	}
}

// CallExtend implements callExtend operation.
//
// Call plugin extend.
//
// GET /rest/extend/{path}
func (Handler) CallExtend(ctx context.Context, params api.CallExtendParams) (r api.CallExtendRes, _ error) {
	session := ctx.(*clu.Context)
	log.Log.Debugf("Generate JWT token")
	log.Log.Debugf("Call extend: %s,%#v", params.Path, params.Params)
	log.Log.Debugf("User Agent: %#v", session.CurrentRequest.UserAgent())
	log.Log.Debugf("Current request form value %#v", session.CurrentRequest.FormValue("xx"))
	e := filepath.Clean(params.Path)
	parts := strings.Split(e, "/")

	if entryPoint, ok := entryMap.Load(parts[0]); ok {
		user := session.User()
		log.Log.Debugf("User: %#v", user.User)
		log.Log.Debugf("Search records for fields %s -> %s", user.User, "#"+parts[0])
		if !Validate(session, auth.UserRole, "#"+parts[0]) {
			return &api.CallExtendForbidden{}, nil
		}
		return entryPoint.(RestExtend).CallExtendGet(e, session.CurrentRequest)
	}
	return r, ht.ErrNotImplemented
}

// CallPostExtend implements callPostExtend operation.
//
// Post extend/plugin.
//
// POST /rest/extend/{path}
func (Handler) CallPostExtend(ctx context.Context, req *api.CallPostExtendReq, params api.CallPostExtendParams) (r api.CallPostExtendRes, _ error) {
	session := ctx.(*clu.Context)
	log.Log.Debugf("Generate JWT token")
	e := filepath.Clean(params.Path)
	parts := strings.Split(e, "/")

	if entryPoint, ok := entryMap.Load(parts[0]); ok {
		user := session.User()
		log.Log.Debugf("User: %#v", user.User)
		log.Log.Debugf("Search records for fields %s -> %s", user.User, "#"+parts[0])
		if !Validate(session, auth.UserRole, "#"+parts[0]) {
			return &api.CallPostExtendForbidden{}, nil
		}
		return entryPoint.(RestExtend).CallExtendPost(e, session.CurrentRequest)
	}
	return r, ht.ErrNotImplemented
}

// TriggerExtend implements triggerExtend operation.
//
// Put extend/plugin request.
//
// PUT /rest/extend/{path}
func (Handler) TriggerExtend(ctx context.Context, params api.TriggerExtendParams) (r api.TriggerExtendRes, _ error) {
	session := ctx.(*clu.Context)
	log.Log.Debugf("Generate JWT token")
	e := filepath.Clean(params.Path)
	parts := strings.Split(e, "/")

	if entryPoint, ok := entryMap.Load(parts[0]); ok {
		user := session.User()
		log.Log.Debugf("User: %#v", user.User)
		log.Log.Debugf("Search records for fields %s -> %s", user.User, "#"+parts[0])
		if !Validate(session, auth.UserRole, "#"+parts[0]) {
			return &api.TriggerExtendForbidden{}, nil
		}
		return entryPoint.(RestExtend).CallExtendPut(e, session.CurrentRequest)
	}
	return r, ht.ErrNotImplemented
}

// DeleteExtend implements deleteExtend operation.
//
// Delete extend/plugin data.
//
// DELETE /rest/extend/{path}
func (Handler) DeleteExtend(ctx context.Context, params api.DeleteExtendParams) (r api.DeleteExtendRes, _ error) {
	return r, ht.ErrNotImplemented
}

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

package server

import (
	"context"
	"path/filepath"
	"strings"
	"sync"

	ht "github.com/ogen-go/ogen/http"
	"github.com/tknie/clu/api"
)

var entryMap = sync.Map{}

// RestExtend Adabas method to send to plugin
type RestExtend interface {
	EntryPoint() string
	Call(path string) (r api.CallExtendRes, _ error)
}

func RegisterExtend(extend RestExtend) {
	entryMap.Store(extend.EntryPoint(), extend)
}

// CallExtend implements callExtend operation.
//
// Call plugin extend.
//
// GET /rest/extend/{path}
func (Handler) CallExtend(ctx context.Context, params api.CallExtendParams) (r api.CallExtendRes, _ error) {
	e := filepath.Clean(params.Path)
	parts := strings.Split(e, "/")

	if entryPoint, ok := entryMap.Load(parts[0]); ok {
		return entryPoint.(RestExtend).Call(e)
	}
	return r, ht.ErrNotImplemented
}

// CallPostExtend implements callPostExtend operation.
//
// Post extend/plugin.
//
// POST /rest/extend/{path}
func (Handler) CallPostExtend(ctx context.Context, req *api.CallPostExtendReq, params api.CallPostExtendParams) (r api.CallPostExtendRes, _ error) {
	return r, ht.ErrNotImplemented
}

// TriggerExtend implements triggerExtend operation.
//
// Put extend/plugin request.
//
// PUT /rest/extend/{path}
func (Handler) TriggerExtend(ctx context.Context, params api.TriggerExtendParams) (r api.TriggerExtendRes, _ error) {
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
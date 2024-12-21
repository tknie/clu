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

package server

import (
	"context"

	ht "github.com/ogen-go/ogen/http"
	"github.com/tknie/clu"
	"github.com/tknie/clu/api"
	"github.com/tknie/services"
	"github.com/tknie/services/auth"
)

// GetVersion implements getVersion operation.
//
// Retrieves the current version.
//
// GET /version
func (Handler) GetVersion(ctx context.Context) (r api.GetVersionRes, _ error) {
	r = &api.Versions{Product: api.NewOptString(clu.ServiceName),
		Version: api.NewOptString(services.BuildVersion)}
	return r, nil
}

// GetUserInfo implements getUserInfo operation.
//
// Retrieves the user information.
//
// GET /rest/user
func (Handler) GetUserInfo(ctx context.Context) (r api.GetUserInfoRes, _ error) {
	return r, ht.ErrNotImplemented
}

// GetMaps implements getMaps operation.
//
// Retrieves a list of available views.
//
// GET /rest/view
func (Handler) GetMaps(ctx context.Context) (r api.GetMapsRes, _ error) {
	session := ctx.(*clu.Context)
	if !auth.ValidUser(auth.UserRole, false, session.User(), "*Maps") {
		return &api.GetMapsForbidden{}, nil
	}

	maps := make([]api.Map, 0)
	for _, m := range clu.GetAllViews() {
		maps = append(maps, api.Map(m))
	}
	r = &api.Maps{Maps: maps}
	return r, nil
}

// GetDatabaseSessions implements getDatabaseSessions operation.
//
// Retrieve a list of user queue entries.
//
// GET /admin/database/{table}/sessions
func (Handler) GetDatabaseSessions(ctx context.Context, params api.GetDatabaseSessionsParams) (r api.GetDatabaseSessionsRes, _ error) {
	return r, ht.ErrNotImplemented
}

// ShutdownServer implements shutdownServer operation.
//
// Init shutdown procedure.
//
// PUT /shutdown/{hash}
func (Handler) ShutdownServer(ctx context.Context, params api.ShutdownServerParams) (r api.ShutdownServerRes, _ error) {
	session := ctx.(*clu.Context)
	if !auth.ValidUser(auth.AdministratorRole, true, session.User(), "") {
		return &api.ShutdownServerForbidden{}, nil
	}

	return r, ht.ErrNotImplemented
}

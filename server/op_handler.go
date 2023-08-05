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
	r = &api.Versions{Product: api.NewOptString("REST-API"),
		Version: api.NewOptString(services.BuildVersion)}
	return r, nil
}

// GetMaps implements getMaps operation.
//
// Retrieves a list of available views.
//
// GET /rest/view
func (Handler) GetMaps(ctx context.Context) (r api.GetMapsRes, _ error) {
	session := ctx.(*clu.Context)
	if !auth.ValidUser(auth.UserRole, false, session.User, "*Maps") {
		return &api.GetMapsForbidden{}, nil
	}

	maps := make([]api.Map, 0)
	for _, m := range GetAllViews() {
		maps = append(maps, api.Map(m))
	}
	r = &api.Maps{Maps: maps}
	return r, nil
}

// GetEnvironments implements getEnvironments operation.
//
// Retrieves the list of environments.
//
// GET /version/env
func (Handler) GetEnvironments(ctx context.Context) (r api.GetEnvironmentsRes, _ error) {
	return r, ht.ErrNotImplemented
}

// GetDatabaseSessions implements getDatabaseSessions operation.
//
// Retrieve a list of user queue entries.
//
// GET /admin/database/{table}/sessions
func (Handler) GetDatabaseSessions(ctx context.Context, params api.GetDatabaseSessionsParams) (r api.GetDatabaseSessionsRes, _ error) {
	return r, ht.ErrNotImplemented
}

// GetDatabaseStats implements getDatabaseStats operation.
//
// Retrieve SQL statistics.
//
// GET /admin/database/{table}/stats
func (Handler) GetDatabaseStats(ctx context.Context, params api.GetDatabaseStatsParams) (r api.GetDatabaseStatsRes, _ error) {
	return r, ht.ErrNotImplemented
}

// DatabaseOperation implements databaseOperation operation.
//
// Retrieve the current status of database with the given dbid.
//
// GET /admin/database/{table_operation}
func (Handler) DatabaseOperation(ctx context.Context, params api.DatabaseOperationParams) (r api.DatabaseOperationRes, _ error) {
	return r, ht.ErrNotImplemented
}

// DatabasePostOperations implements databasePostOperations operation.
//
// Initiate operations on the given dbid.
//
// POST /admin/database/{table_operation}
func (Handler) DatabasePostOperations(ctx context.Context, params api.DatabasePostOperationsParams) (r api.DatabasePostOperationsRes, _ error) {
	return r, ht.ErrNotImplemented
}

// DeleteDatabase implements deleteDatabase operation.
//
// Delete the database.
//
// DELETE /admin/database/{table_operation}
func (Handler) DeleteDatabase(ctx context.Context, params api.DeleteDatabaseParams) (r api.DeleteDatabaseRes, _ error) {
	return r, ht.ErrNotImplemented
}

// ShutdownServer implements shutdownServer operation.
//
// Init shutdown procedure.
//
// PUT /shutdown/{hash}
func (Handler) ShutdownServer(ctx context.Context, params api.ShutdownServerParams) (r api.ShutdownServerRes, _ error) {
	session := ctx.(*clu.Context)
	if !auth.ValidUser(auth.AdministratorRole, false, session.User, "") {
		return &api.ShutdownServerForbidden{}, nil
	}

	return r, ht.ErrNotImplemented
}

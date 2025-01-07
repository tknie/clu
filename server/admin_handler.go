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

	ht "github.com/ogen-go/ogen/http"
	"github.com/tknie/clu/api"
)

// AdaptPermission implements adaptPermission operation.
//
// Add RBAC role.
//
// PUT /admin/database/{table}/permission
func (Handler) AdaptPermission(ctx context.Context, params api.AdaptPermissionParams) (r api.AdaptPermissionRes, _ error) {
	return r, ht.ErrNotImplemented
}

// AddView implements addView operation.
//
// Add configuration in View repositories.
//
// POST /admin/config/views
func (Handler) AddView(ctx context.Context, params api.AddViewParams) (r api.AddViewRes, _ error) {
	return r, ht.ErrNotImplemented
}

// DeleteView implements deleteView operation.
//
// Delete entry in configuration.
//
// DELETE /admin/config/views
func (Handler) DeleteView(ctx context.Context, params api.DeleteViewParams) (r api.DeleteViewRes, _ error) {
	return r, ht.ErrNotImplemented
}

// DisconnectTCP implements disconnectTCP operation.
//
// Disconnect connection in the database with the given dbid.
//
// DELETE /admin/database/{table}/connection
func (Handler) DisconnectTCP(ctx context.Context, params api.DisconnectTCPParams) (r api.DisconnectTCPRes, _ error) {
	return r, ht.ErrNotImplemented
}

// GetConfig implements getConfig operation.
//
// Get configuration.
//
// GET /admin/config
func (Handler) GetConfig(ctx context.Context) (r api.GetConfigRes, _ error) {
	return r, ht.ErrNotImplemented
}

// GetConnections implements getConnections operation.
//
// Retrieve the current TCP connection.
//
// GET /admin/database/{table}/connection
func (Handler) GetConnections(ctx context.Context, params api.GetConnectionsParams) (r api.GetConnectionsRes, _ error) {
	return r, ht.ErrNotImplemented
}

// GetDatabases implements getDatabases operation.
//
// Retrieves a list of databases known by Interface.
//
// GET /admin/database
func (Handler) GetDatabases(ctx context.Context) (r api.GetDatabasesRes, _ error) {
	return r, ht.ErrNotImplemented
}

// GetPermission implements getPermission operation.
//
// List RBAC assignments permission.
//
// GET /admin/database/{table}/permission
func (Handler) GetPermission(ctx context.Context, params api.GetPermissionParams) (r api.GetPermissionRes, _ error) {
	return r, ht.ErrNotImplemented
}

// GetViews implements getViews operation.
//
// Defines the current views.
//
// GET /admin/config/views
func (Handler) GetViews(ctx context.Context) (r api.GetViewsRes, _ error) {
	return r, ht.ErrNotImplemented
}

// PostDatabase implements postDatabase operation.
//
// Create a new database, the input need to be JSON. A structure level parameter indicate version to
// be used.
//
// POST /admin/database
func (Handler) PostDatabase(ctx context.Context, req *api.Database) (r api.PostDatabaseRes, _ error) {
	return r, ht.ErrNotImplemented
}

// RemovePermission implements removePermission operation.
//
// Add RBAC role.
//
// DELETE /admin/database/{table}/permission
func (Handler) RemovePermission(ctx context.Context, params api.RemovePermissionParams) (r api.RemovePermissionRes, _ error) {
	return r, ht.ErrNotImplemented
}

// SetConfig implements setConfig operation.
//
// Store configuration.
//
// PUT /admin/config
func (Handler) SetConfig(ctx context.Context, req api.SetConfigReq) (r api.SetConfigRes, _ error) {
	return r, ht.ErrNotImplemented
}

// SetJobsConfig implements setJobsConfig operation.
//
// Set the ADADATADIR.
//
// PUT /admin/config/jobs
func (Handler) SetJobsConfig(ctx context.Context, req api.OptJobStore) (r api.SetJobsConfigRes, _ error) {
	return r, ht.ErrNotImplemented
}

// StoreConfig implements storeConfig operation.
//
// Store configuration.
//
// POST /admin/config
func (Handler) StoreConfig(ctx context.Context) (r api.StoreConfigRes, _ error) {
	return r, ht.ErrNotImplemented
}

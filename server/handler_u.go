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
	"github.com/tknie/clu/api"
)

// Access implements access operation.
//
// Retrieve the list of users who are allowed to access data.
//
// GET /admin/access/{role}
func (ServerHandler) Access(ctx context.Context, params api.AccessParams) (r api.AccessRes, _ error) {
	return r, ht.ErrNotImplemented
}

// AdaptPermission implements adaptPermission operation.
//
// Add RBAC role.
//
// PUT /admin/database/{table}/permission
func (ServerHandler) AdaptPermission(ctx context.Context, params api.AdaptPermissionParams) (r api.AdaptPermissionRes, _ error) {
	return r, ht.ErrNotImplemented
}

// AddAccess implements addAccess operation.
//
// Insert user in the list of users who are allowed to access data.
//
// POST /admin/access/{role}
func (ServerHandler) AddAccess(ctx context.Context, params api.AddAccessParams) (r api.AddAccessRes, _ error) {
	return r, ht.ErrNotImplemented
}

// AddRBACResource implements addRBACResource operation.
//
// Add permission role.
//
// PUT /admin/database/{table}/permission/{resource}/{name}
func (ServerHandler) AddRBACResource(ctx context.Context, params api.AddRBACResourceParams) (r api.AddRBACResourceRes, _ error) {
	return r, ht.ErrNotImplemented
}

// AddView implements addView operation.
//
// Add configuration in View repositories.
//
// POST /admin/config/views
func (ServerHandler) AddView(ctx context.Context, params api.AddViewParams) (r api.AddViewRes, _ error) {
	return r, ht.ErrNotImplemented
}

// Browse implements browse operation.
//
// Retrieves a list of files in the defined location.
//
// GET /rest/file/browse/{location}
func (ServerHandler) Browse(ctx context.Context, params api.BrowseParams) (r api.BrowseRes, _ error) {
	return r, ht.ErrNotImplemented
}

// BrowseList implements browseList operation.
//
// Retrieves a list of Browseable locations.
//
// GET /rest/file/browse
func (ServerHandler) BrowseList(ctx context.Context) (r api.BrowseListRes, _ error) {
	return r, ht.ErrNotImplemented
}

// CreateDirectory implements createDirectory operation.
//
// Create a new directory.
//
// PUT /rest/file/{location}
func (ServerHandler) CreateDirectory(ctx context.Context, params api.CreateDirectoryParams) (r api.CreateDirectoryRes, _ error) {
	return r, ht.ErrNotImplemented
}

// DatabaseOperation implements databaseOperation operation.
//
// Retrieve the current status of database with the given dbid.
//
// GET /admin/database/{table_operation}
func (ServerHandler) DatabaseOperation(ctx context.Context, params api.DatabaseOperationParams) (r api.DatabaseOperationRes, _ error) {
	return r, ht.ErrNotImplemented
}

// DatabasePostOperations implements databasePostOperations operation.
//
// Initiate operations on the given dbid.
//
// POST /admin/database/{table_operation}
func (ServerHandler) DatabasePostOperations(ctx context.Context, params api.DatabasePostOperationsParams) (r api.DatabasePostOperationsRes, _ error) {
	return r, ht.ErrNotImplemented
}

// DelAccess implements delAccess operation.
//
// Delete user in the list of users who are allowed to access data.
//
// DELETE /admin/access/{role}
func (ServerHandler) DelAccess(ctx context.Context, params api.DelAccessParams) (r api.DelAccessRes, _ error) {
	return r, ht.ErrNotImplemented
}

// DeleteDatabase implements deleteDatabase operation.
//
// Delete the database.
//
// DELETE /admin/database/{table_operation}
func (ServerHandler) DeleteDatabase(ctx context.Context, params api.DeleteDatabaseParams) (r api.DeleteDatabaseRes, _ error) {
	return r, ht.ErrNotImplemented
}

// DeleteFileLocation implements deleteFileLocation operation.
//
// Delete the file on the given location.
//
// DELETE /rest/file/{location}
func (ServerHandler) DeleteFileLocation(ctx context.Context, params api.DeleteFileLocationParams) (r api.DeleteFileLocationRes, _ error) {
	return r, ht.ErrNotImplemented
}

// DeleteJobResult implements deleteJobResult operation.
//
// Delete a specific job result.
//
// DELETE /tasks/{jobName}/result/{jobId}
func (ServerHandler) DeleteJobResult(ctx context.Context, params api.DeleteJobResultParams) (r api.DeleteJobResultRes, _ error) {
	return r, ht.ErrNotImplemented
}

// DeleteRBACResource implements deleteRBACResource operation.
//
// Delete RBAC role.
//
// DELETE /admin/database/{table}/permission/{resource}/{name}
func (ServerHandler) DeleteRBACResource(ctx context.Context, params api.DeleteRBACResourceParams) (r api.DeleteRBACResourceRes, _ error) {
	return r, ht.ErrNotImplemented
}

// DeleteView implements deleteView operation.
//
// Delete entry in configuration.
//
// DELETE /admin/config/views
func (ServerHandler) DeleteView(ctx context.Context, params api.DeleteViewParams) (r api.DeleteViewRes, _ error) {
	return r, ht.ErrNotImplemented
}

// DisconnectTCP implements disconnectTCP operation.
//
// Disconnect ADATCP connection in the database with the given dbid.
//
// DELETE /admin/database/{table}/connection
func (ServerHandler) DisconnectTCP(ctx context.Context, params api.DisconnectTCPParams) (r api.DisconnectTCPRes, _ error) {
	return r, ht.ErrNotImplemented
}

// DownloadFile implements downloadFile operation.
//
// Download a file out of file location.
//
// GET /rest/file/{location}
func (ServerHandler) DownloadFile(ctx context.Context, params api.DownloadFileParams) (r api.DownloadFileRes, _ error) {
	return r, ht.ErrNotImplemented
}

// GetConfig implements getConfig operation.
//
// Get configuration.
//
// GET /admin/config
func (ServerHandler) GetConfig(ctx context.Context) (r api.GetConfigRes, _ error) {
	return r, ht.ErrNotImplemented
}

// GetConnections implements getConnections operation.
//
// Retrieve the current TCP connection.
//
// GET /admin/database/{table}/connection
func (ServerHandler) GetConnections(ctx context.Context, params api.GetConnectionsParams) (r api.GetConnectionsRes, _ error) {
	return r, ht.ErrNotImplemented
}

// GetDatabaseSessions implements getDatabaseSessions operation.
//
// Retrieve a list of user queue entries.
//
// GET /admin/database/{table}/sessions
func (ServerHandler) GetDatabaseSessions(ctx context.Context, params api.GetDatabaseSessionsParams) (r api.GetDatabaseSessionsRes, _ error) {
	return r, ht.ErrNotImplemented
}

// GetDatabaseStats implements getDatabaseStats operation.
//
// Retrieve SQL statistics.
//
// GET /admin/database/{table}/stats
func (ServerHandler) GetDatabaseStats(ctx context.Context, params api.GetDatabaseStatsParams) (r api.GetDatabaseStatsRes, _ error) {
	return r, ht.ErrNotImplemented
}

// GetDatabases implements getDatabases operation.
//
// Retrieves a list of databases known by Interface.
//
// GET /admin/database
func (ServerHandler) GetDatabases(ctx context.Context) (r api.GetDatabasesRes, _ error) {
	return r, ht.ErrNotImplemented
}

// GetEnvironments implements getEnvironments operation.
//
// Retrieves the list of environments.
//
// GET /version/env
func (ServerHandler) GetEnvironments(ctx context.Context) (r api.GetEnvironmentsRes, _ error) {
	return r, ht.ErrNotImplemented
}

// GetFields implements getFields operation.
//
// Retrieves all fields of an file.
//
// GET /rest/tables/{table}/fields
func (ServerHandler) GetFields(ctx context.Context, params api.GetFieldsParams) (r api.GetFieldsRes, _ error) {
	return r, ht.ErrNotImplemented
}

// GetJobExecutionResult implements getJobExecutionResult operation.
//
// Retrieves a specific job result.
//
// GET /tasks/results
func (ServerHandler) GetJobExecutionResult(ctx context.Context, params api.GetJobExecutionResultParams) (r api.GetJobExecutionResultRes, _ error) {
	return r, ht.ErrNotImplemented
}

// GetJobFullInfo implements getJobFullInfo operation.
//
// Retrieves a full job definition.
//
// GET /tasks/{jobName}/full
func (ServerHandler) GetJobFullInfo(ctx context.Context, params api.GetJobFullInfoParams) (r api.GetJobFullInfoRes, _ error) {
	return r, ht.ErrNotImplemented
}

// GetJobs implements getJobs operation.
//
// Retrieves a list of jobs known by the Interface.
//
// GET /tasks
func (ServerHandler) GetJobs(ctx context.Context, params api.GetJobsParams) (r api.GetJobsRes, _ error) {
	return r, ht.ErrNotImplemented
}

// GetJobsConfig implements getJobsConfig operation.
//
// Read job configuration section.
//
// GET /admin/config/jobs
func (ServerHandler) GetJobsConfig(ctx context.Context) (r api.GetJobsConfigRes, _ error) {
	return r, ht.ErrNotImplemented
}

// GetMapMetadata implements getMapMetadata operation.
//
// Retrieves metadata of a Map definition.
//
// GET /rest/metadata/view/{table}
func (ServerHandler) GetMapMetadata(ctx context.Context, params api.GetMapMetadataParams) (r api.GetMapMetadataRes, _ error) {
	return r, ht.ErrNotImplemented
}

// GetPermission implements getPermission operation.
//
// List RBAC assignments permission.
//
// GET /admin/database/{table}/permission
func (ServerHandler) GetPermission(ctx context.Context, params api.GetPermissionParams) (r api.GetPermissionRes, _ error) {
	return r, ht.ErrNotImplemented
}

// GetViews implements getViews operation.
//
// Defines the current views.
//
// GET /admin/config/views
func (ServerHandler) GetViews(ctx context.Context) (r api.GetViewsRes, _ error) {
	return r, ht.ErrNotImplemented
}

// InsertMapFileRecords implements insertMapFileRecords operation.
//
// Store send records into Map definition.
//
// POST /rest/view
func (ServerHandler) InsertMapFileRecords(ctx context.Context, req api.OptInsertMapFileRecordsReq) (r api.InsertMapFileRecordsRes, _ error) {
	return r, ht.ErrNotImplemented
}

// ListRBACResource implements listRBACResource operation.
//
// Add RBAC role.
//
// GET /admin/database/{table}/permission/{resource}
func (ServerHandler) ListRBACResource(ctx context.Context, params api.ListRBACResourceParams) (r api.ListRBACResourceRes, _ error) {
	return r, ht.ErrNotImplemented
}

// PostDatabase implements postDatabase operation.
//
// Create a new database, the input need to be JSON. A structure level parameter indicate version to
// be used.
//
// POST /admin/database
func (ServerHandler) PostDatabase(ctx context.Context, req *api.Database) (r api.PostDatabaseRes, _ error) {
	return r, ht.ErrNotImplemented
}

// PostJob implements postJob operation.
//
// Create a new Job database.
//
// POST /tasks
func (ServerHandler) PostJob(ctx context.Context, req api.PostJobReq) (r api.PostJobRes, _ error) {
	return r, ht.ErrNotImplemented
}

// PutDatabaseResource implements putDatabaseResource operation.
//
// Change resource of the database.
//
// PUT /admin/database/{table_operation}
func (ServerHandler) PutDatabaseResource(ctx context.Context, params api.PutDatabaseResourceParams) (r api.PutDatabaseResourceRes, _ error) {
	return r, ht.ErrNotImplemented
}

// RemovePermission implements removePermission operation.
//
// Add RBAC role.
//
// DELETE /admin/database/{table}/permission
func (ServerHandler) RemovePermission(ctx context.Context, params api.RemovePermissionParams) (r api.RemovePermissionRes, _ error) {
	return r, ht.ErrNotImplemented
}

// RemoveSessionCompat implements removeSessionCompat operation.
//
// Remove the session.
//
// GET /logoff
func (ServerHandler) RemoveSessionCompat(ctx context.Context) (r api.RemoveSessionCompatRes, _ error) {
	return r, ht.ErrNotImplemented
}

// SearchTable implements searchTable operation.
//
// Retrieves all fields of an file.
//
// GET /rest/tables/{table}/{search}
func (ServerHandler) SearchTable(ctx context.Context, params api.SearchTableParams) (r api.SearchTableRes, _ error) {
	return r, ht.ErrNotImplemented
}

// SetConfig implements setConfig operation.
//
// Store configuration.
//
// PUT /admin/config
func (ServerHandler) SetConfig(ctx context.Context, req api.SetConfigReq) (r api.SetConfigRes, _ error) {
	return r, ht.ErrNotImplemented
}

// SetJobsConfig implements setJobsConfig operation.
//
// Set the ADADATADIR.
//
// PUT /admin/config/jobs
func (ServerHandler) SetJobsConfig(ctx context.Context, req api.OptJobStore) (r api.SetJobsConfigRes, _ error) {
	return r, ht.ErrNotImplemented
}

// ShutdownServer implements shutdownServer operation.
//
// Init shutdown procedure.
//
// PUT /shutdown/{hash}
func (ServerHandler) ShutdownServer(ctx context.Context, params api.ShutdownServerParams) (r api.ShutdownServerRes, _ error) {
	return r, ht.ErrNotImplemented
}

// StoreConfig implements storeConfig operation.
//
// Store configuration.
//
// POST /admin/config
func (ServerHandler) StoreConfig(ctx context.Context) (r api.StoreConfigRes, _ error) {
	return r, ht.ErrNotImplemented
}

// UpdateLobByMap implements updateLobByMap operation.
//
// Set a lob at a specific ISN of an field in a Map.
//
// PUT /binary/{table}/{field}/{search}
func (ServerHandler) UpdateLobByMap(ctx context.Context, req api.UpdateLobByMapReq, params api.UpdateLobByMapParams) (r api.UpdateLobByMapRes, _ error) {
	return r, ht.ErrNotImplemented
}

// UploadFile implements uploadFile operation.
//
// Upload a new file to the given location.
//
// POST /rest/file/{location}
func (ServerHandler) UploadFile(ctx context.Context, req *api.UploadFileReq, params api.UploadFileParams) (r api.UploadFileRes, _ error) {
	return r, ht.ErrNotImplemented
}

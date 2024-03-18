// Code generated by ogen, DO NOT EDIT.

package api

import (
	"context"
)

// Handler handles operations described by OpenAPI v3 specification.
type Handler interface {
	// Access implements access operation.
	//
	// Retrieve the list of users who are allowed to access data.
	//
	// GET /admin/access/{role}
	Access(ctx context.Context, params AccessParams) (AccessRes, error)
	// AdaptPermission implements adaptPermission operation.
	//
	// Add RBAC role.
	//
	// PUT /rest/database/{table}/permission
	AdaptPermission(ctx context.Context, params AdaptPermissionParams) (AdaptPermissionRes, error)
	// AddAccess implements addAccess operation.
	//
	// Insert user in the list of users who are allowed to access data.
	//
	// POST /admin/access/{role}
	AddAccess(ctx context.Context, params AddAccessParams) (AddAccessRes, error)
	// AddRBACResource implements addRBACResource operation.
	//
	// Add permission role.
	//
	// PUT /rest/database/{table}/permission/{resource}/{name}
	AddRBACResource(ctx context.Context, params AddRBACResourceParams) (AddRBACResourceRes, error)
	// AddView implements addView operation.
	//
	// Add configuration in View repositories.
	//
	// POST /config/views
	AddView(ctx context.Context, params AddViewParams) (AddViewRes, error)
	// BatchParameterQuery implements batchParameterQuery operation.
	//
	// Call a SQL query batch command posted in query.
	//
	// GET /rest/batch/{table}/{query}
	BatchParameterQuery(ctx context.Context, params BatchParameterQueryParams) (BatchParameterQueryRes, error)
	// BatchQuery implements batchQuery operation.
	//
	// Call a SQL query batch command posted in body.
	//
	// POST /rest/batch/{table}
	BatchQuery(ctx context.Context, req BatchQueryReq, params BatchQueryParams) (BatchQueryRes, error)
	// BatchSelect implements batchSelect operation.
	//
	// Call a SQL query batch command out of the stored query list.
	//
	// GET /rest/batch/{table}
	BatchSelect(ctx context.Context, params BatchSelectParams) (BatchSelectRes, error)
	// BrowseList implements browseList operation.
	//
	// Retrieves a list of Browseable locations.
	//
	// GET /rest/file/browse
	BrowseList(ctx context.Context) (BrowseListRes, error)
	// BrowseLocation implements browseLocation operation.
	//
	// Retrieves a list of files in the defined location.
	//
	// GET /rest/file/browse/{path}
	BrowseLocation(ctx context.Context, params BrowseLocationParams) (BrowseLocationRes, error)
	// CallExtend implements callExtend operation.
	//
	// Call plugin extend.
	//
	// GET /rest/extend/{path}
	CallExtend(ctx context.Context, params CallExtendParams) (CallExtendRes, error)
	// CallPostExtend implements callPostExtend operation.
	//
	// Post extend/plugin.
	//
	// POST /rest/extend/{path}
	CallPostExtend(ctx context.Context, req *CallPostExtendReq, params CallPostExtendParams) (CallPostExtendRes, error)
	// CreateDirectory implements createDirectory operation.
	//
	// Create a new directory.
	//
	// PUT /rest/file/{path}
	CreateDirectory(ctx context.Context, params CreateDirectoryParams) (CreateDirectoryRes, error)
	// DatabaseOperation implements databaseOperation operation.
	//
	// Retrieve the current status of database with the given dbid.
	//
	// GET /rest/database/{table_operation}
	DatabaseOperation(ctx context.Context, params DatabaseOperationParams) (DatabaseOperationRes, error)
	// DatabasePostOperations implements databasePostOperations operation.
	//
	// Initiate operations on the given dbid.
	//
	// POST /rest/database/{table_operation}
	DatabasePostOperations(ctx context.Context, params DatabasePostOperationsParams) (DatabasePostOperationsRes, error)
	// DelAccess implements delAccess operation.
	//
	// Delete user in the list of users who are allowed to access data.
	//
	// DELETE /admin/access/{role}
	DelAccess(ctx context.Context, params DelAccessParams) (DelAccessRes, error)
	// DeleteDatabase implements deleteDatabase operation.
	//
	// Delete the database.
	//
	// DELETE /rest/database/{table_operation}
	DeleteDatabase(ctx context.Context, params DeleteDatabaseParams) (DeleteDatabaseRes, error)
	// DeleteExtend implements deleteExtend operation.
	//
	// Delete extend/plugin data.
	//
	// DELETE /rest/extend/{path}
	DeleteExtend(ctx context.Context, params DeleteExtendParams) (DeleteExtendRes, error)
	// DeleteFileLocation implements deleteFileLocation operation.
	//
	// Delete the file on the given location.
	//
	// DELETE /rest/file/{path}
	DeleteFileLocation(ctx context.Context, params DeleteFileLocationParams) (DeleteFileLocationRes, error)
	// DeleteJobResult implements deleteJobResult operation.
	//
	// Delete a specific job result.
	//
	// DELETE /rest/tasks/{jobName}/{jobId}
	DeleteJobResult(ctx context.Context, params DeleteJobResultParams) (DeleteJobResultRes, error)
	// DeleteRBACResource implements deleteRBACResource operation.
	//
	// Delete RBAC role.
	//
	// DELETE /rest/database/{table}/permission/{resource}/{name}
	DeleteRBACResource(ctx context.Context, params DeleteRBACResourceParams) (DeleteRBACResourceRes, error)
	// DeleteRecordsSearched implements deleteRecordsSearched operation.
	//
	// Delete a record with a given search.
	//
	// DELETE /rest/view/{table}/{search}
	DeleteRecordsSearched(ctx context.Context, params DeleteRecordsSearchedParams) (DeleteRecordsSearchedRes, error)
	// DeleteView implements deleteView operation.
	//
	// Delete entry in configuration.
	//
	// DELETE /config/views
	DeleteView(ctx context.Context, params DeleteViewParams) (DeleteViewRes, error)
	// DisconnectTCP implements disconnectTCP operation.
	//
	// Disconnect connection in the database with the given dbid.
	//
	// DELETE /rest/database/{table}/connection
	DisconnectTCP(ctx context.Context, params DisconnectTCPParams) (DisconnectTCPRes, error)
	// DownloadFile implements downloadFile operation.
	//
	// Download a file out of file location.
	//
	// GET /rest/file/{path}
	DownloadFile(ctx context.Context, params DownloadFileParams) (DownloadFileRes, error)
	// GetConfig implements getConfig operation.
	//
	// Get configuration.
	//
	// GET /config
	GetConfig(ctx context.Context) (GetConfigRes, error)
	// GetConnections implements getConnections operation.
	//
	// Retrieve the current TCP connection.
	//
	// GET /rest/database/{table}/connection
	GetConnections(ctx context.Context, params GetConnectionsParams) (GetConnectionsRes, error)
	// GetDatabaseSessions implements getDatabaseSessions operation.
	//
	// Retrieve a list of user queue entries.
	//
	// GET /rest/database/{table}/sessions
	GetDatabaseSessions(ctx context.Context, params GetDatabaseSessionsParams) (GetDatabaseSessionsRes, error)
	// GetDatabaseStats implements getDatabaseStats operation.
	//
	// Retrieve SQL statistics.
	//
	// GET /rest/database/{table}/stats
	GetDatabaseStats(ctx context.Context, params GetDatabaseStatsParams) (GetDatabaseStatsRes, error)
	// GetDatabases implements getDatabases operation.
	//
	// Retrieves a list of databases known by Interface.
	//
	// GET /rest/database
	GetDatabases(ctx context.Context) (GetDatabasesRes, error)
	// GetEnvironments implements getEnvironments operation.
	//
	// Retrieves the list of environments.
	//
	// GET /rest/env
	GetEnvironments(ctx context.Context) (GetEnvironmentsRes, error)
	// GetFields implements getFields operation.
	//
	// Retrieves all fields of an file.
	//
	// GET /rest/tables/{table}/fields
	GetFields(ctx context.Context, params GetFieldsParams) (GetFieldsRes, error)
	// GetImage implements getImage operation.
	//
	// Retrieves a field of a specific ISN of a Map definition.
	//
	// GET /image/{table}/{field}/{search}
	GetImage(ctx context.Context, params GetImageParams) (GetImageRes, error)
	// GetJobExecutionResult implements getJobExecutionResult operation.
	//
	// Retrieves a specific job result.
	//
	// GET /rest/tasks/results
	GetJobExecutionResult(ctx context.Context, params GetJobExecutionResultParams) (GetJobExecutionResultRes, error)
	// GetJobFullInfo implements getJobFullInfo operation.
	//
	// Retrieves a full job definition.
	//
	// GET /rest/tasks/{jobName}
	GetJobFullInfo(ctx context.Context, params GetJobFullInfoParams) (GetJobFullInfoRes, error)
	// GetJobResult implements getJobResult operation.
	//
	// Delete a specific job result.
	//
	// GET /rest/tasks/{jobName}/{jobId}
	GetJobResult(ctx context.Context, params GetJobResultParams) (GetJobResultRes, error)
	// GetJobs implements getJobs operation.
	//
	// Retrieves a list of jobs known by the Interface.
	//
	// GET /rest/tasks
	GetJobs(ctx context.Context, params GetJobsParams) (GetJobsRes, error)
	// GetJobsConfig implements getJobsConfig operation.
	//
	// Read job configuration section.
	//
	// GET /config/jobs
	GetJobsConfig(ctx context.Context) (GetJobsConfigRes, error)
	// GetLobByMap implements getLobByMap operation.
	//
	// Retrieves a lob of a specific ISN of an field in a Map.
	//
	// GET /binary/{table}/{field}/{search}
	GetLobByMap(ctx context.Context, params GetLobByMapParams) (GetLobByMapRes, error)
	// GetLoginSession implements getLoginSession operation.
	//
	// Login receiving JWT.
	//
	// GET /login
	GetLoginSession(ctx context.Context) (GetLoginSessionRes, error)
	// GetMapMetadata implements getMapMetadata operation.
	//
	// Retrieves metadata of a Map definition.
	//
	// GET /rest/metadata/view/{table}
	GetMapMetadata(ctx context.Context, params GetMapMetadataParams) (GetMapMetadataRes, error)
	// GetMapRecordsFields implements getMapRecordsFields operation.
	//
	// Retrieves a field of a specific ISN of a Map definition.
	//
	// GET /rest/view/{table}/{fields}/{search}
	GetMapRecordsFields(ctx context.Context, params GetMapRecordsFieldsParams) (GetMapRecordsFieldsRes, error)
	// GetMaps implements getMaps operation.
	//
	// Retrieves a list of available views.
	//
	// GET /rest/view
	GetMaps(ctx context.Context) (GetMapsRes, error)
	// GetPermission implements getPermission operation.
	//
	// List RBAC assignments permission.
	//
	// GET /rest/database/{table}/permission
	GetPermission(ctx context.Context, params GetPermissionParams) (GetPermissionRes, error)
	// GetUserInfo implements getUserInfo operation.
	//
	// Retrieves the user information.
	//
	// GET /rest/user
	GetUserInfo(ctx context.Context) (GetUserInfoRes, error)
	// GetVersion implements getVersion operation.
	//
	// Retrieves the current version.
	//
	// GET /version
	GetVersion(ctx context.Context) (GetVersionRes, error)
	// GetVideo implements getVideo operation.
	//
	// Retrieves a video stream of a specific ISN of a Map definition.
	//
	// GET /video/{table}/{field}/{search}
	GetVideo(ctx context.Context, params GetVideoParams) (GetVideoRes, error)
	// GetViews implements getViews operation.
	//
	// Defines the current views.
	//
	// GET /config/views
	GetViews(ctx context.Context) (GetViewsRes, error)
	// InsertMapFileRecords implements insertMapFileRecords operation.
	//
	// Store send records into Map definition.
	//
	// POST /rest/view
	InsertMapFileRecords(ctx context.Context, req OptInsertMapFileRecordsReq) (InsertMapFileRecordsRes, error)
	// InsertRecord implements insertRecord operation.
	//
	// Insert given record.
	//
	// POST /rest/view/{table}
	InsertRecord(ctx context.Context, req OptInsertRecordReq, params InsertRecordParams) (InsertRecordRes, error)
	// ListModelling implements listModelling operation.
	//
	// Retrieves all tables, views or data representation objects.
	//
	// GET /rest/map
	ListModelling(ctx context.Context) (ListModellingRes, error)
	// ListRBACResource implements listRBACResource operation.
	//
	// Add permission role.
	//
	// GET /rest/database/{table}/permission/{resource}
	ListRBACResource(ctx context.Context, params ListRBACResourceParams) (ListRBACResourceRes, error)
	// ListTables implements listTables operation.
	//
	// Retrieves all tables of databases.
	//
	// GET /rest/tables
	ListTables(ctx context.Context) (ListTablesRes, error)
	// LoginSession implements loginSession operation.
	//
	// Login receiving JWT.
	//
	// PUT /login
	LoginSession(ctx context.Context) (LoginSessionRes, error)
	// LogoutSessionCompat implements logoutSessionCompat operation.
	//
	// Logout the session.
	//
	// PUT /logout
	LogoutSessionCompat(ctx context.Context) (LogoutSessionCompatRes, error)
	// PostDatabase implements postDatabase operation.
	//
	// Create a new database, the input need to be JSON. A structure level parameter indicate version to
	// be used.
	//
	// POST /rest/database
	PostDatabase(ctx context.Context, req *Database) (PostDatabaseRes, error)
	// PostJob implements postJob operation.
	//
	// Create a new Job database.
	//
	// POST /rest/tasks
	PostJob(ctx context.Context, req PostJobReq) (PostJobRes, error)
	// PushLoginSession implements pushLoginSession operation.
	//
	// Login receiving JWT.
	//
	// POST /login
	PushLoginSession(ctx context.Context) (PushLoginSessionRes, error)
	// PutDatabaseResource implements putDatabaseResource operation.
	//
	// Change resource of the database.
	//
	// PUT /rest/database/{table_operation}
	PutDatabaseResource(ctx context.Context, params PutDatabaseResourceParams) (PutDatabaseResourceRes, error)
	// RemovePermission implements removePermission operation.
	//
	// Add RBAC role.
	//
	// DELETE /rest/database/{table}/permission
	RemovePermission(ctx context.Context, params RemovePermissionParams) (RemovePermissionRes, error)
	// RemoveSessionCompat implements removeSessionCompat operation.
	//
	// Remove the session.
	//
	// GET /logoff
	RemoveSessionCompat(ctx context.Context) (RemoveSessionCompatRes, error)
	// SearchModelling implements searchModelling operation.
	//
	// Retrieves all columns, fields of a tables, views or data representation.
	//
	// GET /rest/map/{path}
	SearchModelling(ctx context.Context, params SearchModellingParams) (SearchModellingRes, error)
	// SearchRecordsFields implements searchRecordsFields operation.
	//
	// Query a record with a given SQL query.
	//
	// GET /rest/view/{table}/{search}
	SearchRecordsFields(ctx context.Context, params SearchRecordsFieldsParams) (SearchRecordsFieldsRes, error)
	// SearchTable implements searchTable operation.
	//
	// Retrieves all fields of an file.
	//
	// GET /rest/tables/{table}/{fields}/{search}
	SearchTable(ctx context.Context, params SearchTableParams) (SearchTableRes, error)
	// SetConfig implements setConfig operation.
	//
	// Store configuration.
	//
	// PUT /config
	SetConfig(ctx context.Context, req SetConfigReq) (SetConfigRes, error)
	// SetJobsConfig implements setJobsConfig operation.
	//
	// Set the ADADATADIR.
	//
	// PUT /config/jobs
	SetJobsConfig(ctx context.Context, req OptJobStore) (SetJobsConfigRes, error)
	// ShutdownServer implements shutdownServer operation.
	//
	// Init shutdown procedure.
	//
	// PUT /rest/shutdown/{hash}
	ShutdownServer(ctx context.Context, params ShutdownServerParams) (ShutdownServerRes, error)
	// StoreConfig implements storeConfig operation.
	//
	// Store configuration.
	//
	// POST /config
	StoreConfig(ctx context.Context) (StoreConfigRes, error)
	// TriggerExtend implements triggerExtend operation.
	//
	// Put extend/plugin request.
	//
	// PUT /rest/extend/{path}
	TriggerExtend(ctx context.Context, params TriggerExtendParams) (TriggerExtendRes, error)
	// TriggerJob implements triggerJob operation.
	//
	// Trigger a job.
	//
	// PUT /rest/tasks/{jobName}
	TriggerJob(ctx context.Context, params TriggerJobParams) (TriggerJobRes, error)
	// UpdateLobByMap implements updateLobByMap operation.
	//
	// Set a lob at a specific ISN of an field in a Map.
	//
	// PUT /binary/{table}/{field}/{search}
	UpdateLobByMap(ctx context.Context, req UpdateLobByMapReq, params UpdateLobByMapParams) (UpdateLobByMapRes, error)
	// UpdateRecordsByFields implements updateRecordsByFields operation.
	//
	// Update a record dependent on field(s) of a specific table.
	//
	// PUT /rest/view/{table}/{search}
	UpdateRecordsByFields(ctx context.Context, req OptUpdateRecordsByFieldsReq, params UpdateRecordsByFieldsParams) (UpdateRecordsByFieldsRes, error)
	// UploadFile implements uploadFile operation.
	//
	// Upload a new file to the given location.
	//
	// POST /rest/file/{path}
	UploadFile(ctx context.Context, req *UploadFileReq, params UploadFileParams) (UploadFileRes, error)
	// NewError creates *ErrorStatusCode from error returned by handler.
	//
	// Used for common default response.
	NewError(ctx context.Context, err error) *ErrorStatusCode
}

// Server implements http server based on OpenAPI v3 specification and
// calls Handler to handle requests.
type Server struct {
	h   Handler
	sec SecurityHandler
	baseServer
}

// NewServer creates new Server.
func NewServer(h Handler, sec SecurityHandler, opts ...ServerOption) (*Server, error) {
	s, err := newServerConfig(opts...).baseServer()
	if err != nil {
		return nil, err
	}
	return &Server{
		h:          h,
		sec:        sec,
		baseServer: s,
	}, nil
}

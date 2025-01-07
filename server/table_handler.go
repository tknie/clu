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

// GetFields implements getFields operation.
//
// Retrieves all fields of an file.
//
// GET /rest/tables/{table}/fields
func (Handler) GetFields(ctx context.Context, params api.GetFieldsParams) (r api.GetFieldsRes, _ error) {
	return r, ht.ErrNotImplemented
}

// GetMapMetadata implements getMapMetadata operation.
//
// Retrieves metadata of a Map definition.
//
// GET /rest/metadata/view/{table}
func (Handler) GetMapMetadata(ctx context.Context, params api.GetMapMetadataParams) (r api.GetMapMetadataRes, _ error) {
	return r, ht.ErrNotImplemented
}

// InsertMapFileRecords implements insertMapFileRecords operation.
//
// Store send records into Map definition.
//
// POST /rest/view
func (Handler) InsertMapFileRecords(ctx context.Context, req api.OptInsertMapFileRecordsReq) (r api.InsertMapFileRecordsRes, _ error) {
	return r, ht.ErrNotImplemented
}

// ListTables implements listTables operation.
//
// Retrieves all tables of databases.
//
// GET /rest/tables
func (Handler) ListTables(ctx context.Context) (r api.ListTablesRes, _ error) {
	return r, ht.ErrNotImplemented
}

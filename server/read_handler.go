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
	"strings"

	ht "github.com/ogen-go/ogen/http"
	"github.com/tknie/clu"
	"github.com/tknie/clu/api"
	"github.com/tknie/flynn/common"
	"github.com/tknie/log"
	"github.com/tknie/services/auth"
)

// Handler server handler to ogen API
type Handler struct {
}

// SearchRecordsFields implements searchRecordsFields operation.
//
// Retrieves a field of a specific ISN of a Map definition.
//
// GET /rest/view/{table}/{fields}
func (Handler) SearchRecordsFields(ctx context.Context, params api.SearchRecordsFieldsParams) (r api.SearchRecordsFieldsRes, _ error) {
	session := ctx.(*clu.Context)
	log.Log.Debugf("Search records for fields %s -> %s", session.User, params.Table)
	if !auth.ValidUser(auth.UserRole, false, session.User, params.Table) {
		return &api.SearchRecordsFieldsForbidden{}, nil
	}
	log.Log.Debugf("SQL search fields %s - %v", params.Table, params.Search)
	d, err := ConnectTable(session, params.Table)
	if err != nil {
		log.Log.Errorf("Error search table %s:%v", params.Table, err)
		return nil, err // NewAPIError(err), nil
	}
	defer CloseTable(d)

	descriptor := false
	if params.Descriptor.Value {
		descriptor = true
	}
	limit := uint32(0)
	if params.Limit.Set {
		limit = uint32(params.Limit.Value)
	}
	q := &common.Query{TableName: params.Table,
		Fields:     extractFieldList(params.Search),
		Search:     "",
		Limit:      limit,
		Descriptor: descriptor,
		Order:      checkOrderBy(params.Orderby)}
	data, err := query(d, q)
	if err != nil {
		log.Log.Errorf("Error during query on %s:%v", params.Table, err)
		return nil, err
	}

	resp := api.Response{Records: data}
	respH := &api.ResponseHeaders{Response: resp, XToken: api.NewOptString(session.Token)}
	return respH, nil
}

// GetMapRecordsFields implements getMapRecordsFields operation.
//
// Retrieves a field of a specific ISN of a Map definition.
//
// GET /rest/view/{table}/{fields}/{search}
func (Handler) GetMapRecordsFields(ctx context.Context, params api.GetMapRecordsFieldsParams) (r api.GetMapRecordsFieldsRes, _ error) {
	session := ctx.(*clu.Context)
	if !auth.ValidUser(auth.UserRole, false, session.User, params.Table) {
		return &api.GetMapRecordsFieldsForbidden{}, nil
	}
	log.Log.Debugf("SQL search %s - %v", params.Table, params.Fields)
	d, err := ConnectTable(session, params.Table)
	if err != nil {
		log.Log.Errorf("Error search table %s:%v", params.Table, err)
		return nil, err
	}
	defer CloseTable(d)

	descriptor := false
	if params.Descriptor.Set && params.Descriptor.Value {
		descriptor = true
	}
	limit := uint32(0)
	if params.Limit.Set {
		limit = uint32(params.Limit.Value)
	}
	q := &common.Query{TableName: params.Table,
		Fields:     extractFieldList(params.Fields),
		Search:     params.Search,
		Descriptor: descriptor,
		Limit:      limit,
		Order:      checkOrderBy(params.Orderby)}
	data, err := query(d, q)
	if err != nil {
		log.Log.Errorf("Error during query on %s:%v", params.Table, err)
		return nil, err
	}

	log.Log.Debugf("Return payload")
	resp := api.Response{Records: data}
	respH := &api.ResponseHeaders{Response: resp, XToken: api.NewOptString(session.Token)}
	return respH, nil
}

func extractFieldList(fields string) []string {
	f := []string{"*"}
	if fields != "" {
		f = strings.Split(fields, ",")
	}
	return f
}

func checkOrderBy(orderby api.OptString) []string {
	if !orderby.Set || orderby.Value == "" {
		return make([]string, 0)
	}
	orderbyList := strings.Split(orderby.Value, ",")
	return orderbyList
}

// SearchTable implements searchTable operation.
//
// Retrieves all fields of an file.
//
// GET /rest/tables/{table}/{fields}/{search}
func (Handler) SearchTable(ctx context.Context, params api.SearchTableParams) (r api.SearchTableRes, _ error) {

	return r, ht.ErrNotImplemented
}

// SearchModelling implements searchModelling operation.
//
// Retrieves all fields of an file.
//
// GET /rest/map/{path}
func (Handler) SearchModelling(ctx context.Context, params api.SearchModellingParams) (r api.SearchModellingRes, _ error) {
	session := ctx.(*clu.Context)
	if !auth.ValidUser(auth.UserRole, false, session.User, params.Path) {
		return &api.SearchModellingForbidden{}, nil
	}
	log.Log.Debugf("SQL search field of an table %s - %v", params.Path, params.Path)
	d, err := ConnectTable(session, params.Path)
	if err != nil {
		log.Log.Errorf("Error search table %s:%v", params.Path, err)
		return nil, err
	}
	defer CloseTable(d)

	fields, err := d.GetTableColumn(params.Path)
	if err != nil {
		return nil, err
	}
	res := &api.Response{MapName: api.NewOptString(params.Path), FieldNames: fields}
	return res, nil
}

// ListModelling implements listModelling operation.
//
// Retrieves all fields of an file.
//
// GET /rest/map
func (Handler) ListModelling(ctx context.Context) (r api.ListModellingRes, _ error) {
	return r, ht.ErrNotImplemented
}

// BatchQuery implements batchQuery operation.
//
// Call a SQL query batch command posted in body.
//
// POST /rest/batch
func (Handler) BatchQuery(ctx context.Context, req api.OptSQLQuery) (r api.BatchQueryRes, _ error) {
	return r, ht.ErrNotImplemented
}

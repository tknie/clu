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
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	ht "github.com/ogen-go/ogen/http"
	"github.com/tknie/clu"
	"github.com/tknie/clu/api"
	"github.com/tknie/errorrepo"
	"github.com/tknie/flynn/common"
	"github.com/tknie/log"
	"github.com/tknie/services/auth"
)

// TimeFormat time format for date time representation
const TimeFormat = "2006-01-02 15:04:05"

var csvDelimiter = ","

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
	user := session.User()
	log.Log.Debugf("Search records for fields %s -> %s", user.User, params.Table)
	if !Validate(session, auth.UserRole, params.Table) {
		return &api.SearchRecordsFieldsForbidden{}, nil
	}
	log.Log.Debugf("SQL search fields %s - %v", params.Table, params.Search)
	d, err := ConnectTable(session, params.Table)
	if err != nil {
		log.Log.Errorf("Error search table %s:%v", params.Table, err)
		return nil, err // NewAPIError(err), nil
	}

	descriptor := false
	if params.Descriptor.Value {
		descriptor = true
	}
	limit := "ALL"
	if params.Limit.Set {
		limit = params.Limit.Value
		if params.Limit.Value == "-1" {
			limit = "ALL"
		}
	}
	q := &common.Query{TableName: params.Table,
		Fields:     extractFieldList(params.Search),
		Search:     "",
		Limit:      limit,
		Descriptor: descriptor,
		Order:      checkOrderBy(params.Orderby)}
	req := session.CurrentRequest
	accept := req.Header.Get("Accept")
	if accept == "text/csv" {
		piper, pipew := io.Pipe()
		hs := api.SearchRecordsFieldsOKTextCsv{Data: piper}
		s := &api.SearchRecordsFieldsOKTextCsvHeaders{Response: hs}
		go parallelQuery(d, q, piper, pipew)

		if err != nil {
			return nil, err
		}
		return s, nil
	}
	data, fields, err := query(d, q)
	if err != nil {
		log.Log.Errorf("Error during query on %s:%v", params.Table, err)
		return nil, err
	}

	if len(fields) == 0 {
		fields = extractFieldList(params.Search)
	}
	resp := api.Response{Records: data, FieldNames: fields}
	respH := &api.ResponseHeaders{Response: resp, XToken: api.NewOptString(session.Token)}
	log.Log.Debugf("DONE SQL search fields %s - %v", params.Table, params.Search)
	log.Log.Debugf("DONE SQL result %#v", respH)
	return respH, nil

}

func parallelQuery(d common.RegDbID, q *common.Query, piper *io.PipeReader, pipew *io.PipeWriter) {
	defer CloseTable(d)

	_, err := d.Query(q, func(search *common.Query, result *common.Result) error {
		if result == nil {
			return errorrepo.NewError("REST00011")
		}
		if result.Counter == 1 {
			x := strings.Join(result.Fields, csvDelimiter)
			pipew.Write([]byte(x + "\n"))
		}
		str := ""
		for i, x := range result.Rows {
			if i > 0 {
				str += csvDelimiter
			}
			switch t := x.(type) {
			case pgtype.Numeric:
				f, err := t.Float64Value()
				if err == nil {
					str += fmt.Sprintf("%v", f.Float64)
				}
			case time.Time:
				str += t.Format(TimeFormat)
			case nil:
			default:
				log.Log.Debugf("Default CSV %T\n", t)
				str += fmt.Sprintf("%v", t)
			}
		}
		pipew.Write([]byte(str + "\n"))
		return nil
	})
	if err != nil {
		pipew.Write([]byte(fmt.Sprintf("Error: %v", err)))
	}
	//	piper.Close()
	pipew.Close()
}

// GetMapRecordsFields implements getMapRecordsFields operation.
//
// Retrieves a field of a specific ISN of a Map definition.
//
// GET /rest/view/{table}/{fields}/{search}
func (Handler) GetMapRecordsFields(ctx context.Context, params api.GetMapRecordsFieldsParams) (r api.GetMapRecordsFieldsRes, _ error) {
	session := ctx.(*clu.Context)
	if !Validate(session, auth.UserRole, params.Table) {
		return &api.GetMapRecordsFieldsForbidden{}, nil
	}
	log.Log.Debugf("SQL search %s - %v -> %s", params.Table, params.Fields, params.Search)
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
	limit := "ALL"
	if params.Limit.Set {
		limit = params.Limit.Value
	}
	q := &common.Query{TableName: params.Table,
		Fields:     extractFieldList(params.Fields),
		Search:     params.Search,
		Descriptor: descriptor,
		Limit:      limit,
		Order:      checkOrderBy(params.Orderby)}
	req := session.CurrentRequest
	accept := req.Header.Get("Accept")
	if accept == "text/csv" {
		piper, pipew := io.Pipe()
		hs := api.GetMapRecordsFieldsOKTextCsv{Data: piper}
		s := &api.GetMapRecordsFieldsOKTextCsvHeaders{Response: hs}
		go parallelQuery(d, q, piper, pipew)

		if err != nil {
			return nil, err
		}
		return s, nil
	}
	data, fields, err := query(d, q)
	if err != nil {
		log.Log.Errorf("Error during query on %s:%v", params.Table, err)
		return nil, err
	}
	log.Log.Debugf("Return SQL search payload %d entries", len(data))

	resp := api.Response{Records: data, FieldNames: fields,
		MapName:   api.NewOptString(params.Table),
		NrRecords: api.NewOptInt(len(data))}
	respH := &api.ResponseHeaders{Response: resp, XToken: api.NewOptString(session.Token)}
	log.Log.Debugf("Return SQL search %s - %v -> %s", params.Table, params.Fields, params.Search)
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
	if !Validate(session, auth.UserRole, params.Path) {
		return &api.SearchModellingForbidden{}, nil
	}
	log.Log.Debugf("SQL modelling field of an table %s - %v", params.Path, params.Path)
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
	log.Log.Debugf("Return SQL modelling field of an table %s - %v", params.Path, params.Path)
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

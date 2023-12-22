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
	"fmt"
	"io"

	"github.com/go-faster/jx"
	"github.com/tknie/clu"
	"github.com/tknie/clu/api"
	"github.com/tknie/flynn/common"
	"github.com/tknie/log"
	"github.com/tknie/services"
	"github.com/tknie/services/auth"
)

// BatchQuery implements batchQuery operation.
//
// Call a SQL query batch command posted in body.
//
// POST /rest/batch/{table}
func (Handler) BatchQuery(ctx context.Context, req api.BatchQueryReq,
	params api.BatchQueryParams) (r api.BatchQueryRes, _ error) {
	sqlStatement := ""
	switch sqlQuery := req.(type) {
	case *api.SQLQuery:
		sqlStatement = sqlQuery.Batch.Value.SQL.Value
	case *api.BatchQueryReqTextPlain:
		b := make([]byte, 1024)
		for k, _ := range b {
			b[k] = ' '
		}
		n, err := sqlQuery.Read(b)
		if err != io.EOF && err != nil {
			return nil, err
		}
		log.Log.Debugf("Receive buffer %d", n)
		sqlStatement = string(b[:n])
	}
	session := ctx.(*clu.Context)
	if !auth.ValidUser(auth.UserRole, false, session.User, "/batch") {
		log.Log.Debugf("SQL statemant forbidden")
		return &api.BatchQueryForbidden{}, nil
	}
	log.Log.Debugf("SQL statement on table %s - %v", params.Table, sqlStatement)
	services.ServerMessage("SQL query by user %s: %s", session.User.User, sqlStatement)

	d, err := ConnectTable(session, params.Table)
	if err != nil {
		log.Log.Errorf("Error search table %s:%v", params.Table, err)
		return nil, err
	}
	defer CloseTable(d)
	batch := &common.Query{Search: sqlStatement}

	rria := make([]api.ResponseRecordsItem, 0)
	var fields []string
	err = d.BatchSelectFct(batch, func(search *common.Query, result *common.Result) error {
		if fields == nil {
			fields = result.Fields
		}
		rri := api.ResponseRecordsItem{}
		for d, r := range result.Rows {
			rri[result.Fields[d]] = jx.Raw(fmt.Sprintf("%v", r))
		}
		rria = append(rria, rri)
		return nil
	})
	if err != nil {
		return nil, err
	}
	resp := api.Response{NrRecords: api.NewOptInt(int(len(rria))),
		FieldNames: fields,
		MapName:    api.NewOptString(params.Table), Records: rria}
	respH := &api.ResponseHeaders{Response: resp, XToken: api.NewOptString(session.Token)}
	return respH, nil
}

// BatchParameterQuery implements batchParameterQuery operation.
//
// Call a SQL query batch command posted in query.
//
// GET /rest/batch/{table}/{query}
func (Handler) BatchParameterQuery(ctx context.Context, params api.BatchParameterQueryParams) (r api.BatchParameterQueryRes, _ error) {
	session := ctx.(*clu.Context)
	if !auth.ValidUser(auth.UserRole, false, session.User, "/batch") {
		log.Log.Debugf("SQL statemant forbidden")
		return &api.BatchParameterQueryForbidden{}, nil
	}
	log.Log.Debugf("SQL statement on table %s - %v", params.Table, params.Query)
	services.ServerMessage("SQL query by user %s: %s", session.User.User, params.Query)

	d, err := ConnectTable(session, params.Table)
	if err != nil {
		log.Log.Errorf("Error search table %s:%v", params.Table, err)
		return nil, err
	}
	defer CloseTable(d)
	batch := &common.Query{Search: params.Query}

	rria := make([]api.ResponseRecordsItem, 0)
	var fields []string
	err = d.BatchSelectFct(batch, func(search *common.Query, result *common.Result) error {
		if fields == nil {
			fields = result.Fields
		}
		rri := api.ResponseRecordsItem{}
		for d, r := range result.Rows {
			rri[result.Fields[d]] = jx.Raw(fmt.Sprintf("%v", r))
		}
		rria = append(rria, rri)
		return nil
	})
	if err != nil {
		return nil, err
	}
	resp := api.Response{NrRecords: api.NewOptInt(int(len(rria))),
		FieldNames: fields,
		MapName:    api.NewOptString(params.Table), Records: rria}
	respH := &api.ResponseHeaders{Response: resp, XToken: api.NewOptString(session.Token)}
	return respH, nil
}

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
	"io"
	"strings"

	"github.com/tknie/clu"
	"github.com/tknie/clu/api"
	"github.com/tknie/flynn/common"
	"github.com/tknie/log"
	"github.com/tknie/services/auth"
)

type batchSelect struct {
	session   *clu.Context
	table     string
	parameter []string
	query     *clu.BatchEntry
}

// BatchSelect implements batchSelect operation.
//
// Call a SQL query batch command out of the stored query list.
//
// GET /rest/batch/{table}
func (Handler) BatchSelect(ctx context.Context,
	params api.BatchSelectParams) (r api.BatchSelectRes, _ error) {
	session := ctx.(*clu.Context)

	if !Validate(session, auth.UserRole, "^"+params.Table) {
		log.Log.Debugf("Validator forbidden")
		return &api.BatchSelectForbidden{}, nil
	}
	log.Log.Debugf("Query batchname %s -> %#v", params.Table, params.Param)
	entry, err := clu.BatchSelect(params.Table)
	if err != nil {
		return nil, err
	}
	if entry == nil {
		log.Log.Fatal("Query entry empty")
	}
	respH, err := querySQLstatement(&batchSelect{session: session, table: params.Table,
		parameter: params.Param, query: entry})
	if err != nil {
		return nil, err
	}
	return respH, nil
}

// BatchQuery implements batchQuery operation.
//
// Call a SQL query batch command posted in body.
//
// POST /rest/batch/{table}
func (Handler) BatchQuery(ctx context.Context, req api.BatchQueryReq,
	params api.BatchQueryParams) (r api.BatchQueryRes, _ error) {
	sqlStatement := ""
	var p []string
	switch sqlQuery := req.(type) {
	case *api.SQLQuery:
		sqlStatement = sqlQuery.Batch.Value.SQL.Value
	case *api.BatchQueryReqTextPlain:
		b := make([]byte, 1024)
		for k := range b {
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
	if !Validate(session, auth.UserRole, "^"+params.Table) {
		log.Log.Debugf("SQL statemant forbidden")
		return &api.BatchQueryForbidden{}, nil
	}
	log.Log.Debugf("SQL statement on table %s - %v", params.Table, sqlStatement)
	// services.ServerMessage("SQL query by user %s: %s", session.User.User, sqlStatement)

	respH, err := querySQLstatement(&batchSelect{session: session, table: params.Table,
		parameter: p,
		query:     &clu.BatchEntry{Query: sqlStatement, Database: params.Table}})
	if err != nil {
		return nil, err
	}
	return respH, nil

}

func querySQLstatement(query *batchSelect) (*api.ResponseHeaders, error) {
	log.Log.Debugf("Query/Batch SQL statemant %s: %#v", query.query.Query, query.parameter)
	d, err := ConnectTable(query.session, query.query.Database)
	if err != nil {
		log.Log.Errorf("Error search table %s:%v", query.query.Database, err)
		return nil, err
	}
	defer CloseTable(d)
	batch := &common.Query{Search: sqlInParameter(query.query.Query, query.parameter)}

	rria := make([]api.ResponseRecordsItem, 0)
	var fields []string
	err = d.BatchSelectFct(batch, func(search *common.Query, result *common.Result) error {
		if fields == nil {
			fields = result.Fields
		}
		rri := generateItem(result.Fields, result.Rows)
		rria = append(rria, rri)
		return nil
	})
	if err != nil {
		return nil, err
	}
	resp := api.Response{NrRecords: api.NewOptInt(int(len(rria))),
		FieldNames: fields,
		MapName:    api.NewOptString(query.table), Records: rria}
	respH := &api.ResponseHeaders{Response: resp, XToken: api.NewOptString(query.session.Token)}
	log.Log.Debugf("Return Query/Batch SQL statemant %s: %#v", query.query.Query, query.parameter)
	return respH, nil
}

// BatchParameterQuery implements batchParameterQuery operation.
//
// Call a SQL query batch command posted in query.
//
// GET /rest/batch/{table}/{query}
func (Handler) BatchParameterQuery(ctx context.Context, params api.BatchParameterQueryParams) (r api.BatchParameterQueryRes, _ error) {
	session := ctx.(*clu.Context)
	if !Validate(session, auth.UserRole, "^"+params.Table) {
		log.Log.Debugf("Batch SQL statement for user forbidden, returning forbidden‚")
		return &api.BatchParameterQueryForbidden{}, nil
	}
	log.Log.Debugf("SQL statement on table %s - %v", params.Table, params.Query)
	// services.ServerMessage("SQL query by user %s: %s", session.User.User, params.Query)

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
		rri := generateItem(result.Fields, result.Rows)
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
	log.Log.Debugf("Return SQL statement on table %s - %v", params.Table, params.Query)
	return respH, nil
}

func sqlInParameter(statement string, params []string) string {
	st := statement
	for _, p := range params {
		log.Log.Debugf("Given parameter '%s'", p)
		np := strings.Trim(p, "\"")
		if np[0] == '^' {
			pv := strings.Split(np, ":")
			if len(np) > len(pv[0]) {
				log.Log.Debugf("Handle parameter %s : %s", pv[0], np[len(pv[0])+1:])
				st = strings.Replace(st, "<"+pv[0][1:]+">", np[len(pv[0])+1:], -1)
			}
		}
	}
	log.Log.Debugf("SQL in : %s", statement)
	log.Log.Debugf("SQL out: %s", st)
	return st
}

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
	"strings"

	"github.com/go-faster/jx"
	"github.com/tknie/clu"
	"github.com/tknie/clu/api"
	"github.com/tknie/flynn/common"
	"github.com/tknie/log"
	"github.com/tknie/services/auth"
)

// InsertRecord implements insertRecord operation.
//
// Insert given record.
//
// POST /rest/view/{table}
func (ServerHandler) InsertRecord(ctx context.Context, req api.OptInsertRecordReq, params api.InsertRecordParams) (r api.InsertRecordRes, _ error) {
	session := ctx.(*clu.Context)
	log.Log.Debugf("Insert records for fields %s -> %s", session.User, params.Table)
	if !auth.ValidUser(auth.UserRole, false, session.User, params.Table) {
		return &api.InsertRecordForbidden{}, nil
	}
	log.Log.Debugf("SQL insert%s", params.Table)
	d, err := ConnectTable(session, params.Table)
	if err != nil {
		log.Log.Errorf("Error search table %s:%v", params.Table, err)
		return &api.Error{Error: api.NewOptErrorError(api.ErrorError{Message: api.NewOptString(err.Error())})}, nil
	}
	defer CloseTable(d)

	records := make([]any, 0)
	nameMap := make(map[string]bool)
	for _, r := range req.Value.Records {
		m := make(map[string]any)
		for n, v := range r {
			v, err := parseJx(v)
			if err != nil {
				log.Log.Debugf("Error %s: %v", n, err)
				return &api.Error{Error: api.NewOptErrorError(api.ErrorError{Message: api.NewOptString(err.Error())})}, nil
			}
			log.Log.Debugf("[%s]=%v", n, v)
			m[n] = v
			nameMap[n] = true
		}
		records = append(records, m)
	}
	fields := make([]string, 0)
	for n := range nameMap {
		fields = append(fields, n)
	}
	list := make([][]any, 0)
	for _, r := range records {
		subList := make([]any, 0)
		m := r.(map[string]any)
		for _, n := range fields {
			subList = append(subList, m[n])
		}
		list = append(list, subList)
	}
	// list := [][]any{{vId1, "xxxxxx", 1}, {vId2, "yyywqwqwqw", 2}}
	input := &common.Entries{Fields: fields,
		Values: list}
	fmt.Printf("%#v ->>>\n", input)
	err = d.Insert(params.Table, input)
	if err != nil {
		log.Log.Debugf("Error: %v", err)
		return &api.Error{Error: api.NewOptErrorError(api.ErrorError{Message: api.NewOptString(err.Error())})}, nil
	}
	fmt.Println("INSERT:", records)
	resp := api.Response{NrRecords: api.NewOptInt(1)}
	respH := &api.ResponseHeaders{Response: resp, XToken: api.NewOptString(session.Token)}
	return respH, nil
}

func parseJx(v jx.Raw) (any, error) {
	d := jx.DecodeBytes(v)
	switch v.Type() {
	case jx.Number:
		x, err := d.Int()
		return x, err
	case jx.String:
		x, err := d.Str()
		return x, err
	case jx.Array:
		values := make([]any, 0)
		d.Arr(func(d *jx.Decoder) error {
			o, err := parseJx(v)
			if err != nil {
				return err
			}
			values = append(values, o)
			return nil
		})
		return values, nil
	case jx.Object:
		ms := make(map[string]any)
		d.Obj(func(d *jx.Decoder, key string) error {
			r, err := d.Raw()
			if err != nil {
				return err
			}
			v, err := parseJx(r)
			if err != nil {
				return err
			}
			ms[key] = v
			return nil
		})
		return ms, nil
	default:
		fmt.Println("->>", v.Type().String())
	}
	return nil, fmt.Errorf("json type unknown")
}

// DeleteRecordsSearched implements deleteRecordsSearched operation.
//
// Delete a record with a given search.
//
// DELETE /rest/view/{table}/{search}
func (ServerHandler) DeleteRecordsSearched(ctx context.Context, params api.DeleteRecordsSearchedParams) (r api.DeleteRecordsSearchedRes, _ error) {
	session := ctx.(*clu.Context)
	log.Log.Debugf("Search records for fields %s -> %s", session.User, params.Table)
	if !auth.ValidUser(auth.UserRole, false, session.User, params.Table) {
		return &api.DeleteRecordsSearchedForbidden{}, nil
	}
	log.Log.Debugf("SQL search fields %s - %v", params.Table, params.Search)
	d, err := ConnectTable(session, params.Table)
	if err != nil {
		log.Log.Errorf("Error search table %s:%v", params.Table, err)
		return &api.Error{Error: api.NewOptErrorError(api.ErrorError{Message: api.NewOptString(err.Error())})}, nil
	}
	defer CloseTable(d)
	dr, err := d.Delete(params.Table, &common.Entries{Criteria: params.Search})
	if err != nil {
		log.Log.Errorf("Error delete search %s->%s:%v", params.Table, params.Search, err)
		return &api.Error{Error: api.NewOptErrorError(api.ErrorError{Message: api.NewOptString(err.Error())})}, nil
	}
	fmt.Println("DR", dr)
	resp := api.Response{NrRecords: api.NewOptInt(int(dr))}
	respH := &api.ResponseHeaders{Response: resp, XToken: api.NewOptString(session.Token)}
	return respH, nil
}

// UpdateRecordsByFields implements updateRecordsByFields operation.
//
// Update a record dependent on field(s) of a specific table.
//
// PUT /rest/view/{table}/{search}
func (ServerHandler) UpdateRecordsByFields(ctx context.Context, req api.OptUpdateRecordsByFieldsReq,
	params api.UpdateRecordsByFieldsParams) (r api.UpdateRecordsByFieldsRes, _ error) {
	session := ctx.(*clu.Context)
	log.Log.Debugf("Update records for fields %s -> %s", session.User, params.Table)
	if !auth.ValidUser(auth.UserRole, false, session.User, params.Table) {
		return &api.UpdateRecordsByFieldsForbidden{}, nil
	}
	log.Log.Debugf("SQL update %s", params.Table)
	d, err := ConnectTable(session, params.Table)
	if err != nil {
		log.Log.Errorf("Error update table %s:%v", params.Table, err)
		return &api.Error{Error: api.NewOptErrorError(api.ErrorError{Message: api.NewOptString(err.Error())})}, nil
	}
	defer CloseTable(d)

	records := make([]any, 0)
	nameMap := make(map[string]bool)
	for _, r := range req.Value.Records {
		m := make(map[string]any)
		for n, v := range r {
			v, err := parseJx(v)
			if err != nil {
				log.Log.Debugf("Error %s: %v", n, err)
				return &api.Error{Error: api.NewOptErrorError(api.ErrorError{Message: api.NewOptString(err.Error())})}, nil
			}
			log.Log.Debugf("[%s]=%v", n, v)
			m[n] = v
			nameMap[n] = true
		}
		records = append(records, m)
	}
	fields := make([]string, 0)
	for n := range nameMap {
		fields = append(fields, n)
	}
	list := make([][]any, 0)
	for _, r := range records {
		subList := make([]any, 0)
		m := r.(map[string]any)
		for _, n := range fields {
			subList = append(subList, m[n])
		}
		list = append(list, subList)
	}
	updateFields := strings.Split(params.Search, ",")
	input := &common.Entries{Fields: fields,
		Update: updateFields,
		Values: list}
	fmt.Printf("%#v ->>>\n", input)
	uNr, err := d.Update(params.Table, input)
	if err != nil {
		log.Log.Debugf("Error: %v", err)
		return &api.Error{Error: api.NewOptErrorError(api.ErrorError{Message: api.NewOptString(err.Error())})}, nil
	}
	fmt.Println("Update:", records)
	resp := api.Response{NrRecords: api.NewOptInt(int(uNr))}
	respH := &api.ResponseHeaders{Response: resp, XToken: api.NewOptString(session.Token)}
	return respH, nil
}

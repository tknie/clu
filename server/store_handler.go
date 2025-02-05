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
	"math"
	"strings"
	"time"

	"github.com/go-faster/jx"
	ht "github.com/ogen-go/ogen/http"
	"github.com/tknie/clu"
	"github.com/tknie/clu/api"
	"github.com/tknie/errorrepo"
	"github.com/tknie/flynn/common"
	"github.com/tknie/log"
	"github.com/tknie/services/auth"
)

// InsertRecord implements insertRecord operation.
//
// Insert given record.
//
// POST /rest/view/{table}
func (Handler) InsertRecord(ctx context.Context, req api.OptInsertRecordReq, params api.InsertRecordParams) (r api.InsertRecordRes, _ error) {
	session := ctx.(*clu.Context)
	log.Log.Debugf("Insert records for fields %s -> %s", session.User, params.Table)
	if !auth.ValidUser(auth.UserRole, true, session.User(), params.Table) {
		return &api.InsertRecordForbidden{}, nil
	}
	log.Log.Debugf("SQL insert %s", params.Table)
	d, err := ConnectTable(session, params.Table)
	if err != nil {
		log.Log.Errorf("Error search table %s:%v", params.Table, err)
		return nil, err
	}
	defer CloseTable(d)

	log.Log.Debugf("Incoming %#v", req.Value)
	if req.Value.Records == nil {
		return &api.InsertRecordBadRequest{}, nil
	}

	records := make([]any, 0)
	nameMap := make(map[string]bool)
	for _, r := range req.Value.Records {
		m := make(map[string]any)
		for n, v := range r {
			v, err := parseJx(v)
			if err != nil {
				log.Log.Debugf("Error JSON parser %s: %v", n, err)
				return nil, errorrepo.NewError("RERR00015", n, err)
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
	if params.Returning.Set {
		input.Returning = strings.Split(params.Returning.Value, ",")
	}
	retValue, err := d.Insert(params.Table, input)
	if err != nil {
		log.Log.Debugf("Error: %v", err)
		return nil, err
	}

	resp := api.Response{NrRecords: api.NewOptInt(len(records))}
	if len(input.Returning) > 0 {
		log.Log.Debugf("Returning value: %v", retValue)
		data := make([]api.ResponseRecordsItem, 0)
		for _, r := range retValue {
			d := make(api.ResponseRecordsItem)
			for x, field := range input.Returning {
				convertTypeToRaw(d, field, r[x])
			}
			data = append(data, d)
		}
		resp.Records = data
	}
	respH := &api.ResponseHeaders{Response: resp, XToken: api.NewOptString(session.Token)}
	log.Log.Debugf("Return SQL insert %s", params.Table)
	return respH, nil
}

func parseJx(v jx.Raw) (any, error) {
	d := jx.DecodeBytes(v)
	switch v.Type() {
	case jx.Number:
		f, err := d.Float64()
		if f == math.Trunc(f) {
			return int(f), nil
		}
		// x, err := d.Int()
		return f, err
	case jx.String:
		x, err := d.Str()
		if err != nil {
			return nil, err
		}
		t, err := time.Parse(time.RFC3339, x)
		if err == nil {
			return t, err
		}
		return x, nil
	case jx.Bool:
		x, err := d.Bool()
		return x, err
	case jx.Array:
		values := make([]any, 0)
		d.Arr(func(d *jx.Decoder) error {
			o, err := parseJx(v)
			if err != nil {
				log.Log.Debugf("Error JSON parser of array: %v", err)
				return errorrepo.NewError("RERR00015", "<Array>", err)
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
				log.Log.Debugf("Error JSON parser %s: %v", key, err)
				return errorrepo.NewError("RERR00015", key, err)
			}
			ms[key] = v
			return nil
		})
		return ms, nil
	case jx.Null:
		return nil, nil
	default:
		fmt.Println("Unknown type ->>", v.Type().String())
	}
	return nil, errorrepo.NewError("REST00050")
}

// DeleteRecordsSearched implements deleteRecordsSearched operation.
//
// Delete a record with a given search.
//
// DELETE /rest/view/{table}/{search}
func (Handler) DeleteRecordsSearched(ctx context.Context, params api.DeleteRecordsSearchedParams) (r api.DeleteRecordsSearchedRes, _ error) {
	session := ctx.(*clu.Context)
	log.Log.Debugf("Delete records for fields %s -> %s", session.User, params.Table)
	if !auth.ValidUser(auth.UserRole, true, session.User(), params.Table) {
		return &api.DeleteRecordsSearchedForbidden{}, nil
	}
	log.Log.Debugf("SQL search fields %s - %v", params.Table, params.Search)
	d, err := ConnectTable(session, params.Table)
	if err != nil {
		log.Log.Errorf("Error search table %s:%v", params.Table, err)
		return nil, err
	}
	defer CloseTable(d)
	dr, err := d.Delete(params.Table, &common.Entries{Criteria: params.Search})
	if err != nil {
		log.Log.Errorf("Error delete search %s->%s:%v", params.Table, params.Search, err)
		return nil, err
	}
	log.Log.Errorf("%d Data record deleted from %s: %s", dr, params.Table, params.Search)
	resp := api.Response{NrRecords: api.NewOptInt(int(dr))}
	respH := &api.ResponseHeaders{Response: resp, XToken: api.NewOptString(session.Token)}
	log.Log.Debugf("Delete SQL search fields %s - %v", params.Table, params.Search)
	return respH, nil
}

// UpdateRecordsByFields implements updateRecordsByFields operation.
//
// Update a record dependent on field(s) of a specific table.
//
// PUT /rest/view/{table}/{search}
func (Handler) UpdateRecordsByFields(ctx context.Context, req api.OptUpdateRecordsByFieldsReq,
	params api.UpdateRecordsByFieldsParams) (r api.UpdateRecordsByFieldsRes, _ error) {
	session := ctx.(*clu.Context)
	log.Log.Debugf("Update records for fields %s -> %s", session.User, params.Table)
	if !auth.ValidUser(auth.UserRole, true, session.User(), params.Table) {
		return &api.UpdateRecordsByFieldsForbidden{}, nil
	}
	log.Log.Debugf("SQL update %s", params.Table)
	d, err := ConnectTable(session, params.Table)
	if err != nil {
		log.Log.Errorf("Error update table %s:%v", params.Table, err)
		return nil, err
	}
	defer CloseTable(d)

	records := make([]any, 0)
	nameMap := make(map[string]bool)
	for _, r := range req.Value.Records {
		m := make(map[string]any)
		for n, v := range r {
			v, err := parseJx(v)
			if err != nil {
				log.Log.Debugf("Error JSON parser %s: %v", n, err)
				return nil, errorrepo.NewError("RERR00015", n, err)
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
	_, uNr, err := d.Update(params.Table, input)
	if err != nil {
		log.Log.Debugf("Error: %v", err)
		return nil, err
	}
	resp := api.Response{NrRecords: api.NewOptInt(int(uNr))}
	respH := &api.ResponseHeaders{Response: resp, XToken: api.NewOptString(session.Token)}
	log.Log.Debugf("Return Update records for fields %s -> %s", session.User, params.Table)
	return respH, nil
}

// UpdateLobByMap implements updateLobByMap operation.
//
// Set a lob at a specific ISN of an field in a Map.
//
// PUT /binary/{table}/{field}/{search}
func (Handler) UpdateLobByMap(ctx context.Context, req api.UpdateLobByMapReq, params api.UpdateLobByMapParams) (r api.UpdateLobByMapRes, _ error) {
	return r, ht.ErrNotImplemented
}

/*
* Copyright 2022-2024 Thorsten A. Knieling
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
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/go-faster/jx"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/tknie/clu/api"
	"github.com/tknie/errorrepo"
	"github.com/tknie/flynn/common"
	"github.com/tknie/log"
)

// query query SQL tables
func query(d common.RegDbID, query *common.Query) ([]api.ResponseRecordsItem, []string, error) {
	log.Log.Debugf("Query in db ID %04d", d)
	data := make([]api.ResponseRecordsItem, 0)
	var fields []string
	_, err := d.Query(query, func(search *common.Query, result *common.Result) error {
		if result == nil {
			return errorrepo.NewError("REST00011")
		}
		if fields == nil {
			fields = result.Fields
		}
		log.Log.Debugf("Rows: %d", len(result.Rows))
		d := generateItem(result.Fields, result.Rows)
		data = append(data, d)
		return nil
	})
	if err != nil {
		log.Log.Debugf("SQL query error: %v", err)
		return nil, nil, err
	}
	return data, fields, nil
}

func generateItem(fields []string, rows []any) api.ResponseRecordsItem {
	///var d api.ResponseRecordsItem
	d := make(api.ResponseRecordsItem)
	for i, r := range rows {
		s := strings.ToLower(fields[i])
		log.Log.Debugf("%d. row is of type %T", i, r)
		convertTypeToRaw(d, s, r)
	}
	return d
}

func convertTypeToRaw(d api.ResponseRecordsItem, s string, r interface{}) {
	switch t := r.(type) {
	case *string:
		log.Log.Debugf("String Pointer %s", *t)
		raw := jx.Raw([]byte("\"" + *t + "\""))
		d[s] = raw
	case string:
		log.Log.Debugf("String %s", t)
		str := strings.Replace(t, "\"", "\\\"", -1)
		raw := jx.Raw([]byte("\"" + str + "\""))
		d[s] = raw
	case *time.Time:
		d[s] = jx.Raw([]byte("\"" + (*t).UTC().Format(time.RFC3339) + "\""))
	case time.Time:
		d[s] = jx.Raw([]byte("\"" + (t).UTC().Format(time.RFC3339) + "\""))
	case pgtype.Numeric:
		v := uint64(t.Int.Uint64()) * uint64(math.Pow10(int(t.Exp)))
		st := fmt.Sprintf("%d", v)
		d[s] = jx.Raw([]byte("\"" + (st) + "\""))
	default:
		if r != nil {
			log.Log.Debugf("using default ---> %v %T", r, t)
			d[s] = jx.Raw(fmt.Sprintf("%v", r))
		}
	}
}

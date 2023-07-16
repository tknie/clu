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
	"fmt"
	"strings"
	"time"

	"github.com/tknie/clu/api"
	"github.com/tknie/flynn/common"
	"github.com/tknie/log"
)

// query query SQL tables
func query(d common.RegDbID, query *common.Query) ([]api.ResponseRecordsItem, error) {
	log.Log.Debugf("Query in db ID %04d", d)
	var data []api.ResponseRecordsItem
	_, err := d.Query(query, func(search *common.Query, result *common.Result) error {
		if result == nil {
			return fmt.Errorf("result empty")
		}
		log.Log.Debugf("Rows: %d", len(result.Rows))
		///var d api.ResponseRecordsItem
		d := make(api.ResponseRecordsItem)
		for i, r := range result.Rows {
			s := strings.ToLower(result.Fields[i])
			log.Log.Debugf("%d. row is of type %T", i, r)
			switch t := r.(type) {
			case *string:
				log.Log.Debugf("String %s", *t)
				d[s] = *t
			case *time.Time:
				d[s] = (*t).String()
			default:
				d[s] = fmt.Sprintf("%v", r)
			}
		}
		data = append(data, d)
		return nil
	})
	if err != nil {
		return nil, err
	}
	return data, nil
}

// query query SQL tables
func queryBytes(d common.RegDbID, query *common.Query) (map[string]interface{}, error) {
	log.Log.Debugf("Query stream in db ID %04d", d)
	dataMap := make(map[string]interface{})
	found := false
	_, err := d.Query(query, func(search *common.Query, result *common.Result) error {
		if result == nil {
			return fmt.Errorf("result empty")
		}
		if found {
			return fmt.Errorf("result not unique")
		}
		log.Log.Debugf("Rows: %d", len(result.Rows))
		///var d api.ResponseRecordsItem
		for i, r := range result.Rows {
			s := strings.ToLower(result.Fields[i])
			log.Log.Debugf("%d. row is of type %T", i, r)
			switch t := r.(type) {
			case *string:
				log.Log.Debugf("String %s", *t)
				dataMap[s] = *t
			case *time.Time:
				dataMap[s] = *t
			default:
				dataMap[s] = r
			}
		}
		found = true
		return nil
	})
	if err != nil {
		return nil, err
	}
	return dataMap, nil
}

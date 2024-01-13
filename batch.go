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

package clu

import (
	"fmt"
	"os"
	"sync"

	"github.com/tknie/flynn"
	"github.com/tknie/flynn/common"
	"github.com/tknie/log"
	"github.com/tknie/services"
)

// BatchEntry entry of batch repository
type BatchEntry struct {
	Name       string
	Query      string `flynn:":BLOB"`
	Database   string
	ParamCount int
}

var batchDbRef *common.Reference
var batchDbPassword = ""

var batchtablename = ""
var batchStoreOnline = false
var batchLock sync.Mutex

type batchRepository struct {
	stored bool
}

func openBatchRepository() (common.RegDbID, error) {
	var err error
	if userDbPassword == "" {
		userDbPassword = os.Getenv("REST_BATCH_PASS")
	}
	sessionStoreID, err := flynn.Handler(batchDbRef, batchDbPassword)
	if err != nil {
		services.ServerMessage("Register error log: %v", err)
		return 0, err
	}
	return sessionStoreID, nil
}

// InitBatchRepository init batch repository
func InitBatchRepository(dbRef *common.Reference, dbPassword, tablename string) {
	batchDbRef = dbRef
	batchDbPassword = dbPassword
	batchStoreID, err := openBatchRepository()
	if err != nil {
		services.ServerMessage("Register error log: %v", err)
		return
	}
	log.Log.Debugf("Receive batch store handler %s", batchStoreID)
	defer batchStoreID.Close()

	dbTables := flynn.Maps()
	for _, d := range dbTables {
		if d == tablename {
			batchtablename = tablename
			batchStoreOnline = true
			services.ServerMessage("Using batch repository on table '%s'", batchtablename)
			return
		}
	}
	su := &BatchEntry{}
	err = batchStoreID.CreateTable(tablename, su)
	if err != nil {
		services.ServerMessage("Database batch store creating failed: %v", err)
		return
	}
	batchtablename = tablename
	batchStoreOnline = true
	services.ServerMessage("Database batch store '%s' created successfully", batchtablename)
}

// BatchSelect search for batchname in an batch repository
func BatchSelect(batchname string) (*BatchEntry, error) {
	if !batchStoreOnline {
		return nil, fmt.Errorf("error batch repository disabled")
	}
	batchLock.Lock()
	defer batchLock.Unlock()
	batchStoreID, err := openBatchRepository()
	if err != nil {
		services.ServerMessage("Register error log: %v", err)
		return nil, err
	}
	log.Log.Debugf("Receive batch store handler %s", batchStoreID)
	batchStoreID.Close()
	var b *BatchEntry
	q := &common.Query{TableName: batchtablename,
		Search:     "name='" + batchname + "'",
		DataStruct: &BatchEntry{},
		Fields:     []string{"*"}}
	_, err = batchStoreID.Query(q, func(search *common.Query, result *common.Result) error {
		if b == nil {
			b = result.Data.(*BatchEntry)
		}
		return nil
	})
	if err != nil {
		log.Log.Errorf("Query batch store failure: %v", err)
		return nil, err
	}
	return b, nil
}

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
	"strings"
	"time"

	"github.com/tknie/clu"
	"github.com/tknie/errorrepo"
	"github.com/tknie/flynn"
	"github.com/tknie/flynn/common"
	"github.com/tknie/log"
	"github.com/tknie/services"
)

// init config update tracker to tracke changes
func init() {
	clu.RegisterConfigUpdates(initRegister)
}

// initRegister initialize database register getting all current available database
func initRegister() {
	log.Log.Debugf("Register databases")
	initTableOfDatabases()
	log.Log.Debugf("Start table thread for databases")
	go loadTableThread()
}

func loadTableThread() {
	dateTicker := time.NewTicker(60 * time.Second)
	loadTableOfDatabases()
	for {
		<-dateTicker.C
		loadTableOfDatabases()
	}
}

func initTableOfDatabases() {
	for _, dm := range clu.Viewer.Database.DatabaseAccess.Database {
		(&dm).Handles()
	}
}

// loadTableOfDatabases load all tables of databases registered
func loadTableOfDatabases() {
	log.Log.Debugf("Refreshing database list")
	for _, dm := range clu.Viewer.Database.DatabaseAccess.Database {
		log.Log.Debugf("Access database %s with user %s", dm.Target, dm.User)
		id, err := (&dm).Handles()
		if err != nil {
			log.Log.Debugf("Handle creation problem: %v", err)
			continue
		}

		newDatabases := make([]string, 0)
		for _, table := range flynn.Maps() {
			s := strings.ToLower(table)
			if clu.CheckDatabaseRegister(s, id) {
				services.ServerMessage("Found table on different databases: [%s]", s)
			} else {
				if (&dm).RegisterDatabase(s, id) {
					log.Log.Infof("Register database '%s' -> %#v", dm.Table, dm)
					newDatabases = append(newDatabases, s)
				}
			}
		}
		if len(newDatabases) > 0 {
			services.ServerMessage("Collected %04d table(s) in dictionary", len(newDatabases))
		}
	}
	clu.DumpStat()

}

// InitDatabases initialize database reference IDs
func InitDatabases() {
	log.Log.Debugf("Init databases done")
}

// ConnectTable connect table id
func ConnectTable(ctx *clu.Context, table string) (common.RegDbID, error) {
	databaseTableEntry, err := clu.SearchTable(table)
	if err != nil {
		return 0, err
	}
	refCopy := *databaseTableEntry.Reference
	password := databaseTableEntry.Database.Password
	if !databaseTableEntry.Database.AuthenticationGlobal {
		log.Log.Debugf("Using user authentication")
		refCopy.User = ctx.UserName()
		password = ctx.Pass
	}
	log.Log.Debugf("Connect table (register handle) %#v \n-> %#v", databaseTableEntry.Reference, refCopy)
	id, err := flynn.Handler(&refCopy, password)
	if err != nil {
		services.ServerMessage("Error connecting database %s:%d...%v",
			refCopy.Host, refCopy.Port, err)
		return 0, errorrepo.NewError("REST00200", err)
	}
	log.Log.Debugf("Got connectiion to database handle %s", id)
	return id, nil
}

// CloseTable close table id
func CloseTable(id common.RegDbID) {
	log.Log.Debugf("Close table and free database handle %s", id)
	id.Close()
	id.FreeHandler()
}

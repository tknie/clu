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
	"crypto/md5"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/tknie/clu"
	"github.com/tknie/errorrepo"
	"github.com/tknie/flynn"
	"github.com/tknie/flynn/common"
	"github.com/tknie/log"
	"github.com/tknie/services"
)

type databaseRegister struct {
	readCount uint64
	reference *common.Reference
	database  *Database
}

// dbDictionary map of hash to database registry entry
var dbDictionary = sync.Map{}

// dbTableMap map of database table and registry entry
var dbTableMap = sync.Map{}

// init config update tracker to tracke changes
func init() {
	RegisterConfigUpdates(initRegister)
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

func databaseHash(dm *Database) string {
	return fmt.Sprintf("%X", md5.Sum([]byte(dm.String())))
}

// Handles handle database
func Handles(dm *Database) (*common.Reference, error) {
	dHash := databaseHash(dm)
	if e, ok := dbDictionary.Load(dHash); ok {
		regEntry := e.(*databaseRegister)
		log.Log.Debugf("Found database hash %s", dHash)
		atomic.AddUint64(&regEntry.readCount, 1)
		return regEntry.reference, nil
	}
	log.Log.Debugf("Add database hash %s", dHash)
	target := os.ExpandEnv(dm.Target)
	log.Log.Debugf("Handles %s", target)
	ref, _, err := common.NewReference(target)
	if err != nil {
		return nil, fmt.Errorf("error parsing target <%s>: %s -> %s", dm.Target, err, target)
	}
	log.Log.Debugf("Register database handler for target %s", dm.Target)
	_, err = flynn.Handler(ref, os.ExpandEnv(dm.Password))
	if err != nil {
		services.ServerMessage("Error registering database <%s>: %v", dm.Target, err)
		return nil, fmt.Errorf("error registering database")
	}
	dbDictionary.Store(dHash,
		&databaseRegister{reference: ref, readCount: 1, database: dm})
	for i := 0; i < len(dm.Tables); i++ {
		dm.Tables[i] = strings.ToLower(dm.Tables[i])
	}
	services.ServerMessage("Registered database driver=%s to %s:%d/%s",
		dm.Driver, ref.Host, ref.Port, ref.Database)
	return ref, nil
}

func initTableOfDatabases() {
	for _, dm := range Viewer.Database.DatabaseAccess.Database {
		Handles(&dm)
	}
}

// checkFilter checks the filters array if it match to the given table
func checkFilter(filters []string, table string) bool {
	log.Log.Debugf("Check filters: %v search %s", filters, table)
	if len(filters) == 0 {
		return true
	}
	checkTable := strings.ToLower(table)
	for _, filter := range filters {
		if ok, _ := filepath.Match(strings.ToLower(filter), checkTable); ok {
			return true
		}
	}
	return slices.Contains(filters, strings.ToLower(table))
}

func loadTableOfDatabases() {
	log.Log.Debugf("Refreshing database list")
	for _, dm := range Viewer.Database.DatabaseAccess.Database {
		log.Log.Debugf("Access database %s with user %s", dm.Target, dm.User)
		//u := dm.URL
		//m := regexp.MustCompile(`(?m):[^:]*@`)
		//m := regexp.MustCompile(`(?m)\${[^{]*PASS[^}]*}`)
		//res := m.ReplaceAllString(u, ":****@")
		id, err := Handles(&dm)
		if err != nil {
			continue
		}

		newDatabases := make([]string, 0)
		for _, table := range flynn.Maps() {
			s := strings.ToLower(table)
			if sid, ok := dbTableMap.Load(s); ok {
				if sid.(*databaseRegister).reference != id {
					services.ServerMessage("Found table on different databases: [%s]", s)
				}
			} else {
				if checkFilter(dm.Tables, table) {
					log.Log.Debugf("Append table: %s", table)
					newDatabases = append(newDatabases, s)
					dbTableMap.Store(s, &databaseRegister{reference: id, database: &dm})
				} else {
					log.Log.Debugf("Ignore table: %s", table)
				}
			}
		}
		if len(newDatabases) > 0 {
			services.ServerMessage("Collected %04d table(s) in dictionary", len(newDatabases))
		}
	}
	tableStat := "Registered Database access\n  "
	counter := 0
	dbTableMap.Range(func(key, value any) bool {
		tableEntry := value.(*databaseRegister)
		// log.Log.Infof("Database with table %s count: %d", key, tableEntry.readCount)
		if tableEntry.readCount > 0 {
			counter++
			if counter%8 == 0 {
				tableStat += "\n  "
			}
			tableStat += fmt.Sprintf("%s=%d ", key, tableEntry.readCount)
		}
		return true
	})
	if counter > 0 {
		log.Log.Infof(tableStat)
	} else {
		log.Log.Infof("No database access registered")
	}

}

// InitDatabases initialize database reference IDs
func InitDatabases() {
	log.Log.Debugf("Init databases done")
}

// GetAllViews get all table and view names
func GetAllViews() []string {
	viewList := make([]string, 0)
	dbTableMap.Range(func(key, value any) bool {
		viewList = append(viewList, key.(string))
		return true
	})
	return viewList
}

// SearchTable search table ref ID
func searchTable(table string) (*databaseRegister, error) {
	name := strings.ToLower(table)
	if d, ok := dbTableMap.Load(name); ok {
		dicEntry := d.(*databaseRegister)
		atomic.AddUint64(&dicEntry.readCount, 1)
		return dicEntry, nil
	}
	return nil, errorrepo.NewError("RERR01000", table)
}

// ConnectTable connect table id
func ConnectTable(ctx *clu.Context, table string) (common.RegDbID, error) {
	databaseTableEntry, err := searchTable(table)
	if err != nil {
		return 0, err
	}
	refCopy := *databaseTableEntry.reference
	password := databaseTableEntry.database.Password
	if !databaseTableEntry.database.AuthenticationGlobal {
		log.Log.Debugf("Using user authentication")
		refCopy.User = ctx.UserName()
		password = ctx.Pass
	}
	log.Log.Debugf("Connect table (register handle) %#v \n-> %#v", databaseTableEntry.reference, refCopy)
	id, err := flynn.Handler(&refCopy, password)
	if err != nil {
		services.ServerMessage("Error connecting database %s:%d...%v",
			refCopy.Host, refCopy.Port, err)
		return 0, fmt.Errorf("error connecting database")
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

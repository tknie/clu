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
	"time"

	"github.com/tknie/clu"
	"github.com/tknie/errorrepo"
	"github.com/tknie/flynn"
	"github.com/tknie/flynn/common"
	"github.com/tknie/log"
	"github.com/tknie/services"
)

var dbList = make(map[string]*common.Reference)
var dbDictionary = make(map[string]*common.Reference)

func init() {
	RegisterConfigUpdates(initRegister)
}

func initRegister() {
	log.Log.Debugf("Register databases")
	initTableOfDatabases()
	// loadTableOfDatabases()
	log.Log.Debugf("Start table thread for databases")
	go loadTableThread()
}

func loadTableThread() {
	//uptimeTicker := time.NewTicker(5 * time.Second)
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
	if ref, ok := dbList[dHash]; ok {
		log.Log.Debugf("Found database hash %s -> %s", dHash, os.ExpandEnv(dm.String()))
		return ref, nil
	}
	log.Log.Infof("Add database hash %s -> %s", dHash, os.ExpandEnv(dm.String()))
	target := os.ExpandEnv(dm.Target)
	log.Log.Debugf("Handles %s", target)
	ref, _, err := common.NewReference(target)
	if err != nil {
		return nil, fmt.Errorf("error parsing target <%s>: %s -> %s", dm.Target, err, target)
	}
	log.Log.Debugf("Register database handler %#v", dm)
	_, err = flynn.Handler(ref, os.ExpandEnv(dm.Password))
	if err != nil {
		services.ServerMessage("Error registering database <%s>: %v", dm.Target, err)
		return nil, fmt.Errorf("error registering database")
	}
	dbList[dHash] = ref
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
		log.Log.Debugf("Access database %#v", dm)
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
			if sid, ok := dbDictionary[s]; ok {
				if sid != id {
					services.ServerMessage("Found table on different databases: [%s]", s)
				}
			} else {
				if checkFilter(dm.Tables, table) {
					log.Log.Debugf("Append table: %s", table)
					newDatabases = append(newDatabases, s)
					dbDictionary[s] = id
				} else {
					log.Log.Debugf("Ignore table: %s", table)
				}
			}
		}
		if len(newDatabases) > 0 {
			services.ServerMessage("Found table(s) to dictionary:\n%v", newDatabases)
		}
	}

}

// InitDatabases initialize database reference IDs
func InitDatabases() {
	log.Log.Debugf("Init databases done")
}

// GetAllViews get all table and view names
func GetAllViews() []string {
	viewList := make([]string, 0)
	for k := range dbDictionary {
		viewList = append(viewList, k)
	}
	return viewList
}

// SearchTable search table ref ID
func SearchTable(table string) (*common.Reference, error) {
	name := strings.ToLower(table)
	if d, ok := dbDictionary[name]; ok {
		return d, nil
	}
	return nil, errorrepo.NewError("RERR01000", table)
}

// ConnectTable connect table id
func ConnectTable(ctx *clu.Context, table string) (common.RegDbID, error) {
	ref, err := SearchTable(table)
	if err != nil {
		return 0, err
	}
	refCopy := *ref
	refCopy.User = ctx.User.User
	log.Log.Debugf("Connect table (register handle) %#v -> %#v", ref, refCopy)
	id, err := flynn.Handler(&refCopy, ctx.Pass)
	if err != nil {
		services.ServerMessage("Error registering database %s:%d...%v",
			ref.Host, ref.Port, err)
		return 0, fmt.Errorf("error registering database")
	}
	log.Log.Debugf("Got register database handle %s", id)
	return id, nil
}

// CloseTable close table id
func CloseTable(id common.RegDbID) {
	log.Log.Debugf("Close table and free database handle %s", id)
	id.Close()
	id.FreeHandler()
}

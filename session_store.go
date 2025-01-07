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

package clu

import (
	"os"
	"sync"
	"time"

	"github.com/tknie/errorrepo"
	"github.com/tknie/flynn"
	"github.com/tknie/flynn/common"
	"github.com/tknie/log"
	"github.com/tknie/services"
	"github.com/tknie/services/auth"
)

var sessionTableName = ""

var sessionDbRef *common.Reference
var sessionDbPassword = ""

var sessionExpirerDuration = time.Duration(6) * time.Hour
var sessionInfoMap = sync.Map{}
var sessionLock sync.Mutex

var chanUpdateSessionInfo = make(chan *auth.SessionInfo, 10)
var chanRemoveSessionInfo = make(chan *auth.SessionInfo, 10)

// DeleteUUID delete UUID after regular time frame
var DeleteUUID = false

func openSessionStore() (common.RegDbID, error) {
	var err error
	if sessionDbPassword == "" {
		sessionDbPassword = os.Getenv("REST_SESSION_LOG_PASS")
	}
	sessionStoreID, err := flynn.Handler(sessionDbRef, sessionDbPassword)
	if err != nil {
		services.ServerMessage("Register error log: %v", err)
		return 0, err
	}
	return sessionStoreID, nil
}

// InitStoreInfo init session info storage
func InitStoreInfo(ref *common.Reference, password, tablename string) {
	sessionDbRef = ref
	sessionDbPassword = password

	sessionStoreID, err := openSessionStore()
	if err != nil {
		return
	}
	defer sessionStoreID.FreeHandler()
	defer sessionStoreID.Close()
	log.Log.Debugf("Receive session store handler %s", sessionStoreID)

	dbTables := flynn.Maps()
	for _, d := range dbTables {
		if d == tablename {
			go updateSessionInfoThread()
			auth.JWTOperator = &StoreJWTHandler{}
			sessionTableName = tablename
			services.ServerMessage("Storing session data to table '%s'", sessionTableName)
			return
		}
	}
	su := &auth.SessionInfo{}
	err = sessionStoreID.CreateTable(tablename, su)
	if err != nil {
		services.ServerMessage("Database session store creating failed: %v", err)
		return
	}
	services.ServerMessage("Database session store created successfully")
	go updateSessionInfoThread()
	sessionTableName = tablename
	services.ServerMessage("Creating and storing session data to table '%s'", sessionTableName)
	auth.JWTOperator = &StoreJWTHandler{}
}

// StoreJWTHandler store session in a database store
type StoreJWTHandler struct {
}

// UUIDInfo get UUID info User information
func (st *StoreJWTHandler) UUIDInfo(uuid string) (*auth.SessionInfo, error) {
	log.Log.Debugf("Search UUID info for %s", uuid)
	if sessionInfo, ok := sessionInfoMap.Load(uuid); ok {
		return sessionInfo.(*auth.SessionInfo), nil
	}
	sessionLock.Lock()
	defer sessionLock.Unlock()
	sessionStoreID, err := openSessionStore()
	if err != nil {
		return nil, err
	}
	defer sessionStoreID.FreeHandler()
	defer sessionStoreID.Close()

	q := &common.Query{TableName: sessionTableName,
		Search:     "uuid='" + uuid + "'",
		DataStruct: &auth.SessionInfo{},
		Fields:     []string{"*"}}
	result, err := sessionStoreID.Query(q, func(search *common.Query, result *common.Result) error {
		return nil
	})
	if err != nil {
		return nil, err
	}
	si := result.Data.(*auth.SessionInfo)
	log.Log.Debugf("Found session %v", si.UUID)
	return si, nil
}

// Range go through all session entries
func (st *StoreJWTHandler) Range(f func(uuid, value any) bool) error {
	sessionStoreID, err := openSessionStore()
	if err != nil {
		return err
	}
	defer sessionStoreID.FreeHandler()
	defer sessionStoreID.Close()

	q := &common.Query{TableName: sessionTableName,
		Search:     "invalidated  < now()",
		DataStruct: &auth.SessionInfo{},
		Fields:     []string{"*"}}
	_, err = sessionStoreID.Query(q, func(search *common.Query, result *common.Result) error {
		s := result.Data.(*auth.SessionInfo)
		elapsed := s.Invalidated
		if !f(s.UUID, elapsed) {
			return errorrepo.NewError("REST00004")
		}
		return nil
	})
	return err
}

// InvalidateUUID invalidate UUID entry and given elapsed time
func (st *StoreJWTHandler) InvalidateUUID(uuid string, elapsed time.Time) bool {
	log.Log.Debugf("Invalidate session info %s -> %v", uuid, elapsed)
	si, err := st.UUIDInfo(uuid)
	if si == nil || err != nil {
		return false
	}
	log.Log.Debugf("Trigger remove session info %s", uuid)
	si.Invalidated = time.Now()
	chanRemoveSessionInfo <- si
	return true
}

// Store store entry for given input
func (st *StoreJWTHandler) Store(principal auth.PrincipalInterface, user, pass string) (err error) {
	log.Log.Debugf("Store session info %s", user)
	si := principal.Session().(*auth.SessionInfo)
	si.LastAccess = time.Now()
	data, err := auth.EncryptData(pass)
	if err != nil {
		return err
	}
	si.Data = []byte(data)
	si.Invalidated = si.LastAccess.Add(sessionExpirerDuration)
	insert := &common.Entries{Fields: []string{"*"}, DataStruct: si}
	insert.Values = [][]any{{si}}
	sessionLock.Lock()
	defer sessionLock.Unlock()
	sessionStoreID, err := openSessionStore()
	if err != nil {
		return err
	}
	defer sessionStoreID.FreeHandler()
	defer sessionStoreID.Close()
	log.Log.Debugf("Store session value %#v", si.UUID)
	_, err = sessionStoreID.Insert(sessionTableName, insert)
	if err != nil {
		log.Log.Errorf("Error storing user %s: %v", si.UUID, err)
		return err
	}
	log.Log.Debugf("Commiting storing session: %s", si.UUID)
	err = sessionStoreID.Commit()
	return err
}

// ValidateUUID validate JWT claims are in UUID session list
func (st *StoreJWTHandler) ValidateUUID(claims *auth.JWTClaims) (auth.PrincipalInterface, bool) {
	log.Log.Debugf("Valiadte UUID %s", claims.UUID)
	si, err := st.UUIDInfo(claims.UUID)
	if err != nil {
		log.Log.Errorf("Session with UUID %s not found", claims.UUID)
		return nil, false
	}
	log.Log.Debugf("Found valid session for UUID %s", si.UUID)
	pass, err := auth.DecryptData(string(si.Data))
	if err != nil {
		log.Log.Errorf("Error decrypt data %v", err)
	}
	si.LastAccess = time.Now()
	si.Invalidated = si.LastAccess.Add(sessionExpirerDuration)
	chanUpdateSessionInfo <- si
	sessionInfoMap.Delete(claims.UUID)
	p := auth.PrincipalCreater(si, si.User, pass)
	return p, true
}

func updateSessionInfoThread() {
	for {
		log.Log.Debugf("Waiting session info for updates or remove")
		select {
		case si := <-chanUpdateSessionInfo:
			updateSessionInfo(si)
		case si := <-chanRemoveSessionInfo:
			if DeleteUUID {
				deleteUUID(si)
			} else {
				updateSessionInfo(si)
			}
		case <-time.After(30 * time.Second):
			log.Log.Debugf("Shift working 30 seconds(session store)")
		}
	}
}

func updateSessionInfo(si *auth.SessionInfo) {
	sessionLock.Lock()
	defer sessionLock.Unlock()
	sessionStoreID, err := openSessionStore()
	if err != nil {
		return
	}
	defer sessionStoreID.FreeHandler()
	defer sessionStoreID.Close()
	update := &common.Entries{Fields: []string{"Invalidated"}, DataStruct: si}
	update.Values = [][]any{{si}}
	update.Update = []string{"UUID='" + si.UUID + "'"}
	log.Log.Debugf("Update value %#v", si.UUID)
	_, c, err := sessionStoreID.Update(sessionTableName, update)
	if err != nil {
		log.Log.Errorf("Error storing session: %v", err)
		return
	}
	log.Log.Debugf("Commiting session update: (%d)", c)
	err = sessionStoreID.Commit()
	if err == nil {
		return
	}
}

func deleteUUID(si *auth.SessionInfo) {
	sessionLock.Lock()
	defer sessionLock.Unlock()
	sessionStoreID, err := openSessionStore()
	if err != nil {
		return
	}
	defer sessionStoreID.FreeHandler()
	defer sessionStoreID.Close()
	remove := &common.Entries{Criteria: "uuid = '" + si.UUID + "'"}
	log.Log.Debugf("Remove UUID %s", si.UUID)
	_, err = sessionStoreID.Delete(sessionTableName, remove)
	if err != nil {
		log.Log.Errorf("Error deleting session: %v", err)
		return
	}
	err = sessionStoreID.Commit()
	if err == nil {
		return
	}

}

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
	"time"

	"github.com/tknie/flynn"
	"github.com/tknie/flynn/common"
	"github.com/tknie/log"
	"github.com/tknie/services"
	"github.com/tknie/services/auth"
)

var sessionURL = ""
var sessionTableName = ""
var disableStore = false
var sessionStoreID common.RegDbID

var sessionExpirerDuration = time.Duration(6) * time.Hour
var sessionInfoMap = make(map[string]*auth.SessionInfo)
var sessionLock sync.Mutex

var chanUpdateSessionInfo = make(chan *auth.SessionInfo)

// InitStoreInfo init session info storage
func InitStoreInfo(ref *common.Reference) {
	sessionURL = os.Getenv("REST_SESSION_LOG_URL")
	sessionTableName = os.Getenv("REST_SESSION_LOG_TABLENAME")
	if sessionURL == "" || sessionTableName == "" {
		services.ServerMessage("Log parameter storage disabled...")
		log.Log.Debugf("SESSION_STORE: Disable due to URL error")
		disableStore = true
		return
	}
	var err error
	userDbRef, userDbPassword, err = common.NewReference(sessionURL)
	if err != nil {
		log.Log.Fatal("REST audit URL incorrect: " + sessionURL)
	}
	if userDbPassword == "" {
		userDbPassword = os.Getenv("REST_SESSION_LOG_PASS")
	}
	userDbRef.User = "admin"

	services.ServerMessage("Storing session data to table '%s'", sessionTableName)
	sessionStoreID, err = flynn.Handler(userDbRef, userDbPassword)
	if err != nil {
		services.ServerMessage("Register error log: %v", err)
		return
	}
	log.Log.Debugf("Receive session store handler %s", sessionStoreID)

	dbTables := flynn.Maps()
	for _, d := range dbTables {
		if d == sessionTableName {
			auth.JWTOperator = &StoreJWTHandler{}
			return
		}
	}
	su := &auth.SessionInfo{}
	err = sessionStoreID.CreateTable(sessionTableName, su)
	if err != nil {
		services.ServerMessage("Database session store creating failed: %v", err)
		return
	}
	services.ServerMessage("Database session store created successfully")
	go updaterSessionInfo()
	auth.JWTOperator = &StoreJWTHandler{}
}

// StoreJWTHandler store session in a database store
type StoreJWTHandler struct {
}

// UUIDInfo get UUID info User information
func (st *StoreJWTHandler) UUIDInfo(uuid string) (*auth.SessionInfo, error) {
	log.Log.Debugf("Search UUID info for %s", uuid)
	if sessionInfo, ok := sessionInfoMap[uuid]; ok {
		return sessionInfo, nil
	}
	sessionLock.Lock()
	defer sessionLock.Unlock()
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
	log.Log.Debugf("Found session %v", result.Data)
	return result.Data.(*auth.SessionInfo), nil
}

// Range go through all session entries
func (st *StoreJWTHandler) Range(f func(uuid, value any) bool) error {
	q := &common.Query{TableName: sessionTableName,
		DataStruct: &auth.SessionInfo{},
		Fields:     []string{"*"}}
	_, err := sessionStoreID.Query(q, func(search *common.Query, result *common.Result) error {
		s := result.Data.(*auth.SessionInfo)
		elapsed := s.Created.Add(sessionExpirerDuration)
		if !f(s.UUID, elapsed) {
			return fmt.Errorf("aborted Range")
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
	if si.Invalidated.IsZero() && !si.Invalidated.Before(time.Now()) {
		si.Invalidated = time.Now()
		chanUpdateSessionInfo <- si
	}
	return false
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
	log.Log.Debugf("Store session value %#v", si.UUID)
	err = userStoreID.Insert(sessionTableName, insert)
	if err != nil {
		log.Log.Errorf("Error storing user: %v", err)
		return err
	}
	log.Log.Errorf("Insert storing session: %s", si.UUID)
	err = userStoreID.Commit()
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
	p := auth.PrincipalCreater(si, si.User, pass)
	return p, true
}

func updaterSessionInfo() {
	for {
		select {
		case si := <-chanUpdateSessionInfo:
			update := &common.Entries{Fields: []string{"Invalidated"}, DataStruct: si}
			update.Values = [][]any{{si}}
			log.Log.Debugf("Update value %#v", si.UUID)
			c, err := userStoreID.Update(sessionTableName, update)
			if err != nil {
				log.Log.Errorf("Error storing session: %v", err)
				continue
			}
			log.Log.Errorf("Update storing session: (%d)", c)
			err = userStoreID.Commit()
			if err == nil {
				continue
			}
		case <-time.After(30 * time.Second):
			log.Log.Debugf("Shift working 30 seconds")
		}
	}
}

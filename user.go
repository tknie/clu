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
	"net/mail"
	"os"
	"sync"
	"time"

	"github.com/tknie/flynn"
	"github.com/tknie/flynn/common"
	"github.com/tknie/log"
	"github.com/tknie/services"
	"github.com/tknie/services/auth"
)

var userTableName = ""
var disableUser = false
var userDbRef *common.Reference
var userDbPassword = ""

// DefaultRole default role for automatic added user
const DefaultRole = "Reader"

// var userStoreID common.RegDbID

var userFieldList = []string{"Name", "Created", "LastLogin", "Permission"}
var userFieldListUpdate = []string{"Created", "LastLogin"}

type storeUserInfo struct {
	stored   bool
	created  time.Time
	userInfo *auth.UserInfo
}

var userInfoMap = sync.Map{} // make(map[string]*storeUserInfo)

var userLock sync.Mutex

// InitUserInfo init user info evaluation
func InitUserInfo(ref *common.Reference, password, tablename string) {
	userDbRef = ref
	userDbPassword = password
	if userDbPassword == "" {
		userDbPassword = os.Getenv("REST_USER_LOG_PASS")
	}
	userTableName = tablename
	userStoreID, err := openUserStore()
	if err != nil {
		services.ServerMessage("Storing audit data error: %v", err)
		return
	}
	defer userStoreID.FreeHandler()
	defer userStoreID.Close()
	log.Log.Debugf("Receive user store handler %s", userStoreID)
	services.ServerMessage("Storing audit data to table '%s'", userTableName)

	go updaterUserInfo()

	dbTables := flynn.Maps()
	for _, d := range dbTables {
		if d == userTableName {
			return
		}
	}
	su := &auth.UserInfo{}
	err = userStoreID.CreateTable(userTableName, su)
	if err != nil {
		services.ServerMessage("Database user store creating failed: %v", err)
		return
	}
	services.ServerMessage("Database user store created successfully")
}

func updaterUserInfo() {
	for {
		log.Log.Debugf("Waiting user info for updates or remove")
		select {
		case <-time.After(30 * time.Second):
			timeRange := time.Now().Add(time.Duration(-2) * time.Minute)
			log.Log.Debugf("Shift working 30 seconds (user info)")
			userInfoMap.Range(func(key, value any) bool {
				st := value.(*storeUserInfo)
				if st.userInfo.LastLogin.Before(timeRange) {
					userInfoMap.Delete(key)
				}
				return true
			})
		}
	}
}

func openUserStore() (common.RegDbID, error) {
	if userDbPassword == "" {
		userDbPassword = os.Getenv("REST_USER_LOG_PASS")
	}
	userStoreID, err := flynn.Handler(userDbRef, userDbPassword)
	if err != nil {
		services.ServerMessage("Register error log: %v", err)
		return 0, err
	}
	return userStoreID, nil
}

// QueryUser query user information
func QueryUser(user string) *auth.UserInfo {
	if disableUser {
		return nil
	}
	if u, ok := userInfoMap.Load(user); ok {
		return u.(*storeUserInfo).userInfo
	}
	userLock.Lock()
	defer userLock.Unlock()
	userStoreID, err := openUserStore()
	if err != nil {
		return nil
	}
	defer userStoreID.FreeHandler()
	defer userStoreID.Close()
	var userInfo *auth.UserInfo
	count := 0
	q := &common.Query{TableName: userTableName,
		Search:     "name='" + user + "'",
		DataStruct: &auth.UserInfo{},
		Fields:     []string{"*"}}
	_, err = userStoreID.Query(q, func(search *common.Query, result *common.Result) error {
		count++
		if userInfo == nil {
			userInfo = result.Data.(*auth.UserInfo)
		} else {
			services.ServerMessage("%s not unique %03d", user, count)
		}
		return nil
	})
	if err != nil {
		return nil
	}
	if userInfo != nil {
		log.Log.Debugf("Create user info storage %#v", userInfo)
		userInfoMap.Store(userInfo.User, &storeUserInfo{true, time.Now(), userInfo})
		return userInfo
	}
	return nil
}

// CheckUserExist check user already exists and create if not available
func CheckUserExist(user string) *auth.UserInfo {
	userInfo := QueryUser(user)
	if userInfo == nil {
		log.Log.Debugf("No user %s in user info found", user)

		userInfo = &auth.UserInfo{User: user, Created: time.Now(), LastLogin: time.Now(),
			Permission: &auth.User{Name: DefaultRole, Read: "*", Write: ""}}
		if _, err := mail.ParseAddress(user); err == nil {
			userInfo.EMail = user
		}
		log.Log.Debugf("Create user info new %#v", userInfo)
		userInfoMap.Store(userInfo.User, &storeUserInfo{false, time.Now(), userInfo})
	} else {
		log.Log.Debugf("User info %s found: %#v / %#v", user, userInfo, userInfo.Permission)
	}
	return userInfo
}

// AddUserInfo add user if not already exists and create if not available
func AddUserInfo(userInfo *auth.UserInfo) error {
	if disableUser {
		return nil
	}
	if u, ok := userInfoMap.Load(userInfo.User); ok {
		userLock.Lock()
		defer userLock.Unlock()
		userStoreID, err := openUserStore()
		if err != nil {
			return nil
		}
		defer userStoreID.FreeHandler()
		defer userStoreID.Close()
		insert := &common.Entries{Fields: userFieldListUpdate, DataStruct: userInfo}
		insert.Values = [][]any{{userInfo}}
		sui := u.(*storeUserInfo)
		if sui.stored {
			insert.Fields = userFieldListUpdate
			log.Log.Debugf("Update value %#v", userInfo)
			insert.Update = []string{"user='" + userInfo.User + "'"}
			_, _, err := userStoreID.Update(userTableName, insert)
			if err != nil {
				log.Log.Errorf("Error updating user info: %v", err)
				return err
			}
		} else {
			insert.Fields = userFieldList
			log.Log.Debugf("Insert value %#v", userInfo)
			_, err := userStoreID.Insert(userTableName, insert)
			if err != nil {
				log.Log.Errorf("Error inserting user info: %v", err)
				// Could be an outage of database, refresh user information
				CheckUserExist(userInfo.User)
				return err
			}
		}
		err = userStoreID.Commit()
		if err != nil {
			log.Log.Errorf("Error commiting user info: %v", err)
			return err
		}
		sui.stored = true
	}
	return nil
}

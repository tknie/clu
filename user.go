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
	"os"
	"time"

	"github.com/tknie/flynn"
	"github.com/tknie/flynn/common"
	"github.com/tknie/log"
	"github.com/tknie/services"
	"github.com/tknie/services/auth"
)

var userURL = ""
var userTableName = ""
var disableUser = false
var userDbRef *common.Reference
var userDbPassword = ""
var userStoreID common.RegDbID

var userFieldList = []string{"Name", "Created"}

// User user information
type User struct {
	Info      *auth.UserInfo
	Thumbnail []byte
}

// InitUserInfo init user info evaluation
func InitUserInfo() {
	userURL = os.Getenv("REST_USER_LOG_URL")
	userTableName = os.Getenv("REST_USER_LOG_TABLENAME")
	if userURL == "" || userTableName == "" {
		services.ServerMessage("Log parameter storage disabled...")
		log.Log.Debugf("USER_AUDIT: Disable due to URL error")
		disableUser = true
		return
	}
	var err error
	userDbRef, userDbPassword, err = common.NewReference(userURL)
	if err != nil {
		log.Log.Fatal("REST audit URL incorrect: " + userURL)
	}
	if userDbPassword == "" {
		userDbPassword = os.Getenv("REST_USER_LOG_PASS")
	}
	userDbRef.User = "admin"

	services.ServerMessage("Storing audit data to table '%s'", userTableName)
	userStoreID, err = flynn.Handler(userDbRef, userDbPassword)
	if err != nil {
		services.ServerMessage("Register error log: %v", err)
		return
	}
	log.Log.Debugf("Receive user store handler %s", userStoreID)

	dbTables := flynn.Maps()
	for _, d := range dbTables {
		if d == userTableName {
			return
		}
	}
	su := &User{}
	err = userStoreID.CreateTable(userTableName, su)
	if err != nil {
		services.ServerMessage("Database user store creating failed: %v", err)
		return
	}
	services.ServerMessage("Database user store created successfully")
}

// QueryUser query user information
func QueryUser(user string) *User {
	if disableUser {
		return nil
	}
	var userInfo *User
	count := 0
	q := &common.Query{TableName: userTableName,
		Search:     "name='" + user + "'",
		DataStruct: &User{},
		Fields:     []string{"*"}}
	_, err := userStoreID.Query(q, func(search *common.Query, result *common.Result) error {
		count++
		if userInfo == nil {
			userInfo = result.Data.(*User)
		} else {
			services.ServerMessage("%s not unique %03d", user, count)
		}
		return nil
	})
	if err != nil {
		return nil
	}
	if userInfo != nil {
		return userInfo
	}
	return nil
}

// CheckUserExist check user already exists and create if not available
func CheckUserExist(user string, session *auth.SessionInfo) *User {
	if disableUser {
		return nil
	}
	userInfo := QueryUser(user)
	if userInfo == nil {
		userInfo = &User{Info: &auth.UserInfo{User: user, Created: time.Now()}}
		insert := &common.Entries{Fields: userFieldList, DataStruct: userInfo}
		insert.Values = [][]any{{userInfo}}
		log.Log.Debugf("Insert value %#v", userInfo.Info)
		err := userStoreID.Insert(userTableName, insert)
		if err != nil {
			return nil
		}
		log.Log.Errorf("Error storing user: %v", err)
		err = userStoreID.Commit()
		if err != nil {
			return nil
		}
	}
	return userInfo
}

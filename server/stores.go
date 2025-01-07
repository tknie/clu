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
	"log"
	"os"

	"github.com/tknie/clu"
)

// InitDatabaseStores init database session store
func InitDatabaseStores() error {
	dm := clu.Viewer.Database.UserInfo
	if dm != nil {
		r, err := dm.Handles()
		if err == nil {
			clu.InitUserInfo(r, os.ExpandEnv(dm.Password), os.ExpandEnv(dm.Table))
		} else {
			log.Fatal("user information store not being able to start:", err)
		}
	}
	if clu.Viewer.Database.SessionInfo != nil {
		dm = clu.Viewer.Database.SessionInfo.Database
		if dm != nil {
			clu.DeleteUUID = clu.Viewer.Database.SessionInfo.DeleteUUID
			r, err := dm.Handles()
			if err == nil {
				clu.InitStoreInfo(r, os.ExpandEnv(dm.Password), os.ExpandEnv(dm.Table))
			} else {
				log.Fatal("session information store not being able to start:", err)
			}
		}
	}
	go clu.InitBatchWatcherThread()

	return nil
}

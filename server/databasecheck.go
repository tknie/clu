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
	"os"

	"github.com/tknie/log"
	"github.com/tknie/services"
	"github.com/tknie/services/auth"
)

var adadatadir string
var installation []string

// InitDirectAccess init Adabas direct access function, only need in active server
var InitDirectAccess func(*RestServer)

// InitAdmin init Adabas admin function, only need in active server
var InitAdmin func(*RestServer)

// AddLocation add location, only needed in active server
var AddLocation func(name, location string) error

// Loader config loader structure calling static configurations
type Loader struct {
}

var loader = &Loader{}

// Logging logging configuration
func (loader *Loader) Logging(v interface{}) *services.Logging {
	viewer := v.(*RestServer)
	return viewer.Server.LogLocation
}

// SetLogging set logging configuration
func (loader *Loader) SetLogging(l *services.Logging) {
	Viewer.Server.LogLocation = l
}

// Default empty default config
func (loader *Loader) Default() interface{} {
	return &RestServer{}
}

// Current current active config
func (loader *Loader) Current() interface{} {
	return Viewer
}

// IsServer indicate interface to be a server
func (loader *Loader) IsServer() bool {
	return true
}

// Loaded executed after data load
func (loader *Loader) Loaded(nv interface{}) error {
	Viewer = nv.(*RestServer)
	log.Log.Debugf("Loaded configuration active")
	if Viewer.Common.Version == "" || Viewer.Common.Version == "v1" {
		// if Viewer.Version == "" {
		// 	// Viewer.Module.Admin.NoModification = false
		// }
		Viewer.Common.Version = "v2"

		if Viewer.Tasks.Directory == "" {
			Viewer.Tasks.Directory = adadatadir + string(os.PathSeparator) + "logs"
		}
	}

	auth.InitLoginService(&auth.Authentication{AuthenticationServer: Viewer.Server.LoginService.AuthenticationServer})

	// Init authentication/authorization infrastructure like database access and JWT
	Viewer.InitSecurityInfrastructure()

	return nil
}

// InitSecurityInfrastructure init configruation data
func (viewer *RestServer) InitSecurityInfrastructure() {

	if viewer.Server.Content == "" {
		viewer.Server.Content = "./static"
	}

	if viewer.Database.DatabaseAccess.Global {
		services.ServerMessage("Direct access granted to all database (global=true)")
	} else {
		// Init Adabas map, not needed if configuration script is used
		if InitDirectAccess != nil {
			InitDirectAccess(viewer)
		}
	}

	if InitAdmin != nil {
		InitAdmin(viewer)
	}

	// if len(viewer.JobStore.Database) > 0 {
	// 	jobs.Storage = &jobs.JobStore{Dbid: viewer.JobStore.Database[0].Dbid,
	// 		File: viewer.JobStore.Database[0].File,
	// 	}
	// }

	// Add File transfer locations
	if len(viewer.FileTransfer.Directories.Directory) == 0 {
		services.ServerMessage("No File location defined, file transfer not possible\n")
	} else {
		for _, d := range viewer.FileTransfer.Directories.Directory {
			if AddLocation != nil {
				AddLocation(d.Name, d.Location)
			}
		}
	}
	viewer.Server.WebToken.InitWebTokenJose2()
	log.Log.Debugf("Load of configuration finished")
}

// GetAdaDataDir get ADADATADIR configuration
func GetAdaDataDir() string {
	return adadatadir
}

// GetInstallation get defined installations
func GetInstallation() []string {
	return installation
}

// CloseConfig close configuration watcher
func (viewer *RestServer) CloseConfig() {
	// done <- true
	services.CloseConfig()
}

// String representation of Database instance
func (db *Database) String() string {
	return db.User + ":***@" + db.Host + ":" + db.Port
}

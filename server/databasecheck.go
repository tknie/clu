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
	"os"

	"github.com/tknie/clu"
	"github.com/tknie/log"
	"github.com/tknie/services"
	"github.com/tknie/services/auth"
)

// Loader config loader structure calling static configurations
type Loader struct {
}

var loader = &Loader{}

// Logging logging configuration
func (loader *Loader) Logging(v interface{}) *services.Logging {
	viewer := v.(*clu.RestServer)
	return viewer.Server.LogLocation
}

// SetLogging set logging configuration
func (loader *Loader) SetLogging(l *services.Logging) {
	clu.Viewer.Server.LogLocation = l
}

// Default empty default config
func (loader *Loader) Default() interface{} {
	return &clu.RestServer{}
}

// Current current active config
func (loader *Loader) Current() interface{} {
	return clu.Viewer
}

// IsServer indicate interface to be a server
func (loader *Loader) IsServer() bool {
	return true
}

// Loaded executed after data load
func (loader *Loader) Loaded(nv interface{}) error {
	clu.Viewer = nv.(*clu.RestServer)
	log.Log.Debugf("Loaded configuration active")
	if clu.Viewer.Common.Version == "" || clu.Viewer.Common.Version == "v1" {
		// if Viewer.Version == "" {
		// 	// Viewer.Module.Admin.NoModification = false
		// }
		clu.Viewer.Common.Version = "v2"

		if clu.Viewer.Tasks.Directory == "" {
			clu.Viewer.Tasks.Directory = clu.GetAdaDataDir() + string(os.PathSeparator) + "logs"
		}
	}

	auth.InitLoginService(&auth.Authentication{AuthenticationServer: clu.Viewer.Server.LoginService.AuthenticationServer})

	// Init authentication/authorization infrastructure like database access and JWT
	clu.Viewer.InitSecurityInfrastructure()

	return nil
}

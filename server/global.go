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
	"bytes"
	"crypto/sha1"
	"fmt"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/tknie/clu"
	"github.com/tknie/log"
	"github.com/tknie/services"
	"github.com/tknie/services/auth"
)

// INSTALLDIR installation directory environment
const INSTALLDIR = "INSTALLDIR"

// InitConfig load xml configuration file
// The components are used to load and inject the configuration
func InitConfig(watch bool) error {
	err := clu.LoadConfig(watch, loader)
	if err != nil {
		return err
	}
	if len(clu.Viewer.Server.LoginService.AuthenticationServer) == 0 {
		services.ServerMessage("No authentication configuration found, using default configuration")
		a := &auth.AuthenticationServer{Type: "file", PasswordFile: "configuration/realm.properties"}
		clu.Viewer.Server.LoginService.AuthenticationServer = append(clu.Viewer.Server.LoginService.AuthenticationServer, a)
	}
	auth.AuthenticationConfig = &auth.Authentication{AuthenticationServer: clu.Viewer.Server.LoginService.AuthenticationServer}
	if clu.Viewer.Server.WebToken != nil {
		err = clu.Viewer.Server.WebToken.InitWebTokenJose2()
		if err != nil {
			services.ServerErrorMessage("RERR00044", err.Error())
			os.Exit(44)
		}
	}
	err = auth.LoadUsers(auth.AdministratorRole, clu.Viewer.Server.LoginService.Administrators, true)
	if err != nil {
		return err
	}
	err = auth.LoadUsers(auth.UserRole, clu.Viewer.Server.LoginService.Users, true)
	if err != nil {
		return err
	}
	initEnvironmentConfig()
	return nil
}

// initEnvironmentConfig initialize environment configuration to preset
// a number of configuration before go-swagger is initialized.
func initEnvironmentConfig() {
	os.Setenv("PGAPPNAME", "CLU Rest-Server")

	// Check port configuration to be set via environment
	if clu.Viewer != nil {
		host := os.Getenv("HOST")
		if host == "" {
			os.Setenv("HOST", "")
		}
		host = os.Getenv("TLS_HOST")
		if host == "" {
			os.Setenv("TLS_HOST", "")
		}
		for _, s := range clu.Viewer.Server.Service {
			switch strings.ToLower(s.Type) {
			case "http":
				os.Setenv("PORT", strconv.Itoa(s.Port))
			case "https":
				os.Setenv("TLS_PORT", strconv.Itoa(s.Port))
			}
		}
	}

}

// AdaptConfig adapt and reload xml configuration file
func AdaptConfig(config string) {
	log.Log.Debugf("Adapting config %s", config)
	if config != "" && config != clu.CurrentConfig {
		services.LoadConfig(config, loader, true)
	}
	clu.LoadedConfig()
}

// StoreConfig store the current config, wraps the REST server configuration
// and uses components to store the current configuration
func StoreConfig() error {
	err := services.StoreConfig(clu.CurrentConfig, loader)
	return err
}

// ConvertByteArrayToString convert byte array to string by evaluating
// \0 end and create string from sub array
func ConvertByteArrayToString(name []byte, max int) string {
	n := bytes.IndexByte(name[:], 0)
	if n == -1 || n > max {
		n = max
	}
	return string(name[:n])
}

// GenerateShutdownHash generates shutdown hash dependent on pass code time
func GenerateShutdownHash() string {
	t := time.Now()
	tnow := t.Unix()
	torg := tnow - tnow%60
	h := sha1.New()
	dc := fmt.Sprintf("%s-%d", clu.Viewer.Server.Shutdown.Passcode.Value, torg)
	h.Write([]byte(dc))
	return fmt.Sprintf("%x", h.Sum(nil))
}

// DefaultPIDFile default PID file location
func DefaultPIDFile() string {
	p := os.Getenv("TEMP")
	if p == "" {
		p = os.Getenv(INSTALLDIR) + string(os.PathSeparator) + clu.InstallationDirName +
			string(os.PathSeparator) + "tmp"
	}
	return p + string(os.PathSeparator) + "server.pid"
}

// RemoteHost check if X-Forwarded-For is available and use this remote host
func RemoteHost(r *http.Request) string {
	remoteHost := r.Header.Get("X-Forwarded-For")
	if remoteHost == "" {
		remoteHost = r.RemoteAddr
	}
	log.Log.Debugf("Got forward host %s", remoteHost)
	ip6l := strings.LastIndex(remoteHost, "]")
	pa := strings.LastIndex(remoteHost, ":")
	if pa > -1 && pa > ip6l {
		remoteHost = remoteHost[:pa]
	}
	n, err := net.LookupAddr(remoteHost)
	if err != nil {
		log.Log.Debugf("Error evaluating %s: %v", remoteHost, err)
	} else {
		remoteHost = fmt.Sprintf("%s, %v", remoteHost, n)
	}
	log.Log.Debugf("Remote host set to %s", remoteHost)
	return remoteHost
}

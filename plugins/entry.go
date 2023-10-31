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

package plugins

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"plugin"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/tknie/clu"
	"github.com/tknie/clu/server"
	"github.com/tknie/log"
	"github.com/tknie/services"
)

// PluginTypes different types of plugins for
// - auditing
// - database operation
type PluginTypes int

const (
	// NoPlugin no plugin but may be used in module
	NoPlugin PluginTypes = iota
	// AuditPlugin auditing of RESTful server access
	AuditPlugin
	// AdabasPlugin Adabas plugin type
	AdabasPlugin
)

const suffix = ".so"

// Loader plugin Loader module to load plugin features
type Loader interface {
	Name() string
	Version() string
	Types() []int
	Stop()
}

// Audit auditing method to send to plugin
type Audit interface {
	ReceiveAudit(string, string, *http.Request)
	SendAudit(time.Duration, string, string, *http.Request)
	SendAuditError(time.Duration, string, string, *http.Request, error)
}

// AuditLoader auditing loader structure
type AuditLoader struct {
	Loader Loader
	Audit  Audit
}

// Adabas Adabas method to send to plugin
type Adabas interface {
	SendAdabas(time.Duration, interface{})
}

// AdabasLoader adabas plugin loader structure
type AdabasLoader struct {
	Loader Loader
	Adabas Adabas
}

var auditPlugins = make(map[string]*AuditLoader)
var adabasPlugins = make(map[string]*AdabasLoader)
var pluginsFound = false
var disablePlugin = false

var interrupt chan os.Signal
var once = new(sync.Once)
var shutOnce = new(sync.Once)

func signalNotify(interrupt chan<- os.Signal) {
	signal.Notify(interrupt, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
}

func handleInterrupt(interrupt chan os.Signal) {
	once.Do(func() {
		for range interrupt {
			ShutdownPlugins()
		}
	})
}

// InitPlugins initialize plugins in given plugin directory
func InitPlugins() {
	pluginDir, ok := os.LookupEnv("METAVERS_PLUGINS")
	if !ok {
		pluginDir = server.Viewer.Server.Plugins
	}
	pluginDir = os.ExpandEnv(pluginDir)
	if pluginDir == "" {
		return
	}
	pluginEnabled, filterPlugins := os.LookupEnv("METAVERS_PLUGENABLED")
	var plugList []string
	if filterPlugins {
		plugList = strings.Split(pluginEnabled, ",")
	}
	services.ServerMessage("Searching plugins in %s", pluginDir)
	err := filepath.Walk(pluginDir, func(path string, info os.FileInfo, err error) error {
		if info == nil {
			return nil
		}
		if !info.IsDir() && strings.HasSuffix(info.Name(), suffix) {
			plug, err := loadPlugin(pluginDir + "/" + info.Name())
			if err != nil {
				return nil
			}
			symLanguage, err := plug.Lookup("Loader")
			if err != nil {
				services.ServerMessage("Error opening plugin: %v", err)
				return nil
			}
			if loader, ok := symLanguage.(Loader); ok {
				found := !filterPlugins
				if !found && plugList != nil {
					n := loader.Name()
					for _, v := range plugList {
						if n == v {
							found = true
							break
						}
					}
				}
				if found {
					load(loader, info, plug)
				}
			} else {
				services.ServerMessage("Error opening plugin, error loading methods")
			}
		}
		return nil
	})
	if err != nil {
		return
	}

	if len(auditPlugins) > 0 {
		pluginsFound = true
	}

	interrupt = make(chan os.Signal, 1)
	signalNotify(interrupt)
	go handleInterrupt(interrupt)

}

// load loading the plugin
func load(loader Loader, info os.FileInfo, plug *plugin.Plugin) {
	pt := loader.Types()
	for _, t := range pt {
		switch t {
		case int(NoPlugin):
		case int(AuditPlugin):
			symAudit, err := plug.Lookup("Audit")
			if err != nil {
				services.ServerMessage("Error opening Audit plugin %s Version: %s : %v", loader.Name(), loader.Version(), err)
			} else {
				services.ServerMessage("Audit plugin: %s Version: %s", loader.Name(), loader.Version())
				audit := symAudit.(Audit)
				auditPlugins[info.Name()] = &AuditLoader{loader, audit}
			}
		case int(AdabasPlugin):
			symAdabas, err := plug.Lookup("Adabas")
			if err != nil {
				services.ServerMessage("Error opening Adabas plugin %s Version: %s : %v", loader.Name(), loader.Version(), err)
			} else {
				services.ServerMessage("Adabas plugin: %s Version: %s", loader.Name(), loader.Version())
				adaSym := symAdabas.(Adabas)
				adabasPlugins[info.Name()] = &AdabasLoader{loader, adaSym}
			}
		default:
			services.ServerMessage("Error opening plugin, unknown type: %v", t)
		}
	}
}

// ShutdownPlugins shutdown receiving message in plugins
func ShutdownPlugins() {
	shutOnce.Do(func() {
		disablePlugin = true
		services.ServerMessage("Shutdown all plugins ...")

		for _, v := range auditPlugins {
			v.Loader.Stop()
		}
		for _, v := range adabasPlugins {
			v.Loader.Stop()
		}
	})
}

func loadPlugin(mod string) (*plugin.Plugin, error) {
	// load module
	// 1. open the so file to load the symbols
	plug, err := plugin.Open(mod)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return plug, nil
}

// HasPlugins if any plugin is available to send data to
func HasPlugins() bool {
	return pluginsFound
}

// ReceiveAudit send audit information to plugins
func ReceiveAudit(p *clu.Context, r *http.Request) {
	if disablePlugin {
		return
	}
	log.Log.Debugf("Receive auditing plugins request: %v %v %v",
		r.Method, r.URL, server.RemoteHost(r))
	if p != nil {
		c := &http.Cookie{Name: "User", Value: p.UUID()}
		r.AddCookie(c)
		for _, x := range auditPlugins {
			x.Audit.ReceiveAudit(p.User, p.UUID(), r)
		}
		return
	}
	for _, x := range auditPlugins {
		x.Audit.ReceiveAudit(p.User, p.UUID(), r)
	}
}

// SendAudit send audit information to plugins
func SendAudit(started time.Time, w http.ResponseWriter, r *http.Request) {
	if disablePlugin || r.Method == "OPTIONS" {
		return
	}
	log.Log.Debugf("Send auditing plugins request: %v %v %v", r.Method, r.URL, server.RemoteHost(r))
	user := "Unknown"
	uuid := ""
	reqToken := r.Header.Get("Authorization")
	splitToken := strings.Split(reqToken, " ")
	switch strings.ToLower(splitToken[0]) {
	case "basic":
		user, _, _ = r.BasicAuth()
	case "bearer":
		reqToken = strings.TrimSpace(splitToken[1])
		p, err := server.Viewer.Server.WebToken.JWTContainsRoles(reqToken, []string{"admin"})
		if err != nil {
			uuid = "Not available"
		} else {
			c := p.(*clu.Context)
			uuid = c.UUID()
			user = c.User
		}
	}

	elapsed := time.Since(started)
	for _, x := range auditPlugins {
		x.Audit.SendAudit(elapsed, user, uuid, r)
	}
}

// SendAuditError send audit errors to plugins
func SendAuditError(started time.Time, r *http.Request, err error) {
	if disablePlugin {
		return
	}
	username, _, ok := r.BasicAuth()
	if !ok {
		services.ServerMessage("Basic authorization error expanding from %s", r.Host)
		username = "Unknown"
	}
	elapsed := time.Since(started)
	for _, x := range auditPlugins {
		x.Audit.SendAuditError(elapsed, username, "", r, err)
	}
}

// SendAdabasPlugins send adabas information to plugins
func SendAdabasPlugins(used time.Duration, ada interface{}) {
	if disablePlugin {
		return
	}
	for _, x := range adabasPlugins {
		x.Adabas.SendAdabas(used, ada)
	}
}

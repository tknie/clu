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

package plugins

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"plugin"
	"runtime/debug"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/tknie/clu"
	"github.com/tknie/clu/server"
	"github.com/tknie/log"
	"github.com/tknie/services"
	"github.com/tknie/services/auth"
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
	// AuthPlugin authorize and authorize
	AuthPlugin
	// ExtendPlugin extend entry point "/rest/extend"
	ExtendPlugin
	// ValidatorPlugin validate actions
	ValidatorPlugin
)

const suffix = ".so"

// Loader plugin Loader module to load plugin features
type Loader interface {
	Name() string
	Version() string
	Types() []PluginTypes
	Stop()
}

// Auth authenticate and authorize
type Auth interface {
	Init() error
	Authenticate(principal auth.PrincipalInterface, userName, passwd string) error
	Authorize(principal auth.PrincipalInterface, userName, passwd string) error
	CheckToken(token string, scopes []string) (auth.PrincipalInterface, error)
}

// AuthLoader auth loader structure
type AuthLoader struct {
	Loader Loader
	Auth   Auth
}

// Audit auditing method to send to plugin
type Audit interface {
	LoginAudit(string, string, *auth.SessionInfo, *auth.UserInfo)
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

// ExtendLoader extend loader plugin structure
type ExtendLoader struct {
	Loader Loader
	Extend server.RestExtend
}

// ValidatorLoader valiadtor loader plugin structure
type ValidatorLoader struct {
	Loader    Loader
	Validator server.RestValidator
}

var auditPlugins = make(map[string]*AuditLoader)
var adabasPlugins = make(map[string]*AdabasLoader)
var extendPlugins = make(map[string]*ExtendLoader)
var validatorPlugins = make(map[string]*ValidatorLoader)
var authPlugins = make(map[string]*AuthLoader)

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

func init() {
	auth.TriggerInvalidUUID = func(a *auth.SessionInfo, u *auth.UserInfo) {
		log.Log.Debugf("Logoff (invalidate UUID) triggered")
		LoginAudit("LOGIN", "logoff", a, u)
	}
}

// InitPlugins initialize plugins in given plugin directory
func InitPlugins() {
	pluginDir, ok := os.LookupEnv("METAVERS_PLUGINS")
	if !ok {
		pluginDir = clu.Viewer.Server.Plugins
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
		clu.Audit = func(started time.Time, req *http.Request, err error) {
			if HasPlugins() {
				SendAuditError(started, req, err)
			}
		}
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
		case NoPlugin:
		case AuditPlugin:
			symAudit, err := plug.Lookup("Audit")
			if err != nil {
				services.ServerMessage("Error opening Audit plugin %s Version: %s : %v", loader.Name(), loader.Version(), err)
			} else {
				services.ServerMessage("Audit plugin: %s Version: %s", loader.Name(), loader.Version())
				audit := symAudit.(Audit)
				auditPlugins[info.Name()] = &AuditLoader{loader, audit}
			}
		case AdabasPlugin:
			symAdabas, err := plug.Lookup("Adabas")
			if err != nil {
				services.ServerMessage("Error opening Adabas plugin %s Version: %s : %v", loader.Name(), loader.Version(), err)
			} else {
				services.ServerMessage("Adabas plugin: %s Version: %s", loader.Name(), loader.Version())
				adaSym := symAdabas.(Adabas)
				adabasPlugins[info.Name()] = &AdabasLoader{loader, adaSym}
			}
		case AuthPlugin:
			symCallback, err := plug.Lookup("Callback")
			if err != nil {
				services.ServerMessage("Error opening Authenticate/Authorize plugin %s Version: %s : %v", loader.Name(), loader.Version(), err)
			} else {
				services.ServerMessage("Authenticate/Authorize plugin: %s Version: %s", loader.Name(), loader.Version())
				symAuth := symCallback.(Auth)
				authPlugins[info.Name()] = &AuthLoader{loader, symAuth}
			}
		case ExtendPlugin:
			symExtend, err := plug.Lookup("EntryPoint")
			if err != nil {
				services.ServerMessage("Error opening Extend plugin %s Version: %s : %v", loader.Name(), loader.Version(), err)
			} else {
				services.ServerMessage("Extend plugin: %s Version: %s", loader.Name(), loader.Version())
				extendSym := symExtend.(server.RestExtend)
				extendPlugins[info.Name()] = &ExtendLoader{loader, extendSym}
				server.RegisterExtend(extendSym)
			}
		case ValidatorPlugin:
			symValidator, err := plug.Lookup("EntryPoint")
			if err != nil {
				services.ServerMessage("Error opening Validator plugin %s Version: %s : %v", loader.Name(), loader.Version(), err)
			} else {
				services.ServerMessage("Extend plugin: %s Version: %s", loader.Name(), loader.Version())
				validatorSym := symValidator.(server.RestValidator)
				validatorPlugins[info.Name()] = &ValidatorLoader{loader, validatorSym}
				server.RegisterValidator(validatorSym)
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
		if len(adabasPlugins) == 0 && len(auditPlugins) == 0 {
			return
		}
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

// LoginAudit login audit
func LoginAudit(method string, status string, session *auth.SessionInfo, user *auth.UserInfo) {
	for _, x := range auditPlugins {
		x.Audit.LoginAudit(method, status, session, user)
	}
}

// ReceiveAudit send audit information to plugins
func ReceiveAudit(p *clu.Context, r *http.Request) {
	if disablePlugin {
		return
	}
	log.Log.Debugf("Receive auditing plugins request: %v %v %v %T %p",
		r.Method, r.URL, server.RemoteHost(r), p, p)
	if strings.HasSuffix(r.RequestURI, "/version") {
		return
	}
	if p != nil {
		c := &http.Cookie{Name: p.UserName(), Value: p.UUID()}
		r.AddCookie(c)
		for _, x := range auditPlugins {
			x.Audit.ReceiveAudit(p.UserName(), p.UUID(), r)
		}
		return
	}
	for _, x := range auditPlugins {
		debug.PrintStack()
		log.Log.Fatal("Error clu context not defined")
		x.Audit.ReceiveAudit("Unknown", "-", r)
	}
}

// SendAuditEnded send audit information to plugins
func SendAuditEnded(started time.Time, r *http.Request) {
	if disablePlugin || r.Method == "OPTIONS" {
		return
	}
	if strings.HasSuffix(r.RequestURI, "/version") {
		return
	}
	log.Log.Debugf("Send auditing plugins request: %v %v %v", r.Method, r.URL, server.RemoteHost(r))
	user := "<Unknown>"
	uuid := "<UUID undefined>"
	reqToken := r.Header.Get("Authorization")
	if reqToken == "" {
		log.Log.Infof("Call without token: %v %v %v", r.Method, r.URL, server.RemoteHost(r))
		return
	}
	splitToken := strings.Split(reqToken, " ")
	switch strings.ToLower(splitToken[0]) {
	case "basic":
		user, _, _ = r.BasicAuth()
	case "bearer":
		if len(splitToken) > 1 {
			reqToken = strings.TrimSpace(splitToken[1])
			p, err := clu.Viewer.Server.WebToken.JWTContainsRoles(reqToken, []string{"admin", "user"})
			if err != nil {
				uuid = err.Error()
				log.Log.Errorf("Audit error: %v", err)
			} else {
				c := p.(*clu.Context)
				uuid = c.UUID()
				user = c.UserName()
			}
		}
	default:
		log.Log.Debugf("User evaluation failed in " + strings.ToLower(splitToken[0]))
	}

	log.Log.Debugf("Call sendaudit for user %s", user)
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
	if strings.HasSuffix(r.RequestURI, "/version") {
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

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
	"embed"
	"os"
	"path"
	"strconv"
	"sync"
	"time"

	"github.com/go-openapi/runtime/flagext"
	"github.com/tknie/errorrepo"
	"github.com/tknie/flynn/common"
	"github.com/tknie/log"
	"github.com/tknie/services"
	"github.com/tknie/services/auth"
	"gopkg.in/yaml.v3"
)

const (
	// DefaultConfigFileEnv default environment name to search configuration file at
	DefaultConfigFileEnv = "SERVER_CONFIG"
	// InstallationDirName directory name of product installation location
	InstallationDirName = "."
	// ServiceName servoce name of product
	ServiceName = "clutronapi"
)

// CurrentConfig current config name
var CurrentConfig = ""

//go:embed messages
var embedFiles embed.FS

//go:embed server/config.yaml
var embedConfig embed.FS

// ErrorType web request return type
type ErrorType byte

const (
	// Ok web response Ok
	Ok ErrorType = iota
	// BadRequest web response bad request
	BadRequest
)

// RestServer Rest server main node configuration structure containing
// all referenced parameters reflected to the configuration file.
type RestServer struct {
	Common       CommonConfig       `yaml:"rest-server"`
	Server       Server             `yaml:"server"`
	Database     DatabaseConfig     `yaml:"database"`
	Tasks        TaskConfig         `yaml:"tasks"`
	FileTransfer FileTransferConfig `yaml:"fileTransfer"`
	Metrics      []*Database
}

// CommonConfig server config base
type CommonConfig struct {
	Version             string `yaml:"version"`
	ConfigWatcher       bool   `yaml:"configWatcher,omitempty"`
	MaxBinaryBufferSize int    `yaml:"maxBinaryBufferSize,omitempty"`
	StatisticTimer      bool   `yaml:"statisticTimer,omitempty"`
	AppURL              string `yaml:"AppURL,omitempty"`
}

// Server REST server main configuration parameter
type Server struct {
	Service      []*Service        `yaml:"service"`
	WebToken     *auth.WebToken    `yaml:"JWT,omitempty"`
	LoginService LoginService      `yaml:"login"`
	Prefix       string            `yaml:"prefix,omitempty"`
	Content      string            `yaml:"content,omitempty"`
	Plugins      string            `yaml:"plugins,omitempty"`
	LogLocation  *services.Logging `yaml:"location"`
	Shutdown     struct {
		Passcode yaml.Node `yaml:"passcode,omitempty"`
	} `yaml:"shutdown"`
}

// DatabaseConfig database modelling and access
type DatabaseConfig struct {
	Mapping         Mapping        `yaml:"modelling"`
	DatabaseAccess  DatabaseAccess `yaml:"access"`
	SessionInfo     *SessionConfig `yaml:"sessionInfo"`
	UserInfo        *Database      `yaml:"userInfo"`
	BatchRepository *Database      `yaml:"batchRepository"`
}

// SessionConfig session configuration
type SessionConfig struct {
	DeleteUUID bool      `yaml:"deleteUUID"`
	Database   *Database `yaml:"database"`
}

// StatisticConfig statistics configuration
type StatisticConfig struct {
	Metrics *DatabaseAccess `yaml:"metrics"`
}

// Service service entry
type Service struct {
	Host              string           `yaml:"host"`
	Port              int              `yaml:"port"`
	Type              string           `yaml:"type"`
	Certificate       string           `yaml:"certificate,omitempty"`
	Key               string           `yaml:"key,omitempty"`
	MaxHeaderSize     flagext.ByteSize `yaml:"maxHeaderSize"`
	ReadTimeout       time.Duration    `yaml:"readTimeout"`
	WriteTimeout      time.Duration    `yaml:"writeTimeout"`
	KeepAlive         int              `yaml:"keepAive"`
	ListenLimit       int              `yaml:"listenLimit"`
	CleanupTimeout    time.Duration    `yaml:"cleanupTimeout"`
	TLSCACertificate  string           `yaml:"TLSCaCertificate"`
	TLSCertificate    string
	TLSCertificateKey string
}

// TaskConfig job store
type TaskConfig struct {
	Role      string    `yaml:"role,omitempty"`
	UseRole   bool      `yaml:"use_role,omitempty"`
	Directory string    `yaml:"directory,omitempty"`
	Database  *Database `yaml:"database,omitempty"`
}

// LoginService login service
type LoginService struct {
	Type                 string                       `yaml:"type,omitempty"`
	Module               string                       `yaml:"module,omitempty"`
	Administrators       string                       `yaml:"administrators,omitempty"`
	Users                string                       `yaml:"users,omitempty"`
	AuthenticationServer []*auth.AuthenticationServer `yaml:"authenticationServer,omitempty"`
}

// Mapping Adabas Maps
type Mapping struct {
	CacheUpdater int           `yaml:"cacheTimer,omitempty"`
	DatabaseMap  []DatabaseMap `yaml:"Modeling"`
}

// DatabaseMap database modelling configuration
type DatabaseMap struct {
	Name        string `yaml:"Name"`
	SrcDatabase string `yaml:"Database"`
	SQL         string `yaml:"SQL"`
	SrcTable    string `yaml:"SourceTable"`
	SrcField    string `yaml:"SourceField"`
	DestTable   string `yaml:"DestinationTable"`
	DestField   string `yaml:"DestinationField"`
}

// FileTransferConfig file transfer config
type FileTransferConfig struct {
	Admin       Admin `yaml:"Admin"`
	Directories struct {
		Role      string      `yaml:"role,omitempty"`
		UseRole   bool        `yaml:"use_role,omitempty"`
		Directory []Directory `yaml:"directory"`
	} `yaml:"directories"`
}

// Admin suite installation
type Admin struct {
	Role           string `yaml:"role,omitempty"`
	UseRole        bool   `yaml:"use_role,omitempty"`
	NoModification bool   `yaml:"NoModification,omitempty"`
}

// DatabaseAccess database access
type DatabaseAccess struct {
	Global   bool       `yaml:"global,omitempty"`
	Database []Database `yaml:"databases,omitempty"`
}

// Location location attribute
type Location struct {
	Location string `yaml:"location"`
}

// Directory directory entries
type Directory struct {
	Name     string `yaml:"name"`
	Location string `yaml:"location"`
}

// Database database
type Database struct {
	Driver               string   `yaml:"driver"`
	User                 string   `yaml:"user,omitempty"`
	Password             string   `yaml:"password,omitempty"`
	Target               string   `yaml:"target,omitempty"`
	Table                string   `yaml:"table,omitempty"`
	Tables               []string `yaml:"tables,omitempty"`
	Enabled              bool     `yaml:"enabled,omitempty"`
	AuthenticationGlobal bool     `yaml:"global_authentication,omitempty"`
}

var adadatadir string
var installation []string

// InitDirectAccess init Adabas direct access function, only need in active server
var InitDirectAccess func(*RestServer)

// InitAdmin init Adabas admin function, only need in active server
var InitAdmin func(*RestServer)

// Viewer containing server config
var Viewer *RestServer

var allCallbacks = make([]func(), 0)
var lock sync.Mutex
var loadedAlready = false

// RegisterConfigUpdates register configuration trigger function
func RegisterConfigUpdates(f func()) {
	log.Log.Debugf("Registry function")
	allCallbacks = append(allCallbacks, f)
	lock.Lock()
	defer lock.Unlock()
	if loadedAlready {
		f()
	}
}

// LoadedConfig triggered by configuration load
func LoadedConfig() {
	lock.Lock()
	defer lock.Unlock()

	for _, f := range allCallbacks {
		f()
	}
	services.ServerMessage("Load of configuration completed")
}

// String representation of Database instance
func (db *Database) String() string {
	log.Log.Debugf("Datbase target %s", db.Target)
	ref, p, err := common.NewReference(os.ExpandEnv(db.Target))
	if err != nil {
		log.Log.Debugf("Parse error target: %v", db.Target)
		return "<Error: " + err.Error() + ">"
	}

	if db.Password == "" {
		db.Password = p
	}
	port := strconv.Itoa(ref.Port)
	return db.User + ":***@" + ref.Host + ":" + port
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
		log.Log.Infof("No File location defined, file transfer not possible")
	} else {
		for _, d := range viewer.FileTransfer.Directories.Directory {
			if AddLocation != nil {
				AddLocation(d.Name, d.Location)
			}
		}
	}
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

// AddLocation add location, only needed in active server
var AddLocation = func(name, location string) error {
	if name != "" && os.ExpandEnv(location) != "" {
		services.ServerMessage("Add location %s at %s", name, location)
	}
	return nil
}

// LoadMessages load all REST server embed message templates
func LoadMessages() {
	fss, err := embedFiles.ReadDir("messages")
	if err != nil {
		panic("Internal config load error: " + err.Error())
	}
	for _, f := range fss {
		if f.Type().IsRegular() {
			byteValue, err := embedFiles.ReadFile("messages/" + f.Name())
			if err != nil {
				panic("Internal config load error: " + err.Error())
			}
			lang := path.Ext(f.Name())
			errorrepo.RegisterMessage(lang[1:], string(byteValue))
		}
	}
	// errorrepo.RegisterDirectory(fss)
}

// LoadConfig load xml configuration file
// The components are used to load and inject the configuration
func LoadConfig(watch bool, loaderInterface services.ConfigInterface) error {
	CurrentConfig = os.Getenv(DefaultConfigFileEnv)
	if CurrentConfig == "" {
		CurrentConfig = os.ExpandEnv("${SERVER_HOME}/configuration/config.yaml")
	}
	Viewer = &RestServer{}
	err := services.LoadConfig(CurrentConfig, loaderInterface, watch)
	if err != nil {
		services.ServerErrorMessage("RERR00042", err)
		/*if skipTemplate {
			return adaErr
		}*/
		services.ServerMessage("Loading config template (%v)", err)
		Viewer = loadConfigurationTemplate(loaderInterface)
		services.ServerMessage("Using embed template configuration")
	}
	adaptLogInstances()
	return nil
}

func adaptLogInstances() {
	// if log.Log != log.Log {
	// 	// log.Log = log.Log
	// 	log.SetDebugLevel(log.IsDebugLevel())
	// 	log.Log.Debugf("DEBUG: Testing log ....")
	// 	log.Log.Infof("INFO:  Testing log ....")
	// 	log.Log.Errorf("ERROR: Testing log ....")
	// 	log.Log.Debugf("DEBUG: Testing adatype log ....")
	// }
}

// loadConfigurationTemplate load default configuration available as
// embed file in the binary
func loadConfigurationTemplate(loaderInterface services.ConfigInterface) *RestServer {
	byteValue, err := embedConfig.ReadFile("config.yaml")
	if err != nil {
		panic("Internal config access error: " + err.Error())
	}
	viewer := &RestServer{}
	err = services.ParseConfig(byteValue, loaderInterface)
	if err != nil {
		panic("Internal config interpreter error: " + err.Error())
	}

	return viewer
}

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
	"embed"
	"sync"
	"time"

	"github.com/go-openapi/runtime/flagext"
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

var currentConfig = ""

//go:embed messages
var embedFiles embed.FS

//go:embed config.yaml
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
	Mapping        Mapping        `yaml:"modelling"`
	DatabaseAccess DatabaseAccess `yaml:"access"`
	SessionInfo    *Database      `yaml:"sessionInfo"`
	UserInfo       *Database      `yaml:"userInfo"`
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
	Database []Database `yaml:"Database,omitempty"`
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
	Driver string `yaml:"driver"`
	// URL      string `yaml:"url,omitempty"`
	Host     string `yaml:"host,omitempty"`
	Port     string `yaml:"port,omitempty"`
	User     string `yaml:"user,omitempty"`
	Password string `yaml:"password,omitempty"`
	Database string `yaml:"database,omitempty"`
	Table    string `yaml:"table,omitempty"`
	Enabled  bool   `yaml:"enabled,omitempty"`
}

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
	services.ServerMessage("Configuration loading completed")
}

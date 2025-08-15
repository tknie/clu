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

package main

import (
	"net/http"

	"github.com/tknie/clu/plugins"
	"github.com/tknie/log"
)

type greeting string

// This plugin can be used to provide extra information for
// validation to extra systems like ServiceNow or Jira tickets.

// Types type of plugin working with
func (g greeting) Types() []plugins.PluginTypes {
	return []plugins.PluginTypes{plugins.ValidatorPlugin}
}

// Name name of the plugin
func (g greeting) Name() string {
	return "Demo Validator plugin"
}

// Version version of the number
func (g greeting) Version() string {
	return "1.0"
}

// Stop stop plugin
func (g greeting) Stop() {
}

func (g greeting) EntryPoint() []string {
	return []string{"validator test"}
}

func (g greeting) CallValidator(req *http.Request) (_ bool, _ error) {
	log.Log.Debugf("Validate -> %s", req.URL.String())
	for k, v := range req.URL.Query() {
		log.Log.Debugf("%v => %v", k, v)
		if k == "validator" {
			return false, nil
		}
		if k == "xxx" {
			return false, nil
		}
	}
	return true, nil
}

// exported

// Loader loader for initialize plugin
var Loader greeting

// EntryPoint entry point for main structure
var EntryPoint greeting

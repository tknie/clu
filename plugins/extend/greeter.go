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

package main

import (
	"fmt"
	"net/http"

	"github.com/go-faster/jx"
	"github.com/tknie/clu/api"
	"github.com/tknie/clu/plugins"
)

type greeting string

// Types type of plugin working with
func (g greeting) Types() []plugins.PluginTypes {
	return []plugins.PluginTypes{plugins.ExtendPlugin}
}

// Name name of the plugin
func (g greeting) Name() string {
	return "Demo Rest Extend"
}

// Version version of the number
func (g greeting) Version() string {
	return "1.0"
}

// Stop stop plugin
func (g greeting) Stop() {
}

func (g greeting) EntryPoint() string {
	return "test"
}

func (g greeting) Call(path string, req *http.Request) (r api.CallExtendRes, _ error) {
	fmt.Println("Extend plugin call received:" + path)
	var data []api.ResponseRecordsItem
	d := make(api.ResponseRecordsItem)
	t := "XXX"
	s := "FFFF"
	raw := jx.Raw([]byte("\"" + t + "\""))
	d[s] = raw
	data = append(data, d)
	resp := &api.Response{Records: data, FieldNames: []string{s}}
	return resp, nil

}

// exported

// Loader loader for initialize plugin
var Loader greeting

// EntryPoint entry point for main structure
var EntryPoint greeting

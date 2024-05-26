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

	ht "github.com/ogen-go/ogen/http"

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

func (g greeting) CallGet(path string, req *http.Request) (r api.CallExtendRes, _ error) {
	fmt.Println("Extend plugin call received:" + path)
	d := make(api.ResponseRaw)
	t := "XXX"
	s := "FFFF"
	raw := jx.Raw([]byte("\"" + t + "\""))
	d[s] = raw
	var e1 jx.Encoder
	e1.SetIdent(0)
	e1.ObjStart()           // {
	e1.FieldStart("arrval") // "values":
	e1.ArrStart()           // [
	for _, v := range []int{4, 8, 15, 16, 23, 42} {
		e1.Int(v)
	}
	e1.ArrEnd() // ]
	e1.ObjEnd() // }
	d["l"] = e1.Bytes()
	var e2 jx.Encoder
	e2.Reset()
	e2.ObjStart()             // {
	e2.FieldStart("fieldval") // "values":
	e2.Str("XXXX")
	e2.ObjEnd()
	d["a"] = e2.Bytes()
	var e3 jx.Encoder
	e3.ArrStart() // [
	for _, v := range []int{4, 8, 15, 16, 23, 42} {
		e3.Int(v)
	}
	e3.ArrEnd() // ]
	d["x"] = e3.Bytes()
	return &d, nil

}

func (g greeting) CallPut(path string, req *http.Request) (r api.TriggerExtendRes, _ error) {
	return r, ht.ErrNotImplemented
}
func (g greeting) CallPost(path string, req *http.Request) (r api.CallPostExtendRes, _ error) {
	return r, ht.ErrNotImplemented
}

// exported

// Loader loader for initialize plugin
var Loader greeting

// EntryPoint entry point for main structure
var EntryPoint greeting

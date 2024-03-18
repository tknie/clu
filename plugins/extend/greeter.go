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

package main

import (
	"fmt"

	ht "github.com/ogen-go/ogen/http"

	"github.com/tknie/clu/api"
)

type greeting string

// Types type of plugin working with
func (g greeting) Types() []int {
	return []int{4}
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

func (g greeting) Call(path string) (r api.CallExtendRes, _ error) {
	fmt.Println("Extend plugin call received:" + path)
	return r, ht.ErrNotImplemented
}

// exported

// Loader loader for initialize plugin
var Loader greeting

// EntryPoint entry point for main structure
var EntryPoint greeting
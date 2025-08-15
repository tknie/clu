/*
* Copyright 2024-2025 Thorsten A. Knieling
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
	"net/http"
	"sync"

	"github.com/tknie/clu"
	"github.com/tknie/log"
	"github.com/tknie/services/auth"
)

// RestValidator Validator method to send to plugin
type RestValidator interface {
	EntryPoint() []string
	CallValidator(req *http.Request) (_ bool, _ error)
}

var validatorMap = sync.Map{}

// RegisterValidator register the validtor handler
func RegisterValidator(validator RestValidator) {
	for _, v := range validator.EntryPoint() {
		validatorMap.Store(v, validator)
	}
}

// Validate validate current HTTP session received in REST server
func Validate(session *clu.Context, role auth.AccessRole, resource string) bool {
	req := session.CurrentRequest
	validated := true
	validatorMap.Range(func(key, value any) bool {
		v := value.(RestValidator)
		validated, _ = v.CallValidator(req)
		return validated
	})
	if !validated {
		return false
	}
	writeAccess := false
	switch req.Method {
	case http.MethodGet:
		writeAccess = false
	case http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodPatch:
		writeAccess = false
	default:
	}
	if !auth.ValidUser(role, writeAccess, session.User(), resource) {
		log.Log.Debugf("Validate user forbidden")
		return false
	}
	return true
}

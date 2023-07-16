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
	"context"

	"github.com/tknie/clu"
	"github.com/tknie/clu/api"
	"github.com/tknie/log"
)

// GetLoginSession implements getLoginSession operation.
//
// Login receiving JWT.
//
// GET /login
func (ServerHandler) GetLoginSession(ctx context.Context) (r api.GetLoginSessionRes, _ error) {
	x, err := commonLogin(ctx)
	return x.(api.GetLoginSessionRes), err
}

func commonLogin(ctx context.Context) (interface{}, error) {
	session := ctx.(*clu.Context)
	token, err := Viewer.Server.WebToken.GenerateJWToken("*", session)
	if err != nil {
		log.Log.Errorf("Error token generation:%v", err)
		return &api.Error{Error: api.NewOptErrorError(api.ErrorError{Message: api.NewOptString(err.Error())})}, nil
	}

	return &api.AuthorizationToken{Token: api.NewOptString(token)}, nil
}

// LoginSession implements loginSession operation.
//
// Login receiving JWT.
//
// PUT /login
func (ServerHandler) LoginSession(ctx context.Context) (r api.LoginSessionRes, _ error) {
	x, err := commonLogin(ctx)
	return x.(api.LoginSessionRes), err
}

// PushLoginSession implements pushLoginSession operation.
//
// Login receiving JWT.
//
// POST /login
func (ServerHandler) PushLoginSession(ctx context.Context) (r api.PushLoginSessionRes, _ error) {
	x, err := commonLogin(ctx)
	return x.(api.PushLoginSessionRes), err
}

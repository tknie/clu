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
	"fmt"
	"time"

	"github.com/google/uuid"
	ht "github.com/ogen-go/ogen/http"
	"github.com/tknie/clu"
	"github.com/tknie/clu/api"
	"github.com/tknie/log"
	"github.com/tknie/services/auth"
)

// GetLoginSession implements getLoginSession operation.
//
// Login receiving JWT.
//
// GET /login
func (Handler) GetLoginSession(ctx context.Context) (r api.GetLoginSessionRes, _ error) {
	x, err := GenerateJWToken(ctx)
	return x.(api.GetLoginSessionRes), err
}

// GenerateJWToken generate JWT token on context
func GenerateJWToken(ctx context.Context) (interface{}, error) {
	session := ctx.(*clu.Context)
	log.Log.Debugf("Generate JWT token")
	if session.Token == "" {
		token, err := Viewer.Server.WebToken.GenerateJWToken("*", session)
		if err != nil {
			log.Log.Errorf("Error token generation:%v", err)
			return nil, err
		}
		session.Token = token
	}

	uuid, err := uuid.Parse(session.UUID())
	if err != nil {
		fmt.Println("Error parsing UUID", session.UUID())
	}
	t := api.AuthorizationToken{Token: api.NewOptString(session.Token),
		User: api.NewOptUser(api.User{LongName: api.NewOptString(session.User.LongName),
			Created:    api.NewOptDateTime(session.User.Created),
			LastLogin:  api.NewOptDateTime(session.User.LastLogin),
			UUID:       api.NewOptUUID(uuid),
			Name:       api.NewOptString(session.User.User),
			Permission: api.NewOptString(session.User.Permission),
			Email:      api.NewOptString(session.User.EMail)})}
	return &api.AuthorizationTokenHeaders{XToken: api.NewOptString(session.Token),
		Response: t}, nil
}

// LoginSession implements loginSession operation.
//
// Login receiving JWT.
//
// PUT /login
func (Handler) LoginSession(ctx context.Context) (r api.LoginSessionRes, _ error) {
	x, err := GenerateJWToken(ctx)
	return x.(api.LoginSessionRes), err
}

// PushLoginSession implements pushLoginSession operation.
//
// Login receiving JWT.
//
// POST /login
func (Handler) PushLoginSession(ctx context.Context) (r api.PushLoginSessionRes, _ error) {
	x, err := GenerateJWToken(ctx)
	return x.(api.PushLoginSessionRes), err
}

// RemoveSessionCompat implements removeSessionCompat operation.
//
// Remove the session.
//
// GET /logoff
func (Handler) RemoveSessionCompat(ctx context.Context) (r api.RemoveSessionCompatRes, _ error) {
	session := ctx.(*clu.Context)

	auth.InvalidateUUID(session.Token, time.Now())

	return &api.RemoveSessionCompatOK{}, nil
}

// LogoutSessionCompat implements logoutSessionCompat operation.
//
// Logout the session.
//
// PUT /logout
func (Handler) LogoutSessionCompat(ctx context.Context) (r api.LogoutSessionCompatRes, _ error) {
	return r, ht.ErrNotImplemented
}

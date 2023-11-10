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

package webserver

import (
	"context"
	"net/http"
	"strings"

	"github.com/tknie/clu"
	"github.com/tknie/clu/api"
	"github.com/tknie/clu/plugins"
	"github.com/tknie/clu/server"
	"github.com/tknie/log"
	"github.com/tknie/services/auth"
)

// SecurityHandler security handler
type SecurityHandler struct {
}

// HandleBasicAuth handle basic authorisation
func (sec SecurityHandler) HandleBasicAuth(ctx context.Context, operationName string, t api.BasicAuth) (context.Context, error) {
	username := strings.ToLower(strings.Trim(t.Username, " "))
	p, err := auth.BasicAuth(username, t.Password)
	if err != nil {
		log.Log.Errorf("Basic auth error... %v", err)
		return nil, err
	}
	log.Log.Infof("Basic auth... done: %s", p.Name())
	pm := p.(*clu.Context)
	server.GenerateJWToken(pm)
	// plugins.ReceiveAudit(nil, req)
	return pm, nil
}

// HandleBearerAuth handler Bearer authentication
func (sec SecurityHandler) HandleBearerAuth(ctx context.Context, operationName string, t api.BearerAuth) (context.Context, error) {
	// The header: Authorization: Bearer {base64 string} (or ?access_token={base 64 string} param) has already
	// been decoded by the runtime as a token
	p, err := server.Viewer.Server.WebToken.JWTContainsRoles(t.Token, []string{"admin"})
	if err != nil {
		if log.IsDebugLevel() {
			log.Log.Debugf("Bearer auth return: %v", err)
		}
		return nil, err
	}
	if log.IsDebugLevel() {
		log.Log.Debugf("Bearer request return %s", p.Name())
	}
	return p.(*clu.Context), nil
}

// HandleTokenCheck handler for Token to authentication
func (sec SecurityHandler) HandleTokenCheck(ctx context.Context, operationName string, t api.TokenCheck) (context.Context, error) {
	return &clu.Context{}, nil
}

// Request request secruity handler
func (sec SecurityHandler) Request(ctx context.Context, req *http.Request) {
	plugins.ReceiveAudit(ctx.(*clu.Context), req)
}

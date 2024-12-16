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

package plugins

import (
	"github.com/tknie/errorrepo"
	"github.com/tknie/services/auth"
)

// BasicAuth handle basic authentication
func BasicAuth(user string, pass string) (auth.PrincipalInterface, error) {
	sessionUUID := &auth.SessionInfo{}
	principal := auth.PrincipalCreater(sessionUUID, user, pass)
	for _, p := range authPlugins {
		err := p.Auth.Authenticate(principal, user, pass)
		if err == nil {
			return principal, nil
		}
	}
	return nil, errorrepo.NewError("RERR00010")
}

// HandleBearerAuth handle bearer authentication
func HandleBearerAuth(token string) (auth.PrincipalInterface, error) {
	return nil, errorrepo.NewError("RERR00011")
}

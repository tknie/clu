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
	"fmt"

	"github.com/tknie/services/auth"
)

// BasicAuth handle basic authorisation
func BasicAuth(user string, pass string) (auth.PrincipalInterface, error) {
	sessionUUID := &auth.SessionInfo{}
	principal := auth.PrincipalCreater(sessionUUID, user, pass)
	for _, p := range authPlugins {
		err := p.Auth.Authenticate(principal, user, pass)
		if err == nil {
			return principal, nil
		}
	}
	return nil, fmt.Errorf("no plugins")
}

// HandleBearerAuth handle bearer auth
func HandleBearerAuth(token string) (auth.PrincipalInterface, error) {
	return nil, fmt.Errorf("no plugins")
}

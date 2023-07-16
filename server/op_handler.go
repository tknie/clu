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
	"github.com/tknie/services"
	"github.com/tknie/services/auth"
)

// GetVersion implements getVersion operation.
//
// Retrieves the current version.
//
// GET /version
func (ServerHandler) GetVersion(ctx context.Context) (r api.GetVersionRes, _ error) {
	r = &api.Versions{Product: api.NewOptString("REST-API"),
		Version: api.NewOptString(services.BuildVersion)}
	return r, nil
}

// GetMaps implements getMaps operation.
//
// Retrieves a list of available maps.
//
// GET /rest/view
func (ServerHandler) GetMaps(ctx context.Context) (r api.GetMapsRes, _ error) {
	session := ctx.(*clu.Context)
	if !auth.ValidUser(auth.UserRole, false, session.User, "*Maps") {
		return &api.GetMapsForbidden{}, nil
	}

	maps := make([]api.Map, 0)
	for _, m := range GetAllViews() {
		maps = append(maps, api.Map(m))
	}
	r = &api.Maps{Maps: maps}
	return r, nil
}

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

package server

import (
	"context"
	"net/http"
	"time"

	"github.com/ogen-go/ogen/ogenerrors"
	"github.com/tknie/clu"
	"github.com/tknie/clu/api"
	"github.com/tknie/errorrepo"
	"github.com/tknie/log"
)

// NewAPIError new API error
func NewAPIError(code string, err error) *api.Error {
	return &api.Error{Code: api.NewOptString(code),
		Error: api.NewOptErrorError(api.ErrorError{Message: api.NewOptString(err.Error())})}
}

// NewError creates *ErrorStatusCode from error returned by handler.
//
// Used for common default response.
func (Handler) NewError(ctx context.Context, err error) (r *api.ErrorStatusCode) {
	r = new(api.ErrorStatusCode)
	r.StatusCode = http.StatusServiceUnavailable
	switch session := ctx.(type) {
	case *clu.Context:
		log.Log.Errorf("Server handler error: %v/%s -> status=%d", err, session.UserName(), r.StatusCode)
		session.SendAuditError(time.Now(), err)
	default:
		log.Log.Errorf("Unknown error context: %T", ctx)
	}
	switch e := err.(type) {
	case *ogenerrors.SecurityError:
		r.StatusCode = http.StatusForbidden
		r.Response = *NewAPIError("SECERR", err)
		return r
	case *api.ErrorStatusCode:
		return e
	case *errorrepo.Error:
		r.Response = *NewAPIError(e.ID(), e)
	default:
		log.Log.Errorf("Unknown error type %T", e)
		r.StatusCode = http.StatusBadGateway
		r.Response = *NewAPIError("UNKERR", err)
	}
	return r
}

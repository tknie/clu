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

	ht "github.com/ogen-go/ogen/http"
	"github.com/tknie/clu/api"
)

// BatchQuery implements batchQuery operation.
//
// Call a SQL query batch command posted in body.
//
// POST /rest/batch
func (Handler) BatchQuery(ctx context.Context, req api.OptSQLQuery) (r api.BatchQueryRes, _ error) {

	return r, ht.ErrNotImplemented
}

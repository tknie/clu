package server

import (
	"context"
	"fmt"

	ht "github.com/ogen-go/ogen/http"
	"github.com/tknie/clu"
	"github.com/tknie/clu/api"
	"github.com/tknie/flynn/common"
	"github.com/tknie/log"
	"github.com/tknie/services/auth"
)

// InsertRecord implements insertRecord operation.
//
// Insert given record.
//
// POST /rest/view/{table}
func (ServerHandler) InsertRecord(ctx context.Context, params api.InsertRecordParams) (r api.InsertRecordRes, _ error) {
	return r, ht.ErrNotImplemented
}

// DeleteRecordsSearched implements deleteRecordsSearched operation.
//
// Delete a record with a given search.
//
// DELETE /rest/view/{table}/{search}
func (ServerHandler) DeleteRecordsSearched(ctx context.Context, params api.DeleteRecordsSearchedParams) (r api.DeleteRecordsSearchedRes, _ error) {
	session := ctx.(*clu.Context)
	log.Log.Debugf("Search records for fields %s -> %s", session.User, params.Table)
	if !auth.ValidUser(auth.UserRole, false, session.User, params.Table) {
		return &api.DeleteRecordsSearchedForbidden{}, nil
	}
	log.Log.Debugf("SQL search fields %s - %v", params.Table, params.Search)
	d, err := ConnectTable(session, params.Table)
	if err != nil {
		log.Log.Errorf("Error search table %s:%v", params.Table, err)
		return &api.Error{Error: api.NewOptErrorError(api.ErrorError{Message: api.NewOptString(err.Error())})}, nil
	}
	defer CloseTable(d)
	dr, err := d.Delete(params.Table, &common.Entries{})
	fmt.Println("DR", dr)
	return r, ht.ErrNotImplemented
}

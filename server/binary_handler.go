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
	"io"
	"strings"

	"github.com/tknie/clu"
	"github.com/tknie/clu/api"
	"github.com/tknie/flynn/common"
	"github.com/tknie/log"
	"github.com/tknie/services/auth"
)

const blockSize = 65536

type streamRead struct {
	principal *clu.Context
	field     string
	mimetype  string
	data      []byte
	send      int
}

// GetImage implements getImage operation.
//
// Retrieves a field of a specific ISN of a Map definition.
//
// GET /image/{table}/{field}/{search}
func (ServerHandler) GetImage(ctx context.Context, params api.GetImageParams) (r api.GetImageRes, _ error) {
	log.Log.Debugf("GET IMAGE ...")
	session := ctx.(*clu.Context)
	if !auth.ValidUser(auth.UserRole, false, session.User, params.Table) {
		return &api.GetImageForbidden{}, nil
	}

	log.Log.Debugf("SQL image search", params.Table, params.Field, params.Search)
	read, err := initStream(session, params.Table,
		params.Field, params.Search, "")
	if err != nil {
		log.Log.Errorf("Error search table %s:%v", params.Table, err)
		return &api.Error{Error: api.NewOptErrorError(api.ErrorError{Message: api.NewOptString(err.Error())})}, nil
	}
	read.mimetype = "image/jpeg"
	read.field = params.Field
	reader, err := read.streamResponderFunc()
	if err != nil {
		log.Log.Errorf("Error search table %s:%v", params.Table, err)
		return &api.Error{Error: api.NewOptErrorError(api.ErrorError{Message: api.NewOptString(err.Error())})}, nil
	}
	r = &api.GetImageOKImageGIF{Data: reader}
	log.Log.Debugf("Return IMAGE: %#v\n", r)
	return r, nil
}

// GetVideo implements getVideo operation.
//
// Retrieves a video stream of a specific ISN of a Map definition.
//
// GET /video/{table}/{field}/{search}
func (ServerHandler) GetVideo(ctx context.Context, params api.GetVideoParams) (r api.GetVideoRes, _ error) {
	log.Log.Debugf("GET Video ...")
	session := ctx.(*clu.Context)
	if !auth.ValidUser(auth.UserRole, false, session.User, params.Table) {
		return &api.GetVideoForbidden{}, nil
	}
	log.Log.Debugf("SQL video table=%s field=%s search=%s", params.Table, params.Field, params.Search)
	read, err := initStream(session, params.Table,
		params.Field, params.Search, params.MimetypeField)
	if err != nil {
		log.Log.Errorf("Error search table %s:%v", params.Table, err)
		return &api.Error{Error: api.NewOptErrorError(api.ErrorError{Message: api.NewOptString(err.Error())})}, nil
	}
	read.mimetype = params.MimetypeField
	read.field = params.Field
	reader, err := read.streamResponderFunc()
	if err != nil {
		log.Log.Errorf("Error search table %s:%v", params.Table, err)
		return &api.Error{Error: api.NewOptErrorError(api.ErrorError{Message: api.NewOptString(err.Error())})}, nil
	}
	if read.mimetype == "video/mp3" || read.mimetype == "" {
		r = &api.GetVideoOKVideoMP4{Data: reader}
	} else {
		r = &api.GetVideoOKVideoMov{Data: reader}
	}
	log.Log.Debugf("Return VIDEO: %#v\n", r)
	return r, nil
}

// GetLobByMap implements getLobByMap operation.
//
// Retrieves a lob of a specific ISN of an field in a Map.
//
// GET /binary/{table}/{field}/{search}
func (ServerHandler) GetLobByMap(ctx context.Context, params api.GetLobByMapParams) (r api.GetLobByMapRes, _ error) {
	log.Log.Debugf("GET LOB ...")
	session := ctx.(*clu.Context)
	if !auth.ValidUser(auth.UserRole, false, session.User, params.Table) {
		return &api.GetLobByMapForbidden{}, nil
	}

	log.Log.Debugf("SQL image search", params.Table, params.Field, params.Search)
	read, err := initStream(session, params.Table,
		params.Field, params.Search, "")
	if err != nil {
		log.Log.Errorf("Error search table %s:%v", params.Table, err)
		return &api.Error{Error: api.NewOptErrorError(api.ErrorError{Message: api.NewOptString(err.Error())})}, nil
	}
	read.mimetype = "application/octet-stream"
	read.field = params.Field
	reader, err := read.streamResponderFunc()
	if err != nil {
		log.Log.Errorf("Error search table %s:%v", params.Table, err)
		return &api.Error{Error: api.NewOptErrorError(api.ErrorError{Message: api.NewOptString(err.Error())})}, nil
	}
	r = &api.GetLobByMapOK{Data: reader}
	log.Log.Debugf("Return LOB: %#v\n", r)
	return r, nil
}

// initStream init streamed read and open connection to database
func initStream(srvctx *clu.Context, table,
	field, search, mimetypeField string) (read *streamRead, err error) {
	d, err := ConnectTable(srvctx, table)
	if err != nil {
		log.Log.Errorf("Error search table %s:%v", table, err)
		return nil, err
	}
	defer CloseTable(d)

	fields := []string{strings.ToLower(field)}
	if mimetypeField != "" {
		fields = []string{strings.ToLower(field), mimetypeField}
	}

	q := &common.Query{TableName: table,
		Fields: fields,
		Search: search}
	result, err := queryBytes(d, q)
	if err != nil {
		log.Log.Errorf("Error query table %s:%v", table, err)
		return nil, err
	}

	s := strings.ToLower(field)
	if d, ok := result[s]; ok {
		if d == nil {
			return nil, fmt.Errorf("internal error result map nil")
		}
		read = &streamRead{data: d.([]byte)}
		if mimetypeField != "" {
			s := strings.ToLower(mimetypeField)
			if d, ok := result[s]; ok {
				read.mimetype = d.(string)
			}
		}
		return
	}
	err = fmt.Errorf("field not in result map")
	return
}

// nextDataBlock read partial large object segment
func (read *streamRead) nextDataBlock() (data []byte, err error) {
	if read.send >= len(read.data) {
		return nil, io.EOF
	}
	sz := read.send + blockSize
	if sz > len(read.data) {
		sz = len(read.data)
	}
	data = read.data[read.send:sz]
	if err != nil {
		log.Log.Errorf("Error read LOB segment", err)
		return
	}
	read.send += blockSize
	return
}

// streamResponderFunc function to response and send read data
func (read *streamRead) streamResponderFunc() (io.Reader, error) {
	reader, writer := io.Pipe()
	errChan := make(chan error)
	go func() {
		count := 0
		for {
			data, err := read.nextDataBlock()
			if err != nil {
				log.Log.Debugf("Got data block error: %v", err)
				errChan <- err
				writer.CloseWithError(err)
				return
			}
			log.Log.Debugf("Got data block %d", len(data))
			if count == 0 {
				errChan <- nil
			}
			count += len(data)
			_, err = writer.Write(data)
			if err != nil {
				log.Log.Debugf("Got write data block error: %v", err)
				// read.errOut <- err
				writer.CloseWithError(err)
				return
			}
			if len(data) < blockSize {
				log.Log.Debugf("Got end data block %d/%d", len(data), count)

				writer.Close()
				return
			}
			log.Log.Debugf("Need next data block %d/%d", len(data), count)
		}
	}()
	err := <-errChan
	if err != nil {
		return nil, err
	}
	// handle err
	return reader, err
}

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

	"github.com/tknie/clu"
	"github.com/tknie/clu/api"
	"github.com/tknie/log"
	"github.com/tknie/services/auth"
)

// GetImage implements getImage operation.
//
// Retrieves a field of a specific ISN of a Map definition.
//
// GET /image/{table}/{field}/{search}
func (Handler) GetImage(ctx context.Context, params api.GetImageParams) (r api.GetImageRes, _ error) {
	log.Log.Debugf("GET IMAGE ...")
	session := ctx.(*clu.Context)
	if !auth.ValidUser(auth.UserRole, false, session.User, params.Table) {
		return &api.GetImageForbidden{}, nil
	}

	mimeTypeField := ""
	if params.MimetypeField != "" {
		mimeTypeField = params.MimetypeField
	}
	mimeType := ""
	if params.Mimetype.Set {
		mimeType = params.Mimetype.Value
	}

	log.Log.Debugf("SQL image search table=%s field=%s search=%s", params.Table, params.Field, params.Search)
	read := NewStreamRead(params.Table, params.Field, mimeTypeField)
	err := read.initStreamFromTable(session, params.Search, mimeType)
	if err != nil {
		log.Log.Errorf("Error search table %s:%v", params.Table, err)
		return nil, err
	}
	if read.mimetype == "" {
		read.mimetype = "image/jpeg"
	}
	read.field = params.Field
	reader, err := read.streamResponderFunc()
	if err != nil {
		log.Log.Errorf("Error search table %s:%v", params.Table, err)
		return nil, err
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
func (Handler) GetVideo(ctx context.Context, params api.GetVideoParams) (r api.GetVideoRes, _ error) {
	log.Log.Debugf("GET Video ...")
	session := ctx.(*clu.Context)
	if !auth.ValidUser(auth.UserRole, false, session.User, params.Table) {
		return &api.GetVideoForbidden{}, nil
	}
	log.Log.Debugf("SQL video table=%s field=%s search=%s", params.Table, params.Field, params.Search)
	read := NewStreamRead(params.Table, params.Field, params.MimetypeField)
	err := read.initStreamFromTable(session, params.Search, params.Mimetype)
	if err != nil {
		log.Log.Errorf("Error search table %s:%v", params.Table, err)
		return nil, err
	}
	read.mimetype = params.MimetypeField
	read.field = params.Field
	reader, err := read.streamResponderFunc()
	if err != nil {
		log.Log.Errorf("Error search table %s:%v", params.Table, err)
		return nil, err
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
func (Handler) GetLobByMap(ctx context.Context, params api.GetLobByMapParams) (r api.GetLobByMapRes, _ error) {
	log.Log.Debugf("GET LOB ...")
	session := ctx.(*clu.Context)
	if !auth.ValidUser(auth.UserRole, false, session.User, params.Table) {
		return &api.GetLobByMapForbidden{}, nil
	}

	mimeTypeField := ""
	if params.MimetypeField != "" {
		mimeTypeField = params.MimetypeField
	}
	mimeType := ""
	if params.Mimetype.Set {
		mimeType = params.Mimetype.Value
	}
	log.Log.Debugf("SQL image search", params.Table, params.Field, params.Search)
	read := NewStreamRead(params.Table, params.Field, mimeTypeField)
	err := read.initStreamFromTable(session, params.Search, mimeType)
	if err != nil {
		log.Log.Errorf("Error search table %s:%v", params.Table, err)
		return nil, err
	}
	read.mimetype = "application/octet-stream"
	read.field = params.Field
	reader, err := read.streamResponderFunc()
	if err != nil {
		log.Log.Errorf("Error search table %s:%v", params.Table, err)
		return nil, err
	}
	r = &api.GetLobByMapOK{Data: reader}
	log.Log.Debugf("Return LOB: %#v\n", r)
	return r, nil
}

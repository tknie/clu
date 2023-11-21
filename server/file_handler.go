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
	"os"
	"path/filepath"
	"runtime"
	"strings"

	ht "github.com/ogen-go/ogen/http"
	"github.com/tknie/clu"
	"github.com/tknie/clu/api"
	"github.com/tknie/log"
	"github.com/tknie/services/auth"
)

// BrowseList implements browseList operation.
//
// Retrieves a list of Browseable locations.
//
// GET /rest/file/browse
func (Handler) BrowseList(ctx context.Context) (r api.BrowseListRes, _ error) {
	session := ctx.(*clu.Context)
	if !auth.ValidUser(auth.AdministratorRole, false, session.User, "") {
		return &api.BrowseListForbidden{}, nil
	}
	d := &api.Directories{}
	for _, bd := range Viewer.FileTransfer.Directories.Directory {
		dbd := api.Directory{Location: api.NewOptString(bd.Location),
			Name: api.NewOptString(bd.Name)}
		d.Directories = append(d.Directories, dbd)
	}
	return d, nil
}

// CreateDirectory implements createDirectory operation.
//
// Create a new directory.
//
// PUT /rest/file/{path}
func (Handler) CreateDirectory(ctx context.Context, params api.CreateDirectoryParams) (r api.CreateDirectoryRes, _ error) {
	session := ctx.(*clu.Context)
	if !auth.ValidUser(auth.AdministratorRole, false, session.User, "") {
		return &api.CreateDirectoryForbidden{}, nil
	}
	if !auth.ValidUser(auth.AdministratorRole, false, session.User, "") {
		return &api.CreateDirectoryForbidden{}, nil
	}
	d, path, err := extraceLocationPath(params.Path)
	if err != nil {
		err := fmt.Errorf("location reference missing")
		return &api.CreateDirectoryNotFound{Error: api.NewOptErrorError(api.ErrorError{Message: api.NewOptString(err.Error())})}, nil

	}
	fileName := os.ExpandEnv(d.Location + "/" + path)
	log.Log.Debugf("FileName %s", fileName)
	_, ferr := os.Stat(fileName)
	if ferr != nil {
		if !os.IsNotExist(ferr) {
			err := fmt.Errorf("error opening location %s: %v", d.Name, ferr)
			return &api.CreateDirectoryNotFound{Error: api.NewOptErrorError(api.ErrorError{Message: api.NewOptString(err.Error())})}, nil
		}
		os.Mkdir(fileName, os.ModePerm)
		return &api.StatusResponse{Status: api.NewOptStatusResponseStatus(api.StatusResponseStatus{})}, nil
	}
	err = fmt.Errorf("Directory/File already exists")
	return &api.CreateDirectoryBadRequest{Error: api.NewOptErrorError(api.ErrorError{Message: api.NewOptString(err.Error())})}, nil
}

// DeleteFileLocation implements deleteFileLocation operation.
//
// Delete the file on the given location.
//
// DELETE /rest/file/{path}
func (Handler) DeleteFileLocation(ctx context.Context, params api.DeleteFileLocationParams) (r api.DeleteFileLocationRes, _ error) {
	session := ctx.(*clu.Context)
	if !auth.ValidUser(auth.AdministratorRole, false, session.User, "") {
		return &api.DeleteFileLocationForbidden{}, nil
	}
	d, path, err := extraceLocationPath(params.Path)
	if err != nil {
		err := fmt.Errorf("location reference missing")
		return &api.DeleteFileLocationNotFound{Error: api.NewOptErrorError(api.ErrorError{Message: api.NewOptString(err.Error())})}, nil

	}
	fmt.Printf("Try deleting location=%s path=%s\n", d.Location, path)

	return r, ht.ErrNotImplemented
}

// DownloadFile implements downloadFile operation.
//
// Download a file out of file location.
//
// GET /rest/file/{path}
func (Handler) DownloadFile(ctx context.Context, params api.DownloadFileParams) (r api.DownloadFileRes, _ error) {
	session := ctx.(*clu.Context)
	if !auth.ValidUser(auth.AdministratorRole, false, session.User, "") {
		return &api.DownloadFileForbidden{}, nil
	}
	d, path, err := extraceLocationPath(params.Path)
	if err != nil {
		err := fmt.Errorf("location reference missing")
		return &api.DownloadFileNotFound{Error: api.NewOptErrorError(api.ErrorError{Message: api.NewOptString(err.Error())})}, nil

	}
	fmt.Printf("Try download location=%s path=%s\n", d.Location, path)
	fileName := os.ExpandEnv(d.Location + "/" + path)
	log.Log.Debugf("FileName %s", fileName)
	f, ferr := os.Open(fileName)
	if ferr != nil {
		err := fmt.Errorf("error opening location %s", d.Name)
		return &api.DownloadFileNotFound{Error: api.NewOptErrorError(api.ErrorError{Message: api.NewOptString(err.Error())})}, nil
	}
	fileInfo, fierr := f.Stat()
	if fierr != nil {
		err := fmt.Errorf("error opening statistics of location %s", d.Name)
		return &api.DownloadFileNotFound{Error: api.NewOptErrorError(api.ErrorError{Message: api.NewOptString(err.Error())})}, nil
	}
	if fileInfo.IsDir() {
		err := fmt.Errorf("cannot download directory of location %s", d.Name)
		return &api.DownloadFileNotFound{Error: api.NewOptErrorError(api.ErrorError{Message: api.NewOptString(err.Error())})}, nil
	}
	read, err := initStreamFromFile(f)
	if err != nil {
		log.Log.Errorf("Error download file %s:%v", d.Location, err)
		return &api.DownloadFileBadRequest{Error: api.NewOptErrorError(api.ErrorError{Message: api.NewOptString(err.Error())})}, nil
	}
	read.mimetype = "application/octet-stream"
	reader, err := read.streamResponderFunc()
	if err != nil {
		log.Log.Errorf("Error init stream for file %s:%v", d.Location, err)
		return &api.DownloadFileBadRequest{Error: api.NewOptErrorError(api.ErrorError{Message: api.NewOptString(err.Error())})}, nil
	}
	ok := &api.DownloadFileOK{Data: reader}
	return ok, nil
}

func extraceLocationPath(paramsPath string) (*Directory, string, error) {
	location := filepath.Clean(paramsPath)
	log.Log.Debugf("Location %s", location)
	firstSlash := strings.IndexByte(location, '/')
	path := "/"
	if firstSlash == 0 {
		err := fmt.Errorf("location reference empty")
		return nil, "", err
	}
	if firstSlash > 0 {
		location = paramsPath[:firstSlash]
		path = paramsPath[firstSlash:]
	}
	log.Log.Debugf("location=%s path=%s", location, path)
	for _, d := range Viewer.FileTransfer.Directories.Directory {
		if d.Name == location {
			return &d, path, nil
		}
	}
	return nil, "", fmt.Errorf("location %s reference missing", location)
}

// BrowseLocation implements browseLocation operation.
//
// Retrieves a list of files in the defined location.
//
// GET /rest/file/browse/{path}
func (Handler) BrowseLocation(ctx context.Context, params api.BrowseLocationParams) (r api.BrowseLocationRes, _ error) {
	session := ctx.(*clu.Context)
	if !auth.ValidUser(auth.AdministratorRole, false, session.User, "") {
		return &api.BrowseLocationForbidden{}, nil
	}
	d, path, err := extraceLocationPath(params.Path)
	if err != nil {
		return &api.BrowseLocationNotFound{Error: api.NewOptErrorError(api.ErrorError{Message: api.NewOptString(err.Error())})}, nil
	}
	fileName := os.ExpandEnv(d.Location + "/" + path)
	log.Log.Debugf("FileName %s Filter %s", fileName, params.Filter.Value)
	f, ferr := os.Open(fileName)
	if ferr != nil {
		err := fmt.Errorf("error opening location %s", d.Name)
		return &api.BrowseLocationNotFound{Error: api.NewOptErrorError(api.ErrorError{Message: api.NewOptString(err.Error())})}, nil
	}
	fileInfo, fierr := f.Stat()
	if fierr != nil {
		err := fmt.Errorf("error opening statistics of location %s", d.Name)
		return &api.BrowseLocationNotFound{Error: api.NewOptErrorError(api.ErrorError{Message: api.NewOptString(err.Error())})}, nil
	}
	if fileInfo.IsDir() {
		return returnDirectoryInfo(d, path, params.Filter.Value, f)
	}
	return returnFileStream(d, path, f)
}

// returnDirectoryInfo generate directory information list
func returnDirectoryInfo(d *Directory, path, pattern string, f *os.File) (api.BrowseLocationRes, error) {
	files, err := f.ReadDir(0)
	if err != nil {
		err := fmt.Errorf("error reading path of location %s", d.Name)
		return &api.BrowseLocationNotFound{Error: api.NewOptErrorError(api.ErrorError{Message: api.NewOptString(err.Error())})}, nil
	}
	fl := &api.DirectoryFiles{Location: api.NewOptString(d.Location),
		Path:   api.NewOptString(path),
		Files:  make([]api.File, 0),
		System: api.NewOptString(runtime.GOOS)}

	for _, f := range files {
		b := true
		if pattern != "" {
			b, err = filepath.Match(pattern, f.Name())
			if err != nil {
				return nil, err
			}
		}
		if !strings.HasPrefix(f.Name(), ".") && b {
			fileType := "File"
			if f.IsDir() {
				fileType = "Directory"
			}
			fi, err := f.Info()
			if err != nil {
				return nil, err
			}

			content := api.File{Name: api.NewOptString(f.Name()), Type: api.NewOptString(fileType),
				Modified: api.NewOptDateTime(fi.ModTime()), Size: api.NewOptInt64(fi.Size())}
			fl.Files = append(fl.Files, content)
		}
	}
	return fl, nil
}

// returnFileStream return stream from a file to the corresponding HTTP response
func returnFileStream(d *Directory, path string, f *os.File) (api.BrowseLocationRes, error) {
	read, err := initStreamFromFile(f)
	if err != nil {
		log.Log.Errorf("Error download file %s:%v", d.Location, err)
		return &api.BrowseLocationBadRequest{Error: api.NewOptErrorError(api.ErrorError{Message: api.NewOptString(err.Error())})}, nil
	}
	read.mimetype = "application/octet-stream"
	reader, err := read.streamResponderFunc()
	if err != nil {
		log.Log.Errorf("Error init stream for file %s:%v", d.Location, err)
		return &api.BrowseLocationBadRequest{Error: api.NewOptErrorError(api.ErrorError{Message: api.NewOptString(err.Error())})}, nil
	}
	//ok := &api.BrowseLocationOKApplicationOctetStream{Data: reader}
	//reader.Read()
	ok2 := &api.BrowseLocationOKMultipartFormData{File: ht.MultipartFile{Name: f.Name(), File: reader}}
	fmt.Println("ok2", ok2.File.Name)
	return ok2, nil
}

// UploadFile implements uploadFile operation.
//
// Upload a new file to the given location.
//
// POST /rest/file/{location}
func (Handler) UploadFile(ctx context.Context, req *api.UploadFileReq, params api.UploadFileParams) (r api.UploadFileRes, _ error) {
	return r, ht.ErrNotImplemented
}

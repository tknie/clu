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
	"errors"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/tknie/clu"
	"github.com/tknie/clu/api"
	"github.com/tknie/errorrepo"
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
	d := &api.Directories{}
	for _, bd := range clu.Viewer.FileTransfer.Directories.Directory {
		dbd := api.Directory{Location: api.NewOptString(bd.Location),
			Name: api.NewOptString(bd.Name)}
		if !auth.ValidUser(auth.UserRole, false, session.User(), "<"+dbd.Name.Value) {
			return &api.BrowseListForbidden{}, nil
		}
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
	d, path, err := extraceLocationPath(params.Path)
	if err != nil {
		err := errorrepo.NewError("REST00100", params.Path)
		return &api.CreateDirectoryNotFound{Error: api.NewOptErrorError(api.ErrorError{Message: api.NewOptString(err.Error())})}, nil

	}
	session := ctx.(*clu.Context)
	if !auth.ValidUser(auth.UserRole, true, session.User(), ">"+d.Name) {
		return &api.CreateDirectoryForbidden{}, nil
	}
	fileName := os.ExpandEnv(d.Location + "/" + path)
	log.Log.Debugf("FileName %s", fileName)
	_, ferr := os.Stat(fileName)
	if ferr != nil {
		if !os.IsNotExist(ferr) {
			err := errorrepo.NewError("REST00101", d.Name, ferr)
			return &api.CreateDirectoryNotFound{Error: api.NewOptErrorError(api.ErrorError{Message: api.NewOptString(err.Error())})}, nil
		}
		os.Mkdir(fileName, os.ModePerm)
		return &api.StatusResponse{Status: api.NewOptStatusResponseStatus(api.StatusResponseStatus{})}, nil
	}
	err = errorrepo.NewError("REST00102", fileName)
	return &api.CreateDirectoryBadRequest{Error: api.NewOptErrorError(api.ErrorError{Message: api.NewOptString(err.Error())})}, nil
}

// DeleteFileLocation implements deleteFileLocation operation.
//
// Delete the file on the given location.
//
// DELETE /rest/file/{path}
func (Handler) DeleteFileLocation(ctx context.Context, params api.DeleteFileLocationParams) (r api.DeleteFileLocationRes, _ error) {
	d, path, err := extraceLocationPath(params.Path)
	if err != nil {
		err := errorrepo.NewError("REST00103")
		return &api.DeleteFileLocationNotFound{Error: api.NewOptErrorError(api.ErrorError{Message: api.NewOptString(err.Error())})}, nil

	}
	session := ctx.(*clu.Context)
	if !auth.ValidUser(auth.UserRole, true, session.User(), ">"+d.Name) {
		return &api.DeleteFileLocationForbidden{}, nil
	}
	fileName := os.ExpandEnv(d.Location + "/" + path)
	if params.File.IsSet() {
		fileName += params.File.Value
	}
	fileName = filepath.Clean(fileName)

	_, err = os.Stat(fileName)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			err := errorrepo.NewError("REST00104", fileName)
			return &api.DeleteFileLocationNotFound{Error: api.NewOptErrorError(api.ErrorError{Message: api.NewOptString(err.Error())})}, nil
		}
		err := errorrepo.NewError("REST00105", fileName, err)
		return &api.DeleteFileLocationNotFound{Error: api.NewOptErrorError(api.ErrorError{Message: api.NewOptString(err.Error())})}, nil
	}
	log.Log.Debugf("Try deleting location=%s file=%s\n", d.Location, fileName)

	err = os.Remove(fileName)
	if err != nil {
		err := errorrepo.NewError("REST00106", fileName, err)
		return &api.DeleteFileLocationNotFound{Error: api.NewOptErrorError(api.ErrorError{Message: api.NewOptString(err.Error())})}, nil
	}

	v := api.StatusResponseStatus{Message: api.NewOptString("deleted")}
	status := &api.StatusResponse{Status: api.NewOptStatusResponseStatus(v)}
	return status, nil
}

// DownloadFile implements downloadFile operation.
//
// Download a file out of file location.
//
// GET /rest/file/{path}
func (Handler) DownloadFile(ctx context.Context, params api.DownloadFileParams) (r api.DownloadFileRes, _ error) {
	d, path, err := extraceLocationPath(params.Path)
	if err != nil {
		err := errorrepo.NewError("REST00107", params.Path, err)
		return &api.DownloadFileNotFound{Error: api.NewOptErrorError(api.ErrorError{Message: api.NewOptString(err.Error())})}, nil

	}
	session := ctx.(*clu.Context)
	if !auth.ValidUser(auth.UserRole, false, session.User(), "<"+d.Name) {
		return &api.DownloadFileForbidden{}, nil
	}
	log.Log.Debugf("Try download location=%s path=%s", d.Location, path)
	fileName := os.ExpandEnv(d.Location + "/" + path)
	log.Log.Debugf("FileName %s", fileName)
	f, ferr := os.Open(fileName)
	if ferr != nil {
		err := errorrepo.NewError("REST00108", d.Name, ferr)
		return &api.DownloadFileNotFound{Error: api.NewOptErrorError(api.ErrorError{Message: api.NewOptString(err.Error())})}, nil
	}
	fileInfo, fierr := f.Stat()
	if fierr != nil {
		err := errorrepo.NewError("REST00109", d.Name, fierr)
		return &api.DownloadFileNotFound{Error: api.NewOptErrorError(api.ErrorError{Message: api.NewOptString(err.Error())})}, nil
	}
	if fileInfo.IsDir() {
		err := errorrepo.NewError("REST00110", d.Name)
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

func extraceLocationPath(paramsPath string) (*clu.Directory, string, error) {
	location := filepath.Clean(paramsPath)
	log.Log.Debugf("Location %s", location)
	firstSlash := strings.IndexByte(location, '/')
	path := "/"
	if firstSlash == 0 {
		err := errorrepo.NewError("REST00111")
		return nil, "", err
	}
	if firstSlash > 0 {
		location = paramsPath[:firstSlash]
		path = paramsPath[firstSlash:]
	}
	log.Log.Debugf("location=%s path=%s", location, path)
	for _, d := range clu.Viewer.FileTransfer.Directories.Directory {
		if d.Name == location && d.Location != "" {
			return &d, path, nil
		}
	}
	return nil, "", errorrepo.NewError("REST00112", location)
}

// BrowseLocation implements browseLocation operation.
//
// Retrieves a list of files in the defined location.
//
// GET /rest/file/browse/{path}
func (Handler) BrowseLocation(ctx context.Context, params api.BrowseLocationParams) (r api.BrowseLocationRes, _ error) {
	d, path, err := extraceLocationPath(params.Path)
	if err != nil {
		return &api.BrowseLocationNotFound{Error: api.NewOptErrorError(api.ErrorError{Message: api.NewOptString(err.Error())})}, nil
	}
	session := ctx.(*clu.Context)
	if !auth.ValidUser(auth.UserRole, false, session.User(), "<"+d.Name) {
		return &api.BrowseLocationForbidden{}, nil
	}
	fileName := os.ExpandEnv(d.Location + "/" + path)
	log.Log.Debugf("FileName %s Filter %s", fileName, params.Filter.Value)
	f, ferr := os.Open(fileName)
	if ferr != nil {
		err := errorrepo.NewError("REST00113", d.Name, ferr)
		log.Log.Errorf("Error browsing file %v", err)
		return &api.BrowseLocationNotFound{Error: api.NewOptErrorError(api.ErrorError{Message: api.NewOptString(err.Error())})}, nil
	}
	fileInfo, fierr := f.Stat()
	if fierr != nil {
		err := errorrepo.NewError("REST00114", d.Name, fierr)
		return &api.BrowseLocationNotFound{Error: api.NewOptErrorError(api.ErrorError{Message: api.NewOptString(err.Error())})}, nil
	}
	if fileInfo.IsDir() {
		return returnDirectoryInfo(d, path, params.Filter.Value, f)
	}
	return returnFileInfo(f)
}

func returnFileInfo(f *os.File) (r api.BrowseLocationRes, _ error) {
	fl := api.File{}
	fi, err := f.Stat()
	if err != nil {
		return nil, err
	}
	fl.Name = api.NewOptString(f.Name())
	fl.Type = api.NewOptString("File")
	fl.Modified = api.NewOptDateTime(fi.ModTime())
	fl.Size = api.NewOptInt64(fi.Size())
	ok := &api.BrowseLocationOK{Type: api.FileBrowseLocationOK,
		File: fl}
	return ok, nil

}

// returnDirectoryInfo generate directory information list
func returnDirectoryInfo(d *clu.Directory, path, pattern string, f *os.File) (api.BrowseLocationRes, error) {
	files, err := f.ReadDir(0)
	if err != nil {
		err := errorrepo.NewError("REST00115", d.Name, err)
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

// UploadFile implements uploadFile operation.
//
// Upload a new file to the given location.
//
// POST /rest/file/{location}
func (Handler) UploadFile(ctx context.Context, req *api.UploadFileReq, params api.UploadFileParams) (r api.UploadFileRes, _ error) {
	d, path, err := extraceLocationPath(params.Path)
	if err != nil {
		log.Log.Errorf("Error extracting location path: %v (original=%v)", err, params.Path)
		return &api.UploadFileNotFound{Error: api.NewOptErrorError(api.ErrorError{Message: api.NewOptString(err.Error())})}, nil
	}
	session := ctx.(*clu.Context)
	if !auth.ValidUser(auth.UserRole, true, session.User(), ">"+d.Name) {
		return &api.UploadFileForbidden{}, nil
	}

	fileName := os.ExpandEnv(d.Location + "/" + path)
	if params.File.IsSet() {
		fileName += params.File.Value
	}
	fileName = filepath.Clean(fileName)
	log.Log.Debugf("Final file name: " + fileName)
	//	f, err := os.Create(fileName)
	f, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE|os.O_EXCL, 0644)
	if err != nil {
		err = errorrepo.NewError("REST00116", fileName, err)
		return &api.UploadFileBadRequest{Error: api.NewOptErrorError(api.ErrorError{Message: api.NewOptString(err.Error())})}, nil
	}
	n, err := io.Copy(f, req.UploadFile.File)
	if err != nil {
		err = errorrepo.NewError("REST00117", fileName, err)
		return &api.UploadFileBadRequest{Error: api.NewOptErrorError(api.ErrorError{Message: api.NewOptString(err.Error())})}, nil
	}
	log.Log.Debugf("Read/Write bytes %d", n)
	v := api.StatusResponseStatus{Message: api.NewOptString("Wrote")}
	status := &api.StatusResponse{Status: api.NewOptStatusResponseStatus(v)}
	return status, nil
}

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
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	ht "github.com/ogen-go/ogen/http"
	"github.com/tknie/clu/api"
)

// BrowseList implements browseList operation.
//
// Retrieves a list of Browseable locations.
//
// GET /rest/file/browse
func (ServerHandler) BrowseList(ctx context.Context) (r api.BrowseListRes, _ error) {
	return r, ht.ErrNotImplemented
}

// CreateDirectory implements createDirectory operation.
//
// Create a new directory.
//
// PUT /rest/file/{location}
func (ServerHandler) CreateDirectory(ctx context.Context, params api.CreateDirectoryParams) (r api.CreateDirectoryRes, _ error) {
	return r, ht.ErrNotImplemented
}

// DeleteFileLocation implements deleteFileLocation operation.
//
// Delete the file on the given location.
//
// DELETE /rest/file/{location}
func (ServerHandler) DeleteFileLocation(ctx context.Context, params api.DeleteFileLocationParams) (r api.DeleteFileLocationRes, _ error) {
	return r, ht.ErrNotImplemented
}

// DownloadFile implements downloadFile operation.
//
// Download a file out of file location.
//
// GET /rest/file/{location}
func (ServerHandler) DownloadFile(ctx context.Context, params api.DownloadFileParams) (r api.DownloadFileRes, _ error) {
	return r, ht.ErrNotImplemented
}

// BrowseLocation implements browseLocation operation.
//
// Retrieves a list of files in the defined location.
//
// GET /rest/file/browse/{path}
func (ServerHandler) BrowseLocation(ctx context.Context, params api.BrowseLocationParams) (r api.BrowseLocationRes, _ error) {
	location := filepath.Clean(params.Path)
	fmt.Println("LLLL", location)
	firstSlash := strings.IndexByte(location, '/')
	path := "/"
	if firstSlash == 0 {
		err := fmt.Errorf("location reference missing")
		return &api.BrowseLocationNotFound{Error: api.NewOptErrorError(api.ErrorError{Message: api.NewOptString(err.Error())})}, nil
	}
	if firstSlash > 0 {
		location = params.Path[:firstSlash]
		path = params.Path[firstSlash:]
	}
	fmt.Println("location=", location, "path=", path)

	for _, d := range Viewer.FileTransfer.Directories.Directory {
		fmt.Println(d.Location, d.Name, location)
		if d.Name == location {
			fileName := os.ExpandEnv(d.Location + "/" + path)
			fmt.Println("OOO", fileName)
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
				files, err := ioutil.ReadDir(fileName)
				if err != nil {
					err := fmt.Errorf("error reading path of location %s", d.Name)
					return &api.BrowseLocationNotFound{Error: api.NewOptErrorError(api.ErrorError{Message: api.NewOptString(err.Error())})}, nil
				}
				fl := &api.DirectoryFiles{Location: api.NewOptString(location),
					Path:   api.NewOptString(path),
					System: api.NewOptString(runtime.GOOS)}

				for _, f := range files {
					fileType := "File"
					if f.IsDir() {
						fileType = "Directory"
					}
					content := api.File{Name: api.NewOptString(f.Name()), Type: api.NewOptString(fileType),
						Modified: api.NewOptDateTime(f.ModTime()), Size: api.NewOptInt64(f.Size())}
					fl.Files = append(fl.Files, content)
				}
				return fl, nil
			}
			content := api.File{Name: api.NewOptString(fileInfo.Name()), Type: api.NewOptString("File"),
				Modified: api.NewOptDateTime(fileInfo.ModTime()), Size: api.NewOptInt64(fileInfo.Size())}
			return &api.DirectoryFiles{Files: []api.File{content}, System: api.NewOptString(runtime.GOOS)}, nil
		}
	}

	return r, ht.ErrNotImplemented
}

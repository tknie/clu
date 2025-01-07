//go:build cgo && !windows
// +build cgo,!windows

/*
* Copyright 2022-2025 Thorsten A. Knieling
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
	"image"
	"io"

	"github.com/tknie/goheif"
)

func heifdecoder(r io.Reader) (image.Image, error) {
	return goheif.Decode(r)
}

func heifextractor(ra io.ReaderAt) ([]byte, error) {
	return goheif.ExtractExif(ra)
}

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
	"bytes"
	"fmt"
	"image/jpeg"
	"io"
	"os"
	"strings"
	"time"

	"github.com/disintegration/imaging"
	"github.com/rwcarlsen/goexif/exif"
	"github.com/rwcarlsen/goexif/tiff"
	"github.com/tknie/clu"
	"github.com/tknie/errorrepo"
	"github.com/tknie/flynn/common"
	"github.com/tknie/log"
	"github.com/tknie/services"
)

const blockSize = 65536

type StreamRead struct {
	table         string
	field         string
	mimetypeField string
	mimetype      string
	data          []byte
	send          int
}

// NewStreamRead new stream read instance
func NewStreamRead(table, field, mimetypeField string) *StreamRead {
	return &StreamRead{table: table,
		field:         field,
		mimetypeField: mimetypeField}
}

// initStreamFromTable init streamed read and open connection to database
func (read *StreamRead) initStreamFromTable(srvctx *clu.Context, search, destMimeType string) (err error) {
	d, err := ConnectTable(srvctx, read.table)
	if err != nil {
		log.Log.Errorf("Error search table %s:%v", read.table, err)
		return err
	}
	defer CloseTable(d)

	log.Log.Debugf("Init stream for table %s and search %s for field %s", read.table, search, read.field)

	fields := []string{strings.ToLower(read.field)}
	if read.mimetypeField != "" {
		fields = []string{strings.ToLower(read.field), read.mimetypeField}
	}

	q := &common.Query{TableName: read.table,
		Fields: fields,
		Search: search}
	result, err := queryBytes(d, q)
	if err != nil {
		log.Log.Errorf("Error query table %s:%v", read.table, err)
		return err
	}

	if len(result) == 0 {
		err = errorrepo.NewError("REST00002", read.field, read.table)
		// err = fmt.Errorf("field '%s' not in result", field)
		return
	}

	s := strings.ToLower(read.field)
	if d, ok := result[s]; ok {
		if d == nil {
			return fmt.Errorf("stream query result empty")
		}
		read.data = d.([]byte)
		if read.mimetypeField != "" {
			s := strings.ToLower(read.mimetypeField)
			if d, ok := result[s]; ok {
				read.mimetype = d.(string)
			}
		}
		return read.convertMimeType(destMimeType)
	}
	err = errorrepo.NewError("RERR00002", read.field, read.table)
	// err = fmt.Errorf("field '%s' not in result", field)
	return
}

// initStreamFromFile init streamed read and open file
func initStreamFromFile(file *os.File) (read *StreamRead, err error) {
	data, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}
	read = &StreamRead{data: data}
	read.mimetype = "application/octet-stream"

	return
}

// nextDataBlock read partial large object segment
func (read *StreamRead) nextDataBlock() (data []byte, err error) {
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
func (read *StreamRead) streamResponderFunc() (io.Reader, error) {
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

// query query SQL tables
func queryBytes(d common.RegDbID, query *common.Query) (map[string]interface{}, error) {
	log.Log.Debugf("Query stream in db ID %04d", d)
	dataMap := make(map[string]interface{})
	found := false
	_, err := d.Query(query, func(search *common.Query, result *common.Result) error {
		if result == nil {
			return fmt.Errorf("result empty")
		}
		if found {
			return fmt.Errorf("result not unique")
		}
		log.Log.Debugf("Rows: %d", len(result.Rows))
		///var d api.ResponseRecordsItem
		for i, r := range result.Rows {
			s := strings.ToLower(result.Fields[i])
			log.Log.Debugf("%d. row is of type %T", i, r)
			switch t := r.(type) {
			case *string:
				log.Log.Debugf("String %s", *t)
				dataMap[s] = *t
			case *time.Time:
				dataMap[s] = *t
			default:
				dataMap[s] = r
			}
		}
		found = true
		return nil
	})
	if err != nil {
		return nil, err
	}
	return dataMap, nil
}

// convertMimeType convert mime type
func (read *StreamRead) convertMimeType(destMimeType string) error {
	log.Log.Debugf("Convert %s ->  %s len=%d", read.mimetype, destMimeType, len(read.data))
	if destMimeType != "" && read.mimetype != "" && destMimeType != read.mimetype {
		log.Log.Debugf("Check destination mimetype")
		switch strings.ToLower(read.mimetype) {
		case "image/heic":
			r := bytes.NewBuffer(read.data)
			srcImage, err := heifdecoder(r)
			if err != nil {
				log.Log.Debugf("Decode image for conversion error %v", err)
				return err
			}
			ra := bytes.NewReader(read.data)
			exifData, err := heifextractor(ra)
			if err != nil {
				log.Log.Debugf("Extract exif error %v", err)
				return err
			}
			x, err := exif.Decode(bytes.NewBuffer(exifData))
			if err == nil {
				log.Log.Debugf("Decode exif in image")
				t, err := x.Get(exif.Orientation)
				if err == nil {
					switch t.String() {
					case "1":
					case "2":
						srcImage = imaging.FlipV(srcImage)
					case "3":
						srcImage = imaging.Rotate180(srcImage)
					case "4":
						srcImage = imaging.Rotate180(imaging.FlipV(srcImage))
					case "5":
						srcImage = imaging.Rotate270(imaging.FlipV(srcImage))
					case "6":
						srcImage = imaging.Rotate270(srcImage)
					case "7":
						srcImage = imaging.Rotate90(imaging.FlipV(srcImage))
					case "8":
						srcImage = imaging.Rotate90(srcImage)
					}
				}
			} else {
				log.Log.Debugf("Decode exif in image error %v", err)
			}
			buf := new(bytes.Buffer)
			err = jpeg.Encode(buf, srcImage, nil)
			if err != nil {
				log.Log.Debugf("Encode image for jpeg error %v", err)
				return err
			}
			read.data = buf.Bytes()
		case "image/jpeg", "image/jpg", "image/gif":
		default:
			log.Log.Debugf("No convert available -> %s", read.mimetype)
			services.ServerMessage("No convert available -> %s", read.mimetype)
		}
	}
	return nil
}

// Printer temporary
type Printer struct {
}

// Walk walk in printer
func (p *Printer) Walk(name exif.FieldName, tag *tiff.Tag) error {
	fmt.Printf("%s: %s\n", name, tag)
	return nil
}

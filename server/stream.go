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
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/tknie/clu"
	"github.com/tknie/flynn/common"
	"github.com/tknie/log"
)

const blockSize = 65536

type streamRead struct {
	field    string
	mimetype string
	data     []byte
	send     int
}

// initStreamFromTable init streamed read and open connection to database
func initStreamFromTable(srvctx *clu.Context, table,
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
			return nil, fmt.Errorf("stream query result empty")
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
	err = fmt.Errorf("field not in result")
	return
}

// initStreamFromTable init streamed read and open connection to database
func initStreamFromFile(file *os.File) (read *streamRead, err error) {
	data, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}
	read = &streamRead{data: data}
	read.mimetype = "application/octet-stream"

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

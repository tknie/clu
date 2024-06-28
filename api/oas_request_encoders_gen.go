// Code generated by ogen, DO NOT EDIT.

package api

import (
	"bytes"
	"mime"
	"mime/multipart"
	"net/http"

	"github.com/go-faster/errors"
	"github.com/go-faster/jx"

	ht "github.com/ogen-go/ogen/http"
	"github.com/ogen-go/ogen/uri"
)

func encodeBatchQueryRequest(
	req BatchQueryReq,
	r *http.Request,
) error {
	switch req := req.(type) {
	case *BatchQueryReqEmptyBody:
		// Empty body case.
		return nil
	case *SQLQuery:
		const contentType = "application/json"
		e := new(jx.Encoder)
		{
			req.Encode(e)
		}
		encoded := e.Bytes()
		ht.SetBody(r, bytes.NewReader(encoded), contentType)
		return nil
	case *BatchQueryReqTextPlain:
		const contentType = "text/plain"
		body := req
		ht.SetBody(r, body, contentType)
		return nil
	default:
		return errors.Errorf("unexpected request type: %T", req)
	}
}

func encodeCallPostExtendRequest(
	req *CallPostExtendReq,
	r *http.Request,
) error {
	const contentType = "multipart/form-data"
	request := req

	q := uri.NewFormEncoder(map[string]string{})
	body, boundary := ht.CreateMultipartBody(func(w *multipart.Writer) error {
		if val, ok := request.UploadFile.Get(); ok {
			if err := val.WriteMultipart("uploadFile", w); err != nil {
				return errors.Wrap(err, "write \"uploadFile\"")
			}
		}
		if err := q.WriteMultipart(w); err != nil {
			return errors.Wrap(err, "write multipart")
		}
		return nil
	})
	ht.SetCloserBody(r, body, mime.FormatMediaType(contentType, map[string]string{"boundary": boundary}))
	return nil
}

func encodeInsertMapFileRecordsRequest(
	req OptInsertMapFileRecordsReq,
	r *http.Request,
) error {
	const contentType = "multipart/form-data"
	if !req.Set {
		// Keep request with empty body if value is not set.
		return nil
	}
	request := req.Value

	q := uri.NewFormEncoder(map[string]string{})
	body, boundary := ht.CreateMultipartBody(func(w *multipart.Writer) error {
		if val, ok := request.Data.Get(); ok {
			if err := val.WriteMultipart("data", w); err != nil {
				return errors.Wrap(err, "write \"data\"")
			}
		}
		if err := q.WriteMultipart(w); err != nil {
			return errors.Wrap(err, "write multipart")
		}
		return nil
	})
	ht.SetCloserBody(r, body, mime.FormatMediaType(contentType, map[string]string{"boundary": boundary}))
	return nil
}

func encodeInsertRecordRequest(
	req OptInsertRecordReq,
	r *http.Request,
) error {
	const contentType = "application/json"
	if !req.Set {
		// Keep request with empty body if value is not set.
		return nil
	}
	e := new(jx.Encoder)
	{
		if req.Set {
			req.Encode(e)
		}
	}
	encoded := e.Bytes()
	ht.SetBody(r, bytes.NewReader(encoded), contentType)
	return nil
}

func encodePostDatabaseRequest(
	req *Database,
	r *http.Request,
) error {
	const contentType = "application/json"
	e := new(jx.Encoder)
	{
		if req != nil {
			req.Encode(e)
		}
	}
	encoded := e.Bytes()
	ht.SetBody(r, bytes.NewReader(encoded), contentType)
	return nil
}

func encodePostJobRequest(
	req PostJobReq,
	r *http.Request,
) error {
	switch req := req.(type) {
	case *PostJobReqEmptyBody:
		// Empty body case.
		return nil
	case *JobParameter:
		const contentType = "application/json"
		e := new(jx.Encoder)
		{
			req.Encode(e)
		}
		encoded := e.Bytes()
		ht.SetBody(r, bytes.NewReader(encoded), contentType)
		return nil
	case *PostJobReqTextPlain:
		const contentType = "text/plain"
		body := req
		ht.SetBody(r, body, contentType)
		return nil
	default:
		return errors.Errorf("unexpected request type: %T", req)
	}
}

func encodeSetConfigRequest(
	req SetConfigReq,
	r *http.Request,
) error {
	switch req := req.(type) {
	case *Config:
		const contentType = "application/json"
		e := new(jx.Encoder)
		{
			req.Encode(e)
		}
		encoded := e.Bytes()
		ht.SetBody(r, bytes.NewReader(encoded), contentType)
		return nil
	case *SetConfigReqTextPlain:
		const contentType = "text/plain"
		body := req
		ht.SetBody(r, body, contentType)
		return nil
	default:
		return errors.Errorf("unexpected request type: %T", req)
	}
}

func encodeSetJobsConfigRequest(
	req OptJobStore,
	r *http.Request,
) error {
	const contentType = "application/json"
	if !req.Set {
		// Keep request with empty body if value is not set.
		return nil
	}
	e := new(jx.Encoder)
	{
		if req.Set {
			req.Encode(e)
		}
	}
	encoded := e.Bytes()
	ht.SetBody(r, bytes.NewReader(encoded), contentType)
	return nil
}

func encodeUpdateLobByMapRequest(
	req UpdateLobByMapReq,
	r *http.Request,
) error {
	switch req := req.(type) {
	case *UpdateLobByMapReqApplicationOctetStream:
		const contentType = "application/octet-stream"
		body := req
		ht.SetBody(r, body, contentType)
		return nil
	case *UpdateLobByMapReqMultipartFormData:
		const contentType = "multipart/form-data"
		request := req

		q := uri.NewFormEncoder(map[string]string{})
		body, boundary := ht.CreateMultipartBody(func(w *multipart.Writer) error {
			if err := request.UploadLob.WriteMultipart("uploadLob", w); err != nil {
				return errors.Wrap(err, "write \"uploadLob\"")
			}
			if err := q.WriteMultipart(w); err != nil {
				return errors.Wrap(err, "write multipart")
			}
			return nil
		})
		ht.SetCloserBody(r, body, mime.FormatMediaType(contentType, map[string]string{"boundary": boundary}))
		return nil
	default:
		return errors.Errorf("unexpected request type: %T", req)
	}
}

func encodeUpdateRecordsByFieldsRequest(
	req OptUpdateRecordsByFieldsReq,
	r *http.Request,
) error {
	const contentType = "application/json"
	if !req.Set {
		// Keep request with empty body if value is not set.
		return nil
	}
	e := new(jx.Encoder)
	{
		if req.Set {
			req.Encode(e)
		}
	}
	encoded := e.Bytes()
	ht.SetBody(r, bytes.NewReader(encoded), contentType)
	return nil
}

func encodeUploadFileRequest(
	req *UploadFileReq,
	r *http.Request,
) error {
	const contentType = "multipart/form-data"
	request := req

	q := uri.NewFormEncoder(map[string]string{})
	body, boundary := ht.CreateMultipartBody(func(w *multipart.Writer) error {
		if err := request.UploadFile.WriteMultipart("uploadFile", w); err != nil {
			return errors.Wrap(err, "write \"uploadFile\"")
		}
		if err := q.WriteMultipart(w); err != nil {
			return errors.Wrap(err, "write multipart")
		}
		return nil
	})
	ht.SetCloserBody(r, body, mime.FormatMediaType(contentType, map[string]string{"boundary": boundary}))
	return nil
}

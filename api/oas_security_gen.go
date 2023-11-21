package api

import (
	"context"
	"net/http"
	"strings"

	"github.com/go-faster/errors"

	"github.com/ogen-go/ogen/ogenerrors"
)

// SecurityHandler is handler for security parameters.
type SecurityHandler interface {
	// HandleBasicAuth handles BasicAuth security.
	// HTTP Basic Authentication. Works over `HTTP` and `HTTPS`.
	HandleBasicAuth(ctx context.Context, operationName string, t BasicAuth) (context.Context, error)
	// HandleBearerAuth handles BearerAuth security.
	// HTTP Bearer Authentication. Works over `HTTP` and `HTTPS`.
	HandleBearerAuth(ctx context.Context, operationName string, t BearerAuth) (context.Context, error)
	// HandleTokenCheck handles tokenCheck security.
	// HTTP Basic Authentication. Works over `HTTP` and `HTTPS`.
	HandleTokenCheck(ctx context.Context, operationName string, t TokenCheck) (context.Context, error)
	// Request call after receiving
	Request(ctx context.Context, req *http.Request)
}

func findAuthorization(h http.Header, prefix string) (string, bool) {
	v, ok := h["Authorization"]
	if !ok {
		return "", false
	}
	for _, vv := range v {
		scheme, value, ok := strings.Cut(vv, " ")
		if !ok || !strings.EqualFold(scheme, prefix) {
			continue
		}
		return value, true
	}
	return "", false
}

func (s *Server) securityBasicAuth(ctx context.Context, operationName string, req *http.Request) (context.Context, bool, error) {
	var t BasicAuth
	if _, ok := findAuthorization(req.Header, "Basic"); !ok {
		return ctx, false, nil
	}
	username, password, ok := req.BasicAuth()
	if !ok {
		return nil, false, errors.New("invalid basic auth")
	}
	t.Username = username
	t.Password = password
	rctx, err := s.sec.HandleBasicAuth(ctx, operationName, t)
	if errors.Is(err, ogenerrors.ErrSkipServerSecurity) {
		return nil, false, nil
	} else if err != nil {
		return nil, false, err
	}
	// CLU TKN addition
	s.sec.Request(rctx, req)

	return rctx, true, err
}

func (s *Server) securityBearerAuth(ctx context.Context, operationName string, req *http.Request) (context.Context, bool, error) {
	var t BearerAuth
	token, ok := findAuthorization(req.Header, "Bearer")
	if !ok {
		return ctx, false, nil
	}
	t.Token = token
	rctx, err := s.sec.HandleBearerAuth(ctx, operationName, t)
	if errors.Is(err, ogenerrors.ErrSkipServerSecurity) {
		return nil, false, nil
	} else if err != nil {
		return nil, false, err
	}
	// CLU TKN addition
	s.sec.Request(rctx, req)

	return rctx, true, err
}
func (s *Server) securityTokenCheck(ctx context.Context, operationName string, req *http.Request) (context.Context, bool, error) {
	var t TokenCheck
	const parameterName = "X-Tokencheck"
	value := req.Header.Get(parameterName)
	if value == "" {
		return ctx, false, nil
	}
	t.APIKey = value
	rctx, err := s.sec.HandleTokenCheck(ctx, operationName, t)
	if errors.Is(err, ogenerrors.ErrSkipServerSecurity) {
		return nil, false, nil
	} else if err != nil {
		return nil, false, err
	}
	// CLU TKN addition
	s.sec.Request(rctx, req)

	return rctx, true, err
}

// SecuritySource is provider of security values (tokens, passwords, etc.).
type SecuritySource interface {
	// BasicAuth provides BasicAuth security value.
	// HTTP Basic Authentication. Works over `HTTP` and `HTTPS`.
	BasicAuth(ctx context.Context, operationName string) (BasicAuth, error)
	// BearerAuth provides BearerAuth security value.
	// HTTP Bearer Authentication. Works over `HTTP` and `HTTPS`.
	BearerAuth(ctx context.Context, operationName string) (BearerAuth, error)
	// TokenCheck provides tokenCheck security value.
	// HTTP Basic Authentication. Works over `HTTP` and `HTTPS`.
	TokenCheck(ctx context.Context, operationName string) (TokenCheck, error)
}

func (s *Client) securityBasicAuth(ctx context.Context, operationName string, req *http.Request) error {
	t, err := s.sec.BasicAuth(ctx, operationName)
	if err != nil {
		return errors.Wrap(err, "security source \"BasicAuth\"")
	}
	req.SetBasicAuth(t.Username, t.Password)
	return nil
}
func (s *Client) securityBearerAuth(ctx context.Context, operationName string, req *http.Request) error {
	t, err := s.sec.BearerAuth(ctx, operationName)
	if err != nil {
		return errors.Wrap(err, "security source \"BearerAuth\"")
	}
	req.Header.Set("Authorization", "Bearer "+t.Token)
	return nil
}
func (s *Client) securityTokenCheck(ctx context.Context, operationName string, req *http.Request) error {
	t, err := s.sec.TokenCheck(ctx, operationName)
	if err != nil {
		return errors.Wrap(err, "security source \"TokenCheck\"")
	}
	req.Header.Set("X-Tokencheck", t.APIKey)
	return nil
}

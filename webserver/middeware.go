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

package webserver

import (
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/rs/cors"
	"github.com/tknie/clu/api"
	"github.com/tknie/clu/plugins"
	"github.com/tknie/clu/server"
	"github.com/tknie/log"
	"github.com/tknie/services"
)

var prefixesOfServices = []string{"/rest/", "/binary/",
	"/image/", "/video",
	"/adabas/", "/version",
	"/scheduler/", "/file/",
	"/login", "/logout",
	"/shutdown", "/redirect",
	"/docs", "/swagger.json"}
var staticContent = os.Getenv("API_STATIC_CONTENT")

type logWrapper struct {
}

func (lw *logWrapper) Printf(fm string, args ...interface{}) {
	log.Log.Debugf("CORS: "+fm, args...)
}

type myMiddleware struct {
	Next *api.Server
}

func (m *myMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	route, ok := m.Next.FindRoute(r.Method, r.URL.Path)
	if !ok {
		// There is no route for the request.
		//
		// Let server handle 404/405.
		m.Next.ServeHTTP(w, r)
		return
	}
	// Match operation by spec operation ID.
	// Notice that the operation ID is optional and may be empty.
	//
	// You can also use `route.Name()` to get the ogen operation name.
	// Unlike the operation ID, the name is guaranteed to be unique and non-empty.
	switch route.OperationID() {
	case "operation1", "operation2":
		// Middleware logic:
		args := route.Args()
		if args[0] == "debug" {
			w.Header().Set("X-Debug", "true")
		}
	}
	m.Next.ServeHTTP(w, r)
}

// InitMiddleWare init middleware steps
func InitMiddleWare(s *api.Server) http.Handler {
	mainHandler = &myMiddleware{Next: s}
	corsHandler := cors.New(cors.Options{
		Debug:            (services.LogLevel == 1),
		AllowedHeaders:   []string{"*"},
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{http.MethodGet, http.MethodPut, http.MethodPost, http.MethodDelete, http.MethodOptions},
		AllowCredentials: true,
		MaxAge:           1000,
	})
	corsHandler.Log = new(logWrapper)
	mainHandler = corsHandler.Handler(s)
	mainHandler = fileServerMiddleware(mainHandler)
	return mainHandler
}

// fileServerMiddleware file server static file handler
func fileServerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if log.IsDebugLevel() {
			log.Log.Debugf("Serve %s", r.URL.Path)
		}
		path := clearPath(r)

		if log.IsDebugLevel() {
			log.Log.Debugf("Serving %s", path)
		}
		for _, s := range prefixesOfServices {
			if strings.HasPrefix(path, s) {
				r.URL.Path = path
				if plugins.HasPlugins() {
					plugins.ReceiveAudit(nil, r)
					defer plugins.SendAudit(time.Now(), w, r)
				}
				next.ServeHTTP(w, r)
				return
			}
		}
		if log.IsDebugLevel() {
			log.Log.Debugf("Serve file %s on %s -> %s", r.RequestURI,
				server.Viewer.Server.Content, r.URL.Path)
		}
		directory := staticContent
		if directory == "" {
			directory = os.ExpandEnv(server.Viewer.Server.Content)
		}
		log.Log.Debugf("Directory search for: %s", directory)
		if directory == "" {
			w.Write([]byte("Error reading static content"))
			return
		}
		http.FileServer(http.Dir(directory)).ServeHTTP(w, r)
	})
}

func clearPath(r *http.Request) string {
	return r.URL.Path
	// if server.Viewer.Server.Prefix == "" {
	// 	return r.URL.Path
	// }
	// path := r.URL.Path
	// for _, prefix := range strings.Split(server.Viewer.Server.Prefix, ",") {
	// 	path = strings.TrimPrefix(r.URL.Path, strings.Trim(prefix, " "))
	// }
	// return path
}

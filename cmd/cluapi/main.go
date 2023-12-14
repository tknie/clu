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

package main

import (
	"flag"
	"io"
	"net/http"
	"os"
	"runtime"

	"github.com/ogen-go/ogen/middleware"
	"github.com/tknie/clu"
	"github.com/tknie/clu/api"
	"github.com/tknie/clu/plugins"
	"github.com/tknie/clu/server"
	"github.com/tknie/clu/webserver"
	"github.com/tknie/log"
	"github.com/tknie/services"
	"github.com/tknie/services/auth"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var pidFile = ""

func main() {
	var shutdown bool
	var port int

	flag.BoolVar(&shutdown, "S", false, "shutdown API server")
	flag.StringVar(&pidFile, "P", "", "define PID file")
	flag.IntVar(&port, "p", 8080, "define HTTP port")
	flag.IntVar(&port, "s", 8081, "define HTTPS port")
	flag.Parse()

	if pidFile == "" {
		pidFile = server.DefaultPIDFile()
	}

	server.LoadMessages()

	if shutdown {
		services.ServerMessage("Shutdown requested ...")
		services.ShutdownServer(pidFile, 15)
		os.Exit(0)
	}

	auth.PrincipalCreater = func(s *auth.SessionInfo, user, pass string) auth.PrincipalInterface {
		log.Log.Debugf("Create principal %s UUID=%s with password", user, s.UUID)
		u := clu.CheckUserExist(user, s)
		if u == nil {
			log.Log.Fatalf("User info not found for user %s", user)
		}
		m := clu.NewContextUserInfo(u, pass)
		m.User.LongName = user
		m.Auth.Roles = []string{"user", "admin"}
		m.Auth.Session = s
		log.Log.Debugf("Create OGEN principal: %#v", m)
		return m
	}

	services.ServerMessage("Starting CLUAPI server version %s", services.BuildVersion)
	services.ServerMessage("Build date %s", services.BuildDate)
	services.ServerMessage("Go build version %s", runtime.Version())
	webserver.InitServices()

	// Load XML configuration
	server.InitConfig(true)
	plugins.InitPlugins()

	server.InitDatabases()
	dm := server.Viewer.Database.UserInfo
	r, err := server.Handles(dm)
	if err == nil {
		clu.InitUserInfo(r, os.ExpandEnv(dm.Password), os.ExpandEnv(dm.Table))
	}
	dm = server.Viewer.Database.SessionInfo
	r, err = server.Handles(dm)
	if err == nil {
		clu.InitStoreInfo(r, os.ExpandEnv(dm.Password), os.ExpandEnv(dm.Table))
	}

	server.AdaptConfig(os.Getenv(server.DefaultConfigFileEnv))

	s, err := api.NewServer(server.Handler{}, webserver.SecurityHandler{},
		api.WithNotFound(func(w http.ResponseWriter, r *http.Request) {
			_, _ = io.WriteString(w, `{"error_message": "not found"}`)
		}), api.WithPathPrefix(server.Viewer.Server.Prefix))
	if err != nil {
		panic(err)
	}
	mainHandler := webserver.InitMiddleWare(s)
	webserver.StartServices(mainHandler)

	defer services.ServerMessage("Shutdown initiated ...")
}

// Logging logging middleware
func Logging(logger *zap.Logger) middleware.Middleware {
	return func(
		req middleware.Request,
		next func(req middleware.Request) (middleware.Response, error),
	) (middleware.Response, error) {
		logger := logger.With(
			zap.String("operation", req.OperationName),
			zap.String("operationId", req.OperationID),
		)
		logger.Info("Handling request")
		resp, err := next(req)
		if err != nil {
			logger.Error("Fail", zap.Error(err))
		} else {
			var fields []zapcore.Field
			// Some response types may have a status code.
			// ogen provides a getter for it.
			//
			// You can write your own interface to match any response type.
			if tresp, ok := resp.Type.(interface{ GetStatusCode() int }); ok {
				fields = []zapcore.Field{
					zap.Int("status_code", tresp.GetStatusCode()),
				}
			}
			logger.Info("Success", fields...)
		}
		return resp, err
	}
}

package main

import (
	"fmt"

	"github.com/tknie/clu/plugins"
	"github.com/tknie/services/auth"
)

type testCallback struct {
	checkTokenErr    error
	generateTokenErr error
}

type testPrincipal struct {
	testUUID string
}

func (tp *testPrincipal) UUID() string {
	if tp.testUUID != "" {
		return tp.testUUID
	}
	return "TestUUID"
}
func (tp *testPrincipal) Name() string {
	return "TestPrincipal"
}
func (tp *testPrincipal) Remote() string {
	return "RemoteHost"
}
func (tp *testPrincipal) AddRoles(r []string) {
	fmt.Println("Add role", r)
}
func (tp *testPrincipal) SetRemote(r string) {
	fmt.Println("Set remote", r)
}
func (tp *testPrincipal) Roles() []string {
	return []string{"xx", "ME"}
}
func (tp *testPrincipal) Session() interface{} {
	return nil
}
func (tp *testPrincipal) SetSession(interface{}) {
}

func (tc *testCallback) GetName() string { return "testCallback" }
func (tc *testCallback) Init() error     { initAuthCallback(); return nil }
func (tc *testCallback) Authenticate(principal auth.PrincipalInterface, userName, passwd string) error {
	return nil
}
func (tc *testCallback) Authorize(principal auth.PrincipalInterface, userName, passwd string) error {
	return nil
}
func (tc *testCallback) CheckToken(token string, scopes []string) (auth.PrincipalInterface, error) {
	if tc.checkTokenErr != nil {
		return nil, tc.checkTokenErr
	}
	return &testPrincipal{}, tc.checkTokenErr
}
func (tc *testCallback) GenerateToken(IAt string,
	principal auth.PrincipalInterface) (tokenString string, err error) {
	if tc.generateTokenErr != nil {
		return "", tc.generateTokenErr
	}
	return "TESTTOKEN", tc.generateTokenErr
}

func initAuthCallback() {
	auth.RegisterCallback(&testCallback{})
}

type greeting string

// Types type of plugin working with
func (g greeting) Types() []plugins.PluginTypes {
	return []plugins.PluginTypes{plugins.AuthPlugin}
}

// Name name of the plugin
func (g greeting) Name() string {
	return "Auth Access"
}

// Version version of the number
func (g greeting) Version() string {
	return "1.0"
}

// Stop stop plugin
func (g greeting) Stop() {
}

// exported

// Callback test callback entry point
var Callback testCallback

// Loader loader for initialize plugin
var Loader greeting

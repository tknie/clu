#
# Copyright 2022-2023 Thorsten A. Knieling
#
# SPDX-License-Identifier: Apache-2.0
#
#   Licensed under the Apache License, Version 2.0 (the "License");
#   you may not use this file except in compliance with the License.
#   You may obtain a copy of the License at
#
#       http://www.apache.org/licenses/LICENSE-2.0
#

GOARCH            ?= $(shell $(GO) env GOARCH)
GOOS              ?= $(shell $(GO) env GOOS)
GOEXE             ?= $(shell $(GO) env GOEXE)

PACKAGE            = github.com/tknie/restdb
TESTPKGSDIR        = server
INSTALL_DEST       = .
DATE               = $(shell date +%d-%m-%Y'_'%H:%M:%S)
VERSION           ?= v2.0.1.0
MAJORVERS          = v201
RESTVERSION       ?= $(VERSION).$(shell date +%d%m%Y)
#$(shell git describe --tags --always --dirty --match=v* 2> /dev/null || \
#			cat $(CURDIR)/.version 2> /dev/null || echo v0)
BIN                = $(CURDIR)/bin/$(GOOS)_$(GOARCH)
PLUGINSBIN         = $(BIN)/plugins
PROMOTE            = $(CURDIR)/promote/$(GOOS)_$(GOARCH)
BINTOOLS           = $(CURDIR)/bin/tools/$(GOOS)_$(GOARCH)
TARFILE            = cluapi-$(GOOS)_$(GOARCH).tar.gz
LOGPATH            = $(CURDIR)/logs
CURLOGPATH         = $(CURDIR)/logs
TESTOUTPUT         = $(CURDIR)/test
TESTFILES          = $(CURDIR)/files
MESSAGES           = $(CURDIR)/messages
REFERENCES         = $(TESTFILES)/references
RESTEXEC           = $(BIN)/cmd/cluapi
EXECS              = $(RESTEXEC)
REST_SERVER        = $(CURDIR)
SERVER_HOME        = $(CURDIR)
PLUGINS            = $(PLUGINSBIN)/audit
OBJECTS            = cmd/*/*.go server/*.go *.go webserver/*.go plugins/*.go
SWAGGER_SPEC       = $(CURDIR)/swagger/openapi-restserver.yaml	
ENABLE_DEBUG      ?= 0
ARTIFACTORY       ?= http://lion.fritz.box:8081
ARTIFACTORY_PASS  ?= admin:12345
CERTIFICATE        = $(CURDIR)/keys/certificate.pem

CGO_CFLAGS         =
CGO_LDFLAGS        = 

export CGO_CFLAGS
export CGO_LDFLAGS
export LOGPATH CURDIR PLUGINSBIN REST_SERVER SERVER_HOME BIN

include $(CURDIR)/make/common.mk

generatemodels: cleanAPI $(CURDIR)/api

.PHONY: clean
clean: cleanModules cleanVendor cleanCommon ; $(info $(M) cleaning…)	@ ## Cleanup everything
	@rm -rf restadmin.test jobs.test
	@rm -rf bin pkg logs test promote $(CURDIR)/server.pid
	@rm -rf test/tests.* test/coverage.*
	@rm -rf $(CURDIR)/cmd/cluapi-server/logs

cleanVendor: ; $(info $(M) cleaning vendor…)    @ ## Cleanup vendor
	@rm -rf $(CURDIR)/vendor

cleanAPI: ; $(info $(M) cleaning models…)    @ ## Cleanup models
	@rm -rf $(CURDIR)/api

$(BIN)/cluapi: prepare fmt lint lib $(EXECS)

startServer: $(CERTIFICATE) $(BIN)/cmd/cluapi ; $(info $(M) starting server…)
	@rm -f $(CURDIR)/logs/*; \
	if [ ! -d $(CURDIR)/tmp ]; then \
	  mkdir $(CURDIR)/tmp; fi; \
	TEMP=$(CURDIR) \
	$(RESTEXEC) server -c $(CURDIR)/configuration/config.yaml --host= 
#--scheme=https
#--tls-certificate keys/certificate.pem --tls-key keys/key.pem  --port=8130 --tls-port=8131 --host=

stopServer: $(BIN)/cmd/cluapi-client ; $(info $(M) stopping server…)
	@echo "Stop $(VERSION) build at $(DATE)"
	$(RESTEXEC) client -s -c $(CURDIR)/configuration/config.yaml

# Dependency management
.PHONY: generate
$(CURDIR)/api: $(SWAGGER_SPEC) ; $(info $(M) generating code...) @ ## Generate rest go code
	$Q go generate .

$(CERTIFICATE):
	cd keys; sh generate_tls.sh

promote: exec ; $(info $(M) package for promotion…) @ ## package for promotion
	if [ ! -d $(CURDIR)/promote ]; then mkdir $(CURDIR)/promote; fi; \
	if [ ! -d $(PROMOTE) ]; then mkdir -f $(PROMOTE); fi; \
	if [ ! -d $(PROMOTE)/keys ]; then cp -r $(CURDIR)/templateInst $(PROMOTE)/; fi; \
	if [ ! -d $(PROMOTE)/plugins ]; then mkdir -f $(PROMOTE)/plugins; fi; \
	if [ ! -d $(PROMOTE)/static ]; then \
	  mkdir $(PROMOTE)/static/; \
	  cp -r $(CURDIR)/static/* $(PROMOTE)/static/; fi; \
	if [ ! -d $(PROMOTE)/bin ]; then mkdir $(PROMOTE)/bin/; fi; \
	cp -r $(RESTEXEC) $(PROMOTE)/bin/; \
	if [ ! -d $(PROMOTE)/plugins ]; then mkdir $(PROMOTE)/plugins/; fi; \
	cp -r $(PLUGINSBIN)/* $(PROMOTE)/plugins/; \
	if [ ! -d $(PROMOTE)/swagger ]; then mkdir $(PROMOTE)/swagger/; fi; \
	cp $(CURDIR)/swagger/* $(PROMOTE)/swagger/; \
	mkdir $(PROMOTE)/logs/; \
	if [ ! -r $(PROMOTE)/keys/certificate.pem ]; then \
	  cd $(PROMOTE)/keys; sh generate_tls.sh; fi
ifeq ($(shell go env GOOS),windows)
	mv $(PROMOTE)/configuration/config-win.xml $(PROMOTE)/configuration/config.yaml; \
	rm $(PROMOTE)/bin/*.sh
else
#	rm $(PROMOTE)/$(INSTALL_DEST)/configuration/config-win.xml; rm $(PROMOTE)/$(INSTALL_DEST)/bin/*.bat
endif

promoteTest: test-build ; $(info $(M) package for tests…) @ ## package for promotion
	if [ ! -d $(CURDIR)/promote ]; then mkdir $(CURDIR)/promote; fi; \
	if [ ! -d $(PROMOTE) ]; then mkdir $(PROMOTE); fi; \
	if [ ! -d $(PROMOTE)/RestTests ]; then mkdir $(PROMOTE)/RestTests; fi; \
	cp -r $(BINTESTS) $(PROMOTE)/RestTests/bin

upload: prepareUpload uploadNexusInterim uploadNexusPackage; $(info $(M) uploading…) @ ## uploading packages

uploadNexusInterim: ; $(info $(M) uploading…) @ ## uploading packages
ifeq ($(shell go env GOOS),windows)
	cd $(PROMOTE)/; ls; \
	curl -v -u $(ARTIFACTORY_PASS) -X POST '${ARTIFACTORY}/service/rest/v1/components?repository=maven' -F maven2.groupId=com.github.tknie..$(GOOS)_$(GOARCH) -F maven2.artifactId=cluapi-go -F maven2.version=$(MAJORVERS) -F maven2.asset1=@cluapi.zip -F maven2.asset1.extension=zip
else
	cd $(PROMOTE)/; \
	curl -v -u $(ARTIFACTORY_PASS) -X POST '${ARTIFACTORY}/service/rest/v1/components?repository=maven' -F maven2.groupId=com.github.tknie..$(GOOS)_$(GOARCH) -F maven2.artifactId=cluapi-go -F maven2.version=$(MAJORVERS) -F maven2.asset1=@../${TARFILE} -F maven2.asset1.extension=tar.gz
endif

uploadNexusPackage: ; $(info $(M) uploading to Nexus…) @ ## uploading packages
ifeq ($(shell go env GOOS),windows)
	cd $(PROMOTE)/; curl -v -u $(ARTIFACTORY_PASS) -X POST '${ARTIFACTORY}/service/rest/v1/components?repository=maven' -F maven2.groupId=com.github.tknie.Installer.ProductName -F maven2.artifactId=cluapi -F maven2.version=$(GOOS)_$(GOARCH).main -F maven2.asset1=@cluapi.zip -F maven2.asset1.extension=zip
else
	cd $(PROMOTE)/; curl -v -u $(ARTIFACTORY_PASS) -X POST '${ARTIFACTORY}/service/rest/v1/components?repository=maven' -F maven2.groupId=com.github.tknie.Installer.ProductName -F maven2.artifactId=cluapi -F maven2.version=$(GOOS)_$(GOARCH).main -F maven2.asset1=@../${TARFILE} -F maven2.asset1.extension=tar.gz
endif

prepareUpload: ; $(info $(M) prepare uploading…) @ ## uploading packages
ifeq ($(shell go env GOOS),windows)
	cd $(PROMOTE)/; rm -f cluapi.zip; zip -r cluapi.zip $(INSTALL_DEST)
else
	cd $(PROMOTE)/; rm -f ../${TARFILE}; tar cvfz ../${TARFILE} $(INSTALL_DEST)
endif

uploadTest: promoteTest ; $(info $(M) uploading tests…) @ ## uploading packages
ifeq ($(shell go env GOOS),windows)
	cd $(PROMOTE)/; rm -f RestTests.zip; zip -r RestTests.zip RestTests; \
	curl -v -u $(ARTIFACTORY_PASS) -X POST '${ARTIFACTORY}/service/rest/v1/components?repository=maven' -F maven2.groupId=com.github.tknie.$(GOOS)_$(GOARCH) -F maven2.artifactId=RestTests -F maven2.version=$(VERSION) -F maven2.asset1=@RestTests.zip -F maven2.asset1.extension=zip
else
	cd $(PROMOTE)/; rm -f RestTests.tar.gz; tar cfz RestTests.tar.gz RestTests; \
	curl -v -u $(ARTIFACTORY_PASS) -X POST '${ARTIFACTORY}/service/rest/v1/components?repository=maven' -F maven2.groupId=com.github.tknie.$(GOOS)_$(GOARCH) -F maven2.artifactId=RestTests -F maven2.version=$(VERSION) -F maven2.asset1=@RestTests.tar.gz -F maven2.asset1.extension=tar.gz
endif

webapp-update:
	@echo "WebApp not delivered !!!"

.PHONY: docker
docker: ; $(info $(M) genering docker image…)
	cp $(PROMOTE)/../${TARFILE}  $(CURDIR)/docker/cluapi.tar.gz
	cd docker; docker build --platform linux/amd64 -t cluapi .

#!/bin/sh

#=============================================================================
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
#=============================================================================

DIR=`dirname $0`
if [ "${DIR}" = "" ]; then
   DIR=.
fi

cd ${INSTALL_DIR}/
if [ ! -d "logs" ]; then
  mkdir logs
fi

SERVICE_TOOL=bin/cluapi-api-service

if [ "${ACTION}" = "console" ]; then
  "${SERVICE_TOOL}" console
  exit $?
fi

if [ "${ACTION}" = "start" ]; then
  "${SERVICE_TOOL}" start
  exit $?
fi

if [ "${ACTION}" = "stop" ]; then
  "${SERVICE_TOOL}" stop
  exit $?
fi

if [ "${ACTION}" = "install" ]; then
  "${SERVICE_TOOL}" installStart
  exit $?
fi

if [ "${ACTION}" = "installStart" ]; then
  "${SERVICE_TOOL}" installStart
  exit $?
fi

if [ "${ACTION}" = "configure" ]; then
  "${SERVICE_TOOL}" ${ACTION} ${ACTION_PARAMS}
  exit $?
fi

if [ "${ACTION}" = "remove" ]; then
  "${SERVICE_TOOL}" remove
  exit $?
fi

if [ "${ACTION}" = "status" ]; then
  "${SERVICE_TOOL}" status
  exit $?
fi

#USAGE

echo ""
echo "This script manages the CLUTRON RESTful administration system service."
echo ""
echo "Usage:"
echo ""
echo "    system_service.sh <command>"
echo ""
echo "commands:"
echo ""
echo "    install  install and start the CLUTRON RESTful administration system service"
echo "    remove   stop and remove the CLUTRON RESTful administration system service"
echo "    start    start the CLUTRON RESTful administration system service"
echo "    stop     stop the CLUTRON RESTful administration system service"
echo "    status   displays information about the service status"
echo "    console  start the CLUTRON RESTful administration Server in console mode"
echo ""
echo "Return values:"
echo ""
echo "    0     command executed sucessfull"
echo "    not 0 error (check output)"
echo ""

exit 1



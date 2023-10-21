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
#set -xv

usage() {
echo " CLUAPI Rest Server"
echo " Usage:"
echo "    service[options] command [params]"
echo "         options"
echo
echo "         commands"
echo "             run              start server instance in the same console"
echo "             start            start server instance in a new console"
echo "             stop             stop server instance"
echo "             restart          restart server instance"
echo "             ping             ping if server instance is running"
echo "             help             print this help"
exit 1
}

[ $# -gt 0 ] || usage

ACTION=$1
shift 1
ACTION_PARAMS=$*

# Examine script location
SCRIPT=$INSTALLDIR
SCRIPTPATH=$(dirname "$SCRIPT")

#####################################################
# Find a location for the SERVER console
#####################################################
SERVER_HOME=$INSTALLDIR
SERVER_CONFIG=${SERVER_CONFIG:-${DATADIR}/configuration/config.yaml}
export SERVER_CONFIG

LOGPATH=/data/logs/
SERVER_CONSOLE="${LOGPATH}/server.${HOSTNAME}.log"

# Possibility to override ports or TLS topics
#
# Define certificates
#TLS_CERTIFICATE=${CLUTRON_ADMIN_HOME}/tls/certificate.pem
#TLS_PRIVATE_KEY=${CLUTRON_ADMIN_HOME}/tls/key.pem
#export TLS_CERTIFICATE TLS_PRIVATE_KEY

HOST=
export HOST
# Define ports
#PORT=8130
#TLS_PORT=8131
#export PORT TLS_PORT

# starting Rest Interface (kernel)
start() {
  echo $(date +"%Y-%m-%d %H:%m:%S")" Starting API server in background mode"
  echo $(date +"%Y-%m-%d %H:%m:%S")" Server config file: ${SERVER_CONFIG}"
  echo $(date +"%Y-%m-%d %H:%m:%S")" Server startup console out: ${SERVER_CONSOLE}"

  cd ${SERVER_HOME}
  nohup bin/cluapi server -c ${SERVER_CONFIG} $* >${SERVER_CONSOLE} 2>&1 &
}

stop() {
  echo $(date +"%Y-%m-%d %H:%m%s")" Stopping API server"
  cd ${SERVER_HOME}
  bin/cluapi client -s
}

#####################################################
# Action!
#####################################################

case "$ACTION" in

  start)
        start ;;

  stop)
        stop ;;

  run)
        cd ${SERVER_HOME}

        echo $(date +"%Y-%m-%d %H:%m:%S")" Starting API server in foreground"
        echo $(date +"%Y-%m-%d %H:%m:%S")" Server config file: ${SERVER_CONFIG}"
        echo $(date +"%Y-%m-%d %H:%m:%S")" Server console out: ${SERVER_CONSOLE}"

        bin/cluapi server -c ${SERVER_CONFIG} $*
        ;;

  help) usage
        ;;

  *)
        echo "Invalid action: $ACTION"
        usage
        ;;
esac


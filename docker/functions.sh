#!/bin/bash
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

#
# Start offline scripts and delete
#
function offline_scripts {
   if [ -r /data/offline.sh ]; then
      sh /data/offline.sh
      rm /data/offline.sh
   fi
}

#
# Create common Docker specific parameters into INI
#
function init_parameter {
   echo $(date +"%Y-%m-%d %H:%m:%s")" Set default parameter set"
}


#
# Evaluate Docker parameters to be passed to parameter
#
function prepare_parameter {
   echo $(date +"%Y-%m-%d %H:%m:%s")" Prepare parameter set"
}

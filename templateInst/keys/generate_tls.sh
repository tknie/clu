#!/bin/sh
#=============================================================================
# Copyright 2022-2024 Thorsten A. Knieling
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

openssl req -new -x509 -nodes -days 365 -sha256 -newkey rsa:4096 -keyout key.pem -out certificate.pem -config csr_config.cnf
openssl x509 -text -noout -in certificate.pem
openssl pkcs12 -password pass: -inkey key.pem -in certificate.pem -export -out certificate.p12
openssl pkcs12 -password pass: -in certificate.p12 -noout -info


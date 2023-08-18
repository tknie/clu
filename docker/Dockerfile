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

FROM ubuntu:22.04

ENV user_id 1777
ENV group_id 1777
ENV INSTALLDIR /opt/cluapi

#RUN microdnf update -y && microdnf install -y shadow-utils util-linux && microdnf clean all && \ 
#    groupadd -g $group_id admin && \
#    useradd -d $INSTALLDIR -m -s /bin/bash -u $user_id -g admin admin
RUN apt update && apt install -y util-linux openssl && \
    rm -rf /var/lib/apt/lists/* && \
    apt clean && \
    groupadd -g $group_id admin && \
    useradd -d $INSTALLDIR -m -s /bin/bash -u $user_id -g admin admin

COPY startDocker stopDocker functions.sh $INSTALLDIR/docker/
ADD --chown=$user_id:$group_id cluapi.tar.gz $INSTALLDIR
COPY config.yaml $INSTALLDIR/configuration/config.yaml.template

# Remove old scripts and create environment
RUN cd $INSTALLDIR && \
    mkdir $INSTALLDIR/tmp && \
    mkdir -p /data && chown -R admin:admin /data $INSTALLDIR/tmp $INSTALLDIR/logs && \
    chmod 777 /data && \
    chmod 775 $INSTALLDIR && \
    rm $INSTALLDIR/keys/certificate.p12 $INSTALLDIR/keys/key.pem

USER admin

CMD ["/opt/cluapi/docker/startDocker"]

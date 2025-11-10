#
# Copyright 2022-2025 Thorsten A. Knieling
#
# SPDX-License-Identifier: Apache-2.0
#
#   Licensed under the Apache License, Version 2.0 (the "License");
#   you may not use this file except in compliance with the License.
#   You may obtain a copy of the License at
#
#       http://www.apache.org/licenses/LICENSE-2.0
#

FROM alpine:latest

ENV user_id 1777
ENV group_id 1777
ENV INSTALLDIR /opt/cluapi

RUN apk update && apk upgrade && apk add shadow openssl gcompat libde265 && \
     groupadd -g $group_id admin && \
     useradd -d $INSTALLDIR -m -s /bin/bash -u $user_id -g admin admin
COPY startDocker stopDocker functions.sh $INSTALLDIR/docker/
ADD --chown=$user_id:$group_id cluapi.tar.gz $INSTALLDIR
COPY service.sh $INSTALLDIR/bin/
COPY config.yaml $INSTALLDIR/configuration/config.yaml.template

# Remove old scripts and create environment
RUN cd $INSTALLDIR && \
    mkdir $INSTALLDIR/tmp && \
    mkdir -p /data && chown -R admin:admin /data $INSTALLDIR/tmp $INSTALLDIR/logs && \
    chmod 777 /data && \
    chmod 775 $INSTALLDIR && \
    rm $INSTALLDIR/keys/certificate.p12 $INSTALLDIR/keys/key.pem && \
    ln -sf /dev/stdout /data/server.log
#RUN ln -sf /dev/stderr /var/log/nginx/error.log

USER admin

CMD ["/opt/cluapi/docker/startDocker"]

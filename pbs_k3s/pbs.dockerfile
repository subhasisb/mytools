# Copyright (c) 2023 Altair Engineering Inc.
# All Rights Reserved.
# Copyright notice does not imply publication.
# Contains trade secret, proprietary, and confidential Information.

ARG REGISTRY=docker.io
#FROM centos:8 AS pbsmom
FROM iad.ocir.io/idt7ybnr03cb/objectstore_faas_templates:python-faas-templatev10.5 AS pbsmom
USER root
LABEL maintainer="Altair Engineering"

ENV DEBIAN_FRONTEND=noninteractive

RUN sed -i 's/mirrorlist/#mirrorlist/g' /etc/yum.repos.d/CentOS-* \
    && sed -i 's|#baseurl=http://mirror.centos.org|baseurl=http://vault.centos.org|g' /etc/yum.repos.d/CentOS-*

RUN yum update  -y \
    && yum install -y expat libedit postgresql-server postgresql-contrib python3 sendmail sudo tcl tk libical numactl-devel openssh-server munge

# copy PBS rpm
COPY ./pbspro-*.rpm ./pbspro.rpm

RUN yum install -y ./pbspro.rpm

ENV LD_LIBRARY_PATH=/usr/lib:/lib64:/usr/local/lib:/opt/pbs/lib
ENV PATH=/opt/pbs/bin:/opt/pbs/sbin:/opt/ptl/bin:/opt/lqs/bin:$PATH

RUN create-munge-key

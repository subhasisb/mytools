# Copyright (c) 2023 Altair Engineering Inc.
# All Rights Reserved.
# Copyright notice does not imply publication.
# Contains trade secret, proprietary, and confidential Information.

# copy the pbspro-server and pbspro-debug deb files to the current directory
# run using the following command from the current folder
# docker build . --build-context pbsprosrc=/home/subhas/work/forks/pbspro -t pbsimage:latest

ARG REGISTRY=docker.io
FROM ubuntu:latest AS pbsimage
USER root
LABEL maintainer="Altair Engineering"

ENV DEBIAN_FRONTEND=noninteractive

RUN apt-get update -y \
    && apt install -y expat libedit2 libedit-dev postgresql python3 sendmail sudo tcl tk libical3 libnuma-dev openssh-server libmunge-dev libmunge2 munge vim gdb iputils-ping strace

# copy PBS rpm
COPY ./pbspro*.deb /
RUN ls -l 

RUN apt install -y ./pbspro*.deb

ENV LD_LIBRARY_PATH=/usr/lib:/lib64:/usr/local/lib:/opt/pbs/lib
ENV PATH=/opt/pbs/bin:/opt/pbs/sbin:/opt/ptl/bin:/opt/lqs/bin:$PATH

RUN mkdir -p /home/subhas/work/forks/pbspro
COPY --from=pbsprosrc . /home/subhas/work/forks/pbspro/

RUN set -ex \
      && groupadd -g 1900 tstgrp00 \
      && groupadd -g 1901 tstgrp01 \
      && groupadd -g 1902 tstgrp02 \
      && groupadd -g 1903 tstgrp03 \
      && groupadd -g 1904 tstgrp04 \
      && groupadd -g 1905 tstgrp05 \
      && groupadd -g 1906 tstgrp06 \
      && groupadd -g 1907 tstgrp07 \
      && groupadd -g 901 pbs \
      && groupadd -g 1146 agt \
      && groupadd -g 2000 demogroup \
      && useradd  -m -s /bin/bash -u 4357 -g tstgrp00 -G tstgrp00 pbsadmin \
      && useradd  -m -s /bin/bash -u 9000 -g tstgrp00 -G tstgrp00 pbsbuild \
      && useradd  -m -s /bin/bash -u 884 -g tstgrp00 -G tstgrp00 pbsdata \
      && useradd  -m -s /bin/bash -u 4367 -g tstgrp00 -G tstgrp00 pbsmgr \
      && useradd  -m -s /bin/bash -u 4373 -g tstgrp00 -G tstgrp00 pbsnonroot \
      && useradd  -m -s /bin/bash -u 4356 -g tstgrp00 -G tstgrp00 pbsoper \
      && useradd  -m -s /bin/bash -u 4358 -g tstgrp00 -G tstgrp00 pbsother \
      && useradd  -m -s /bin/bash -u 4371 -g tstgrp00 -G tstgrp00 pbsroot \
      && useradd  -m -s /bin/bash -u 4355 -g tstgrp00 -G tstgrp02,tstgrp00 pbstest \
      && useradd  -m -s /bin/bash -u 4359 -g tstgrp00 -G tstgrp00 pbsuser \
      && useradd  -m -s /bin/bash -u 4361 -g tstgrp00 -G tstgrp01,tstgrp02,tstgrp00 pbsuser1 \
      && useradd  -m -s /bin/bash -u 4362 -g tstgrp00 -G tstgrp01,tstgrp03,tstgrp00 pbsuser2 \
      && useradd  -m -s /bin/bash -u 4363 -g tstgrp00 -G tstgrp01,tstgrp04,tstgrp00 pbsuser3 \
      && useradd  -m -s /bin/bash -u 4364 -g tstgrp01 -G tstgrp04,tstgrp05,tstgrp01 pbsuser4 \
      && useradd  -m -s /bin/bash -u 4365 -g tstgrp02 -G tstgrp04,tstgrp06,tstgrp02 pbsuser5 \
      && useradd  -m -s /bin/bash -u 4366 -g tstgrp03 -G tstgrp04,tstgrp07,tstgrp03 pbsuser6 \
      && useradd  -m -s /bin/bash -u 4368 -g tstgrp01 -G tstgrp01 pbsuser7 \
      && useradd  -m -s /bin/bash -u 11000 -g tstgrp00 -G tstgrp00 tstusr00 \
      && useradd  -m -s /bin/bash -u 11001 -g tstgrp00 -G tstgrp00 tstusr01 \
      && useradd  -m -s /bin/bash -u 2000 -g demogroup -G demogroup demouser \
      && chmod g+x,o+x /home/* \
      && echo 'root:pbs' | chpasswd  \
      && rm -rf ~/.ssh \
      && ssh-keygen -t rsa -C root-ssh-keypair -N "" -f ~/.ssh/id_rsa \
      && cp ~/.ssh/id_rsa.pub ~/.ssh/authorized_keys \
      && chmod 0600 ~/.ssh/authorized_keys \
      && for user in $(awk -F: '/^(demo|pbs|tst)/ {print $1}' /etc/passwd); do \
      rm -rf /home/${user}/.ssh; \
      echo "ssh-keygen -t rsa -C ${user}-ssh-keypair -N \"\" -f /home/${user}/.ssh/id_rsa" | sudo -Hiu ${user}; \
      sudo -Hiu ${user} cp /home/${user}/.ssh/id_rsa.pub /home/${user}/.ssh/authorized_keys; \
      sudo -Hiu ${user} chmod 0600 /home/${user}/.ssh/authorized_keys; \
      echo "${user}:pbs" | chpasswd; \
      done \
      && echo 'demouser:demo' | chpasswd \
      && echo 'Host *' >> /etc/ssh/ssh_config \
      && echo '  StrictHostKeyChecking no' >> /etc/ssh/ssh_config \
      && echo '  IdentityFile ~/.ssh/id_rsa' >> /etc/ssh/ssh_config \
      && echo '  PreferredAuthentications publickey,password' >> /etc/ssh/ssh_config \
      && rm -rf /etc/ssh/ssh_host_rsa_key \
      && ssh-keygen -t rsa -b 4096 -N "" -f /etc/ssh/ssh_host_rsa_key

RUN   chown -R root:root /var/lib/munge /etc/munge /var/log/munge \
      && ln -s /usr/lib/x86_64-linux-gnu/libmunge.so.2 /usr/lib/libmunge.so \
      && mkdir -p /run/munge


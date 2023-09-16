FROM ubuntu:20.04
LABEL maintainer="subhasis.bhattacharya@gmail.com"
LABEL description="PBS Professional Open Source dev work env"

ENV DEBIAN_FRONTEND=noninteractive

RUN apt-get -y autoclean \
      && apt-get -y autoremove \
      && apt-get -y update \
      && apt-get -y install software-properties-common sudo wget \
      && apt-get -y autoremove \
      && apt-get -y clean \
      && apt-get -y autoclean \
      && apt-get -y update

# Get all the prerequisites, third parties etc
RUN apt-get -y install gcc make libtool libhwloc-dev libx11-dev \
      libxt-dev libedit-dev libical-dev ncurses-dev perl \
      postgresql-server-dev-all postgresql-contrib python3-dev tcl-dev tk-dev swig \
      libexpat-dev libssl-dev libxext-dev libxft-dev autoconf \
      automake g++ libpam-dev \
      expat libedit2 postgresql python3 postgresql-contrib sendmail-bin \
      tcl tk libical3 munge default-jdk \
      cmake git  libjansson-dev libev-dev wget libhiredis-dev libavro-dev redis librdkafka-dev kafkacat \
      man time csh file lsof curl wget rsync gzip libpam-dev expect dnsutils language-pack-en \
      iputils-ping net-tools vim strace gdb valgrind \
      adduser libfontconfig1 \
      python3-pexpect python3-pytest golang libcurl4-openssl-dev openssh-server \
      && apt-get autoremove

RUN touch /root/.bashrc
RUN echo make install >> /root/.bash_history
RUN echo /opt/pbs/libexec/pbs_postinstall >> /root/.bash_history
RUN echo chmod 4755 /opt/pbs/sbin/pbs_iff /opt/pbs/sbin/pbs_rcp >> /root/.bash_history

# User creation and setup password less ssh
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
      && echo '  PreferredAuthentications publickey,password' >> /etc/ssh/ssh_config

# build and install PBS
COPY . /src/pbspro/
RUN   ls -l /src/pbspro \
      && mkdir -p /tmp/targetpbs \
	&& cd /tmp/targetpbs \
	&& cmake -DCMAKE_INSTALL_PREFIX=/opt/pbs /src/pbspro/ \
	&& cmake --build . --clean-first \
	&& cmake --install . \
	&& chmod 4755 /opt/pbs/sbin/pbs_iff /opt/pbs/sbin/pbs_rcp \
      && echo 'PBS_LOG_HIGHRES_TIMESTAMP=1' >> /etc/pbs.conf \
      && echo 'PBS_COMM_LOG_EVENTS=2047' >> /etc/pbs.conf

RUN mkdir -p /root/certs \
      && cd /root/certs \
      && echo [ cacert ] > /root/certs/my.conf \
      && echo subjectKeyIdentifier = hash >> /root/certs/my.conf \
      && echo authorityKeyIdentifier = keyid:always,issuer  >> /root/certs/my.conf \
      && echo basicConstraints = critical, CA:TRUE  >> /root/certs/my.conf \
      && echo keyUsage = critical, digitalSignature, cRLSign, keyCertSign, keyEncipherment  >> /root/certs/my.conf \
      && echo extendedKeyUsage = clientAuth, serverAuth, emailProtection  >> /root/certs/my.conf \
      && echo nsCertType = server, client, email  >> /root/certs/my.conf \
      && echo nsComment = "CA Certificate Generated By OpenSSL for PBSPro"  >> /root/certs/my.conf \
      && echo [ usrcert ]  >> /root/certs/my.conf \
      && echo subjectKeyIdentifier = hash  >> /root/certs/my.conf \
      && echo authorityKeyIdentifier = keyid:always,issuer:always  >> /root/certs/my.conf \
      && echo basicConstraints = critical, CA:FALSE  >> /root/certs/my.conf \
      && echo keyUsage = critical, nonRepudiation, digitalSignature, keyEncipherment  >> /root/certs/my.conf \
      && echo extendedKeyUsage = clientAuth, serverAuth, emailProtection  >> /root/certs/my.conf \
      && echo nsCertType = server, client, email  >> /root/certs/my.conf \
      && echo nsComment = "User Certificate Generated By OpenSSL for PBSPro" >> /root/certs/my.conf \
      && openssl genrsa -out rootca.key.pem 4096 \
      && openssl req -new -key rootca.key.pem -out rootca.csr.pem -subj "/O=PBSPro/OU=PBSPro/CN=RootCA/" \
      && openssl x509 -req -signkey ./rootca.key.pem -extfile ./my.conf -extensions cacert -days 12775 -in rootca.csr.pem -out rootca.cert.pem \
      && openssl genrsa -out intca.key.pem 4096 \
      && openssl req -new -key intca.key.pem -out intca.csr.pem -subj "/O=PBSPro/OU=PBSPro/CN=IntCA/" \
      && openssl x509 -req -CAkey ./rootca.key.pem -CA ./rootca.cert.pem -CAcreateserial -CAserial ./serials.txt -extfile ./my.conf -extensions cacert -days 9125 -in intca.csr.pem -out intca.cert.pem \
      && cat rootca.cert.pem intca.cert.pem > ca.cert.pem \
      && openssl genrsa -out pbspro.key.pem 2048 \
      && openssl req -new -key pbspro.key.pem -out pbspro.csr.pem -subj "/O=PBSPro/OU=PBSPro/CN=PBSProServices/" \
      && openssl x509 -req -CAkey ./intca.key.pem -CA ./intca.cert.pem -CAcreateserial -CAserial ./serials.txt -extfile ./my.conf -extensions usrcert -days 1825 -in pbspro.csr.pem -out pbspro.cert.pem \
      && openssl verify -CAfile ca.cert.pem pbspro.cert.pem

ENV PATH=/opt/pbs/bin:/opt/pbs/sbin:/opt/ptl/bin:/opt/lqs/bin:$PATH
ENV PYTHONPATH=/opt/ptl/lib/python3.8/site-packages/:/opt/ptl/bin:$PYTHONPATH

COPY --from=tools container-entrypoint /


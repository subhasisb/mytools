FROM centos:8
MAINTAINER Hiren Vadalia <hiren.vadalia@altair.com>
ENV container=docker TERM=xterm TZ=UTC TH_VERSION=3
LABEL maintainer="hiren.vadalia@altair.com"
ENTRYPOINT [ "/container-entrypoint" ]
USER 0
COPY etc/ /workspace/etc
RUN /workspace/etc/install-system-packages centos8
ENV LANGUAGE=C.utf8 LANG=C.utf8
STOPSIGNAL SIGRTMIN+3

FROM opensuse/leap:15
MAINTAINER Hiren Vadalia <hiren.vadalia@altair.com>
ENV container=docker TERM=xterm TZ=UTC TH_VERSION=3 LANGUAGE=en_US.utf8 LANG=en_US.utf8
LABEL maintainer="hiren.vadalia@altair.com"
ENTRYPOINT [ "/container-entrypoint" ]
USER 0
COPY etc/ /workspace/etc
RUN /workspace/etc/install-system-packages opensuse15
STOPSIGNAL SIGRTMIN+3

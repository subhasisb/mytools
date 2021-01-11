FROM ubuntu:18.04
MAINTAINER Hiren Vadalia <hiren.vadalia@altair.com>
LABEL maintainer="hiren.vadalia@altair.com"
ENV container=docker TERM=xterm TZ=UTC TH_VERSION=3 DEBIAN_FRONTEND=noninteractive \
	LANGUAGE=en_US.utf8 LANG=en_US.utf8
ENTRYPOINT [ "/container-entrypoint" ]
USER 0
COPY etc/ /workspace/etc
RUN /workspace/etc/install-system-packages ubuntu18
STOPSIGNAL SIGRTMIN+3

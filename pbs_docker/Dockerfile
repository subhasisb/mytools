FROM centos:7
LABEL maintainer="subhasis.bhattacharya@gmail.com"
LABEL description="PBS Professional Open Source dev work env"

# Get all the prerequisites, third parties etc
RUN yum install -y gcc make rpm-build libtool hwloc-devel \
      libX11-devel libXt-devel libedit-devel libical-devel \
      ncurses-devel perl postgresql-devel postgresql-contrib python3-devel tcl-devel \
      tk-devel swig expat-devel openssl-devel libXext libXft \
      autoconf automake

RUN yum install -y expat libedit postgresql-server postgresql-contrib python3 \
      sendmail sudo tcl tk libical

RUN yum install -y valgrind gdb strace

# install git2
RUN yum remove git*
RUN yum -y install  https://repo.ius.io/ius-release-el7.rpm
RUN yum -y install  git222

# add some routinely needed stuff to the bash history
RUN touch /root/.bashrc
RUN echo ../configure --prefix=/opt/pbs CFLAGS=\"-g -Wall -Werror\" >> /root/.bash_history
RUN echo make install >> /root/.bash_history
RUN echo /opt/pbs/libexec/pbs_postinstall >> /root/.bash_history
RUN echo chmod 4755 /opt/pbs/sbin/pbs_iff /opt/pbs/sbin/pbs_rcp >> /root/.bash_history

ENV PATH="/opt/pbs/bin:/opt/pbs/sbin:${PATH}"


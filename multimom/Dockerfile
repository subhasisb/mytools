FROM centos
ARG JOBUSER=subhas

LABEL maintainer="subhasis.bhattacharya@gmail.com"
LABEL description="PBS mom image, centos + hwloc"

# Get all the prerequisites, third parties etc
RUN yum install -y hwloc valgrind gdb strace
RUN useradd -ms /bin/bash  $JOBUSER


#!/bin/sh
. /etc/pbs.conf
export PATH=$PBS_EXEC/bin:$PATH:/opt/sgi/mpt/mpt-2.02/bin
export LD_LIBRARY_PATH=/opt/sgi/mpt/mpt-2.02/lib


mpiexec /work/tests/multimom/hellompi

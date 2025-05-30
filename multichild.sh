#!/bin/bash

jobs=$1
j=1
source /etc/lqs.conf
while [ $j -le $jobs ]; do
	/opt/pbs/bin/qsub -koe -o /dev/null -e /dev/null -- /bin/true > /dev/null
	j=$(( $j + 1 ))
done


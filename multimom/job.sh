#!/bin/bash -x

echo Server=${XL_SERVER}
echo XL_mom_port=${XL_MOM_PORT}
echo Walltime=${WALLTIME}

echo "Env dump"
env
echo "Env dump end"

pbsdsh_cmd=`which pbsdsh`
dir=`dirname $pbsdsh_cmd`

pbsdsh ${dir}/pbs_xl_mom start ${XL_MOM_PORT} ${XL_SERVER}
sleep ${WALLTIME}
pbsdsh ${dir}/pbs_xl_mom stop

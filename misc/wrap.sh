#!/bin/bash

. ${PBS_CONF_FILE:-/etc/pbs.conf}

if [ $# -lt 1 ]; then
	echo "syntax: $0 <ful-path-pbs-cmd>"
	exit 1
fi

while [ $# -gt 0 ]
do
	cmd=$1
	wrap=val_wrap
	echo Wrapping $cmd 

	if [ -L "${cmd}" ]; then
		echo "Oops! File is already wrapped"
		exit 1
	fi

	if [ ! -f "${cmd}" ]; then
		echo "Oops! binary ${cmd} does not exit!!"
		exit 1
	fi

	cp ./${wrap} ${PBS_EXEC}/bin/${wrap}

	if [ -f "${cmd}" ]; then
		mv ${cmd} ${cmd}.orig
	fi

	ln -s ${PBS_EXEC}/bin/${wrap} ${cmd}

	shift
done

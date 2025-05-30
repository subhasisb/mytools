#!/bin/bash

if [ $# -lt 2 ]; then
	echo "syntax: $0 numjobs numthreads"
	exit 1
fi

numjobs=$1
numthreads=$2
jobs_per_thread=$(( $numjobs / $numthreads ))
rem=$(( $numjobs % $numthreads ))

if [ $rem -gt 0 ]; then
	echo "$numjobs is equally divisible by $numthreads"
	exit 1
fi

declare -a user=("pbsuser1" 
		"pbsuser2" 
		"pbsuser3"
		"pbsuser4"
		"pbsuser5"
		"pbsuser6"
		"pbsuser7"
		)
userlen=${#user[@]}
i=1
while [ $i -le $numthreads ]; do
	userind=$(( $i % $userlen ))
	echo "launching thread $i with $jobs_per_thread jobs as user ${user[$userind]}"
	sudo -u ${user[$userind]} $PWD/multichild.sh $jobs_per_thread &
	#sudo -E -u ${user[$userind]} env | grep PBS &
	i=$(( $i + 1 ))
done
wait




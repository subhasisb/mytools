#!/bin/bash
#


oper=$1
shift

if [ $# -lt 2 ]; then
	echo "syntax: <create/restart/remove> <start-port> <num-comms>"
	exit 1
fi
start_port=$1
num_comms=$2
comm_routers="$HOSTNAME:17001"

port=$start_port
for i in $(seq 1 1 $num_comms)
do
	if [ "$oper" == "create" ]; then
		echo "starting comm with $port"
		cp /etc/pbs.conf /etc/pbscomm_$port.conf
		sed -i "s/PBS_START_SERVER=1/PBS_START_SERVER=0/g" /etc/pbscomm_$port.conf
		sed -i "s/PBS_START_COMM=0/PBS_START_COMM=1/g" /etc/pbscomm_$port.conf
		sed -i "s/PBS_START_SCHED=1/PBS_START_SCHED=0/g" /etc/pbscomm_$port.conf
		sed -i "s/PBS_START_MOM=1/PBS_START_MOM=0/g" /etc/pbscomm_$port.conf
		sed -i "s/PBS_HOME=.*/PBS_HOME=\/var\/spool\/pbscomm_$port/g" /etc/pbscomm_$port.conf
		echo "PBS_COMM_NAME=$HOSTNAME:$port" >> /etc/pbscomm_$port.conf
		echo "PBS_COMM_ROUTERS=$comm_routers" >> /etc/pbscomm_$port.conf

		comm_routers="$comm_routers,$HOSTNAME:$port"
	fi

	if [ "$oper" == "create" -o "$oper" == "restart" ]; then
		if [ -f /etc/pbscomm_$port.conf ]; then	
			PBS_CONF_FILE=/etc/pbscomm_$port.conf /etc/init.d/pbs restart 
		fi
	elif [ "$oper" == "remove" ]; then
		if [ -f /etc/pbscomm_$port.conf ]; then	
			PBS_CONF_FILE=/etc/pbscomm_$port.conf /etc/init.d/pbs stop
		fi
		rm -f /etc/pbscomm_$port.conf
		rm -rf /var/spool/pbscomm_$port
	fi

	port=$((port + 2))
done

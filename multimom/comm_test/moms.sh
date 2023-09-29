#!/bin/bash
#
#


function remove_mom()
{
	PBS_CONF_FILE=$1 /etc/init.d/pbs stop
	rm -f /etc/pbsmom_$port.conf
	rm -rf /var/spool/pbsmom_$port
}


oper=$1
shift

if [ $# -lt 4 ]; then
	echo "syntax: <create/restart/remove> <start-port> <num-moms> <comm-port-start> <num-comms>"
	exit 1
fi
start_port=$1
num_moms=$2
comm_port=$3
num_comms=$4

moms_per_comm=$((num_moms / num_comms))

if [ "$oper" == "create" -o "$oper" == "remove" ]; then
	/opt/pbs/bin/qmgr -c "d n @default"
fi

port=$start_port
count=0
for i in $(seq 1 1 $num_moms)
do
	if [ "$oper" == "create" ]; then
		echo "starting mom with $port"
		cp /etc/pbs.conf /etc/pbsmom_$port.conf
		sed -i "s/PBS_START_SERVER=1/PBS_START_SERVER=0/g" /etc/pbsmom_$port.conf
		sed -i "s/PBS_START_COMM=1/PBS_START_COMM=0/g" /etc/pbsmom_$port.conf
		sed -i "s/PBS_START_SCHED=1/PBS_START_SCHED=0/g" /etc/pbsmom_$port.conf
		sed -i "s/PBS_START_MOM=0/PBS_START_MOM=1/g" /etc/pbsmom_$port.conf
		sed -i "s/PBS_HOME=.*/PBS_HOME=\/var\/spool\/pbsmom_$port/g" /etc/pbsmom_$port.conf
		echo "PBS_MOM_SERVICE_PORT=${port}" >> /etc/pbsmom_$port.conf
		echo "PBS_MANAGER_SERVICE_PORT=$((port + 1))" >> /etc/pbsmom_$port.conf
		echo "PBS_LEAF_ROUTERS=$HOSTNAME:$((comm_port))" >> /etc/pbsmom_$port.conf

		/opt/pbs/bin/qmgr -c "c n pbsmom$port mom=$HOSTNAME,port=$port"
		count=$((count + 1))

		if [ $count -eq $moms_per_comm ]; then
			comm_port=$((comm_port + 2))
			count=0
		fi
	fi

	if [ "$oper" == "create" -o "$oper" == "restart" ]; then
		if [ -f /etc/pbsmom_$port.conf ]; then	
			PBS_CONF_FILE=/etc/pbsmom_$port.conf /etc/init.d/pbs restart &
		fi 
	elif [ "$oper" == "remove" ]; then
		if [ -f /etc/pbsmom_$port.conf ]; then	
			remove_mom /etc/pbsmom_$port.conf &
		fi
	fi

	port=$((port + 2))
done

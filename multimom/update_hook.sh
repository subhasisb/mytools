
if [ $# -ne 2 ]; then
	echo "usage: $0 <hook-file-touches> <delay>"
	exit 1
fi
touches=$1
delay=$2

## Create the hook script file
> tst.py
echo "import pbs" >> tst.py
echo "e=pbs.event()" >> tst.py
echo "e.reject(\"no way Jose\")" >> tst.py

qmgr -c "d h tst1" > /dev/null 2>&1
qmgr -c "d h tst2" > /dev/null 2>&1
hookname="tst1"

cr=0
while :
do
	echo "`date`: Iteration $cr"
	echo "======================================"
	qmgr -c "c h $hookname event=execjob_begin"
	echo "Created hook"

	i=0
	while [ $i -lt $touches ]
	do
		qmgr -c "i h $hookname application/x-python default tst.py"
		i=$(($i+1))	
	done

	qmgr -c "d h $hookname" > /dev/null 2>&1
	echo "Deleted hook"

	if [ "$hookname" = "tst1" ]; then
		hookname="tst2"
	else
		hookname="tst1"
	fi

	if [ $delay -gt 0 ]; then
		sleep $delay
	fi

	cr=$(($cr+1))
done


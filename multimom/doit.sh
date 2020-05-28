while :
do
	./run_moms 1 100 19000
	sleep 2
	#./kill_some.sh
	./stop_all_moms
	ps -ef | grep pbs_mom | wc -l
	sleep 2
done

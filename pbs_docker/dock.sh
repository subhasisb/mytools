status=`docker inspect -f '{{.State.Status}}' devenv`
if [ "$status" = "running" ]; then
	docker-compose exec pbsenv bash
	exit 0
elif [ "$status" = "exited" ]; then
	docker-compose rm -s -f pbsenv
fi
docker-compose up -d pbsenv

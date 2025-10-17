make docker-build
docker rmi subhasisb1974/controller:latest
docker tag controller:latest subhasisb1974/controller:latest
docker push subhasisb1974/controller:latest

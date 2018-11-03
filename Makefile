
docker-build-client:
	docker build -t bulky_client client/

docker-build-server:
	docker build -t bulky_server server/

docker-build: docker-build-client docker-build-server

docker-run:
	docker-compose up


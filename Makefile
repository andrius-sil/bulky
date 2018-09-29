
docker-build-client:
	docker build -t bulky-client client/

docker-build-server:
	docker build -t bulky-server server/

docker-build: docker-build-client docker-build-server

docker-run:
	docker-compose up


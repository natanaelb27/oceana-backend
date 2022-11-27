db-init:
	psql -U postgres -h localhost -d postgres -f db_init.sql
# `docker build` creates docker images from the `Dockerfile`
# -t / --tag used to label the image with string value which readable and recogniseable for human.
# first argument which is "." used to tell which directory that contain Dockerfile to build
docker-build:
	docker build -t oceana-app .
# --publish / -p used to publish a port for our container
# --publish / -p command is [host_port]:[container_port], so if we wanted to expose port 8080 inside container to port 3000 on host then we pass "3000:8000"
# the first argument which is "oceana-app" is the image tag name
# -d / --detach can also be used to run the container in detached mode
docker-run:
	docker run --publish 8080:8080 oceana-app
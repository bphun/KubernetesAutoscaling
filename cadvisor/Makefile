TAG = cadvisor
PREFIX = bphun/k8-auto-scaling
VERSION = v0.36.0

build:
	docker build -t $(PREFIX):$(TAG) .

run: build
	-docker network create k8sAutoscaling
	-docker container rm $(TAG)

	sudo docker run \
	--volume=/:/rootfs:ro \
	--volume=/var/run:/var/run:ro \
	--volume=/sys:/sys:ro \
	--volume=/var/lib/docker/:/var/lib/docker:ro \
	--volume=/dev/disk/:/dev/disk:ro \
	--publish=8080:9094 \
	--detach=true \
	--name=cadvisor \
	--privileged \
	--device=/dev/kmsg \
	--network k8sAutoscaling \
	gcr.io/cadvisor/cadvisor:$(VERSION)

	#  docker run -p 9094:9094 --network k8sAutoscaling --name=$(TAG) -t $(PREFIX):$(TAG)

push: build
	docker push $(PREFIX):$(TAG)
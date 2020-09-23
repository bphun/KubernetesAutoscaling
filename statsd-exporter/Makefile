TAG = statsd
PREFIX = bphun/k8-auto-scaling

build:
	docker build -t $(PREFIX):$(TAG) .

run: build
	-docker network create k8sAutoscaling
	-docker container rm $(TAG)

	 docker run -p 9091:9091 --network k8sAutoscaling --name=$(TAG) -t $(PREFIX):$(TAG)

push: build
	docker push $(PREFIX):$(TAG)
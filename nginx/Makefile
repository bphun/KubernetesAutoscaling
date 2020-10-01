TAG=nginx
PREFIX=bphun/k8-auto-scaling

build:
	docker build -t $(PREFIX):$(TAG) .

run: build
	-docker network create k8sAutoscaling
	-docker container rm $(TAG)
	
	 docker run -p 7080:7080 -p 7443:7443 -p 7081:7081 --network k8sAutoscaling --name $(TAG) -t $(PREFIX):$(TAG)

push: build
	docker push $(PREFIX):$(TAG)
TAG = nginx-exporter
PREFIX = public.ecr.aws/h6v8s1y1/k8s-autoscaling

build:
	docker build -t $(PREFIX):$(TAG) .

run: build
	-docker network create k8sAutoscaling
	-docker container rm $(TAG)

	docker run -p 9093:9093 --network k8sAutoscaling --name $(TAG) -t $(PREFIX):$(TAG)
	
push: build
	docker push $(PREFIX):$(TAG)
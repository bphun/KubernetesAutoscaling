TAG = transaction-api
PREFIX = public.ecr.aws/h6v8s1y1/k8s-autoscaling

build: clean
	docker build -t $(PREFIX):$(TAG) .

run: build
	-docker network create k8sAutoscaling
	-docker container rm $(TAG)
	
	 docker run -p 5001:5001 --network k8sAutoscaling --name=$(TAG) -t $(PREFIX):$(TAG)

transaction-api: clean
	go mod vendor
	go build

clean:
	rm -f main

push: build
	docker push $(PREFIX):$(TAG)
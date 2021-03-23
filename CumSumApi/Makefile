TAG = cum-sum-api
PREFIX = public.ecr.aws/h6v8s1y1/k8s-autoscaling

build: clean
	docker build -t $(PREFIX):$(TAG) .

run: build
	-docker network create k8sAutoscaling
	-docker container rm $(TAG)
	
	 docker run -p 5000:5000 --network k8sAutoscaling --name=$(TAG) -t $(PREFIX):$(TAG)

api-golang: clean
	go get -u github.com/valyala/fasthttp
	go get -u github.com/valyala/fasthttp/reuseport
	go get -u gopkg.in/alexcesaro/statsd.v2
	go get -u github.com/fasthttp/router
	go build

clean:
	rm -f main

push: build
	docker push $(PREFIX):$(TAG)
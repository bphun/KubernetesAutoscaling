TAG = api-golang
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
	go get -u github.com/cactus/go-statsd-client/statsd
	go build

clean:
	rm -f api-golang

push: build
	docker push $(PREFIX):$(TAG)
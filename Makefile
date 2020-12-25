api-build:
	cd api/ && \
	make build

grafana-build:
	cd grafana && \
	make build

nginx-build:
	cd nginx/ && \
	make build

nginx-exporter-build:
	cd nginx-exporter/ && \
	make build

prometheus-build:
	cd prometheus/ && \
	make build

statsd-exporter-build:
	cd statsd-exporter/ && \
	make build

api-push:
	cd api/ && \
	make push

grafana-push:
	cd grafana && \
	make push

nginx-push:
	cd nginx/ && \
	make push

nginx-exporter-push:
	cd nginx-exporter/ && \
	make push

prometheus-push:
	cd prometheus/ && \
	make push

statsd-exporter-push:
	cd statsd-exporter/ && \
	make push

build: api-build grafana-build nginx-build nginx-exporter-build prometheus-build statsd-exporter-build
push: api-push grafana-push nginx-push nginx-exporter-push prometheus-push statsd-exporter-push
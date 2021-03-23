cum-sum-api-build:
	cd CumSumApi/ && \
	make build

transaction-api-build:
	cd TransactionAPI/ && \
	make build

transaction-db-build:
	cd TransactionDB/ && \
	make build

mongodb-exporter-build:
	cd mongodb-exporter/ && \
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

mongodb-exporter-push:
	cd mongodb-exporter/ && \
	make push

cum-sum-api-push:
	cd CumSumApi/ && \
	make push
	
transaction-api-push:
	cd TransactionAPI/ && \
	make push

transaction-db-push:
	cd TransactionDB/ && \
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

build: cum-sum-api-build transaction-api-build transaction-db-build mongodb-exporter-build grafana-build nginx-build nginx-exporter-build prometheus-build statsd-exporter-build
push: cum-sum-api-push transaction-api-push transaction-db-push mongodb-exporter-push grafana-push nginx-push nginx-exporter-push prometheus-push statsd-exporter-push
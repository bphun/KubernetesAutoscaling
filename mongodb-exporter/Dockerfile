FROM alpine:latest

ENV MONGODB_EXPORTER_VERSION 0.20.3

RUN apk update \
    && apk --no-cache add curl \
    && curl -LO https://github.com/percona/mongodb_exporter/releases/download/v${MONGODB_EXPORTER_VERSION}/mongodb_exporter-${MONGODB_EXPORTER_VERSION}.linux-amd64.tar.gz \
    && tar -xvzf mongodb_exporter-${MONGODB_EXPORTER_VERSION}.linux-amd64.tar.gz

FROM alpine:latest

COPY --from=0 mongodb_exporter .

# CMD ./mongodb_exporter-exporter -nginx.scrape-uri http://nginx.nginx:7081/nginx_status -nginx.retries 60 -nginx.retry-interval 10s -web.listen-address :9093
CMD ./mongodb_exporter --mongodb.uri=mongodb://127.0.0.1:27017 --log.level="debug"
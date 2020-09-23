FROM alpine:latest

ENV NGINX_EXPORTER_VERSION 0.8.0

RUN apk update \
    && apk --no-cache add curl \
    && curl -LO https://github.com/nginxinc/nginx-prometheus-exporter/releases/download/v${NGINX_EXPORTER_VERSION}/nginx-prometheus-exporter-${NGINX_EXPORTER_VERSION}-linux-amd64.tar.gz \
    && tar -xvzf nginx-prometheus-exporter-${NGINX_EXPORTER_VERSION}-linux-amd64.tar.gz

FROM alpine:latest

COPY --from=0 nginx-prometheus-exporter .

CMD ./nginx-prometheus-exporter -nginx.scrape-uri http://nginx.nginx:7081/nginx_status -nginx.retries 60 -nginx.retry-interval 10s -web.listen-address :9093
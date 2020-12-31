#!/bin/bash

kubectl apply -f k8s/namespaces/namespaces.yaml

kubectl apply -f k8s/volumes/grafana-volume.yaml
kubectl apply -f k8s/volumes/prometheus-volume.yaml
kubectl apply -f k8s/volumes/api-service-volume.yaml

kubectl apply -f k8s/monitoring/service-account.yaml 
kubectl apply -f k8s/monitoring/node-exporter.yaml 
kubectl apply -f k8s/monitoring/cadvisor.yaml 
kubectl apply -f k8s/monitoring/prometheus-adapter/config-map.yaml 
# kubectl apply -f k8s/monitoring/prometheus-adapter/deployment.yaml 
# kubectl apply -f k8s/monitoring/prometheus-adapter/service.yaml 
kubectl apply -f k8s/api/api.yaml
kubectl apply -f k8s/monitoring/grafana.yaml 
kubectl apply -f k8s/monitoring/prometheus-scrapers.yaml 
kubectl apply -f k8s/nginx/nginx.yaml 

helm install api-autoscaling \
--set prometheus.url=http://prometheus-scrapers.monitoring \
--set prometheus.path=/prometheus \
--set rules.existing=adapter-config \
--set logLevel=10 \
--set metricsRelistInterval=30s \
prometheus-community/prometheus-adapter
# --set rules.default=False \

kubectl apply -f k8s/horizontal-autoscalers/api.yaml 
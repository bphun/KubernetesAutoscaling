#!/bin/bash

kubectl apply -f k8s/namespaces/namespaces.yaml

kubectl apply -f k8s/volumes/grafana-volume.yaml
kubectl apply -f k8s/volumes/prometheus-volume.yaml
kubectl apply -f k8s/volumes/api-service-volume.yaml

kubectl create secret generic -n metallb-system memberlist --from-literal=secretkey="$(openssl rand -base64 128)"
kubectl apply -f k8s/metal-lb/metal-lb.yaml
kubectl apply -f k8s/monitoring/service-account.yaml 
kubectl apply -f k8s/monitoring/node-exporter.yaml 
kubectl apply -f k8s/monitoring/cadvisor.yaml 
kubectl apply -f k8s/monitoring/prometheus-adapter/config-map.yaml 
kubectl apply -f k8s/api/api.yaml
kubectl apply -f k8s/monitoring/grafana.yaml 
kubectl apply -f k8s/monitoring/prometheus.yaml 
kubectl apply -f k8s/nginx/nginx.yaml 
kubectl apply -f k8s/ingress/ingress.yaml

helm install api-autoscaling \
--set prometheus.url=http://prometheus.monitoring \
--set prometheus.path=/prometheus \
--set rules.existing=adapter-config \
--set logLevel=10 \
--set metricsRelistInterval=30s \
prometheus-community/prometheus-adapter
# --set rules.default=False \

helm install nginx-ingress ingress-nginx/ingress-nginx
kubectl apply -f k8s/horizontal-autoscalers/api.yaml 
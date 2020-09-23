#!/bin/bash

kubectl apply -f k8s/namespaces/namespaces.yaml

kubectl apply -f k8s/volumes/grafana-volume.yaml
kubectl apply -f k8s/volumes/prometheus-volume.yaml
kubectl apply -f k8s/volumes/api-service-volume.yaml

kubectl apply -f k8s/api/api.yaml
kubectl apply -f k8s/nginx/nginx.yaml 
kubectl apply -f k8s/monitoring/grafana.yaml 
kubectl apply -f k8s/monitoring/prometheus-scrapers.yaml 
kubectl apply -f k8s/monitoring/cadvisor.yaml 
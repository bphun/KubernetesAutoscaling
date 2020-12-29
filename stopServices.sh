#!/bin/bash

kubectl delete -f k8s/nginx/nginx.yaml 
kubectl delete -f k8s/api/api.yaml
kubectl delete -f k8s/monitoring/grafana.yaml 
kubectl delete -f k8s/monitoring/prometheus-scrapers.yaml 
kubectl delete -f k8s/monitoring/cadvisor.yaml 
kubectl delete -f k8s/monitoring/prometheus-adapter/config-map.yaml 
kubectl delete -f k8s/monitoring/prometheus-adapter/deployment.yaml 
kubectl delete -f k8s/monitoring/prometheus-adapter/service.yaml 

kubectl delete -f k8s/namespaces/namespaces.yaml

kubectl delete -f k8s/volumes/grafana-volume.yaml
kubectl delete -f k8s/volumes/prometheus-volume.yaml
kubectl delete -f k8s/volumes/api-service-volume.yaml
kubectl delete -f k8s/monitoring/serviceaccount.yaml 

kubectl delete -f k8s/horizontal-autoscalers/api.yaml 
#!/bin/bash

kubectl delete -f k8s/nginx/nginx.yaml 
kubectl delete -f k8s/api/api.yaml
kubectl delete -f k8s/monitoring/grafana.yaml 
kubectl delete -f k8s/monitoring/prometheus.yaml 
kubectl delete -f k8s/monitoring/cadvisor.yaml 
kubectl delete -f k8s/monitoring/prometheus-adapter/config-map.yaml 

kubectl delete -f k8s/namespaces/namespaces.yaml

kubectl delete -f k8s/volumes/grafana-volume.yaml
kubectl delete -f k8s/volumes/prometheus-volume.yaml
kubectl delete -f k8s/volumes/api-service-volume.yaml

kubectl delete -f k8s/horizontal-autoscalers/cum-sum-api.yaml 
kubectl delete -f k8s/horizontal-autoscalers/transaction-api.yaml 
kubectl delete -f k8s/metal-lb/metal-lb.yaml
kubectl delete -f k8s/ingress/ingress.yaml

kubectl delete -n monitoring -f https://raw.githubusercontent.com/jaegertracing/jaeger-operator/master/deploy/crds/jaegertracing.io_jaegers_crd.yaml
kubectl delete -n monitoring -f https://raw.githubusercontent.com/jaegertracing/jaeger-operator/master/deploy/service_account.yaml
kubectl delete -n monitoring -f https://raw.githubusercontent.com/jaegertracing/jaeger-operator/master/deploy/role.yaml
kubectl delete -n monitoring -f https://raw.githubusercontent.com/jaegertracing/jaeger-operator/master/deploy/role_binding.yaml
kubectl delete -n monitoring -f https://raw.githubusercontent.com/jaegertracing/jaeger-operator/master/deploy/operator.yaml
kubectl delete -f https://raw.githubusercontent.com/jaegertracing/jaeger-operator/master/deploy/cluster_role.yaml
kubectl delete -f https://raw.githubusercontent.com/jaegertracing/jaeger-operator/master/deploy/cluster_role_binding.yaml

helm uninstall api-autoscaling
helm uninstall nginx-ingress
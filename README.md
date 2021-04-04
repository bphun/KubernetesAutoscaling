# Kubernetes Autoscaling

This is a Kubernetes cluster that horizontally scales a Golang REST API based on how many requests per second (RPS) the API is processing. Currently, the `HorizontalPodAutoscaler` (HPA) threshold is set to 3000 RPS. The API simply takes a JSON with an array and returns the cumulative sum of the array as a JSON response.

## Running the Cluster
I ran the cluster on a Ubuntu 20.04.1 system with Kubernetes bootstrapped using kubeadm 1.20.1 but you can boostrap the cluster however you want. 

To start the services, initialize the Kubernetes cluster then run the following commands.
```
git clone https://github.com/bphun/k8AutoScalingTest.git
cd k8AutoScalingTest
./startServices.sh
```

## Accessing the cluster
Run `kubectl --namespace default get services -o wide -w nginx-ingress-ingress-nginx-controller`. You will see a similar output:
```
NAME                                     TYPE           CLUSTER-IP      EXTERNAL-IP    PORT(S)                      AGE   SELECTOR
nginx-ingress-ingress-nginx-controller   LoadBalancer   10.110.227.68   203.0.113.10   80:31757/TCP,443:31344/TCP   29s   app.kubernetes.io/component=controller,app.kubernetes.io/instan
```

This will give you the external IP address of the Nginx ingress controller (203.0.113.10 in this example). You can use the external IP address to access the API from a node within the Kubernetes cluster by running:
```
curl -X POST --data '{"arr": [1, 2, 3, 4, 5, 5, 6, 7, 8, 10], "mode": "man"}' http://$EXTERNAL_IP/api/ | jq
```

In addition to the API service, the Prometheus service can be accessed at `http://$EXTERNAL_IP/prometheus/` and a Grafana dashboard can be accessed at `http://$EXTERNAL_IP/grafana_dashboard/`. The required dashboards are available as JSON files you can import in the `grafana/` directory.

If you want to access the API service from your LAN, you can use an SSH tunnel like so:
```
ssh -L $MASTER_NODE_IP:8000:$EXTERNAL_IP:80 $USER@$MASTER_NODE_IP
```
This will allow you to access the API service from any computer in your LAN through `$MASTER_NODE_IP:8000`

## Load testing the API
I used [Hey](https://github.com/rakyll/hey) to load test the API, which can trigger the Kubernetes HPA to scale the API services. To run Hey, use the following command:
```
hey -m POST -H "Content-Type: application/json" -n 10000 -c 100 -D api/input.json http://$EXTERNAL_IP/api/ 
```
The `-n` argument is the number of requests that will be sent and `-c` is the number of concurrent requests that will be sent by Hey. Increasing `-n` will trigger the HPA to horizontally scale the API service.

If you want to change the HPA threshold, change the `averageValue` value in `k8s/horizontal-autoscalers/api.yaml` and reload the services using `./startServices.sh`

## Update: 04/03/2021
* Added MongoDB database to store request history. Request history may be used in future for custom data aggregation UI
* Implement gRPC API to allow "Cum-Sum-API" to store requests in MongoDB database
* Use [Jaeger](https://www.jaegertracing.io) to add distributed tracing support
* Add ElasticSearch database to store Jaeger spans

## API Architecture

|--------|    HTTP    |---------------|    gRPC    |------------------|   HTTP    |-----------|
|  User  |  ------->  |  Cum-Sum-API  |  ------->  |  TransactionAPI  |  ------>  |  MongoDB  |  
|--------|            |---------------|            |------------------|           |-----------|

### Distributed tracing pipeline:

                                                                              |-------------|
                                                                              |  Jaeger UI  |
                                                                              |-------------|

                                                                                    |
                                                                                    |
                                                                                    V

  |------------|        |----------------|      |--------------------|      |-----------------|
  |  REST API  |  --->  |  Jaeger Agent  | ---> |  Jaeger Collector  | ---> |  Elasticsearch  |
  |------------|        |----------------|      |--------------------|      |-----------------|

## Future features
* I already have Nginx exporter setup so I would like to add a HPA controller for that service
* Add a VerticalPodAutoscaler (VPA) for the API and Nginx services using Node exporter
* ~~Add another "Microservice" so I can test out distributed tracing frameworks like OpenTracing~~ (*Update: 04/03/2021*)
* Add improved logging to the Go API service
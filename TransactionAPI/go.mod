module github.com/bphun/KubernetesAutoscaling/TransactionAPI

go 1.16

require (
	github.com/bphun/KubernetesAutoscaling/Tracing v0.0.0-20210324211613-db7a49dc633a
	github.com/grpc-ecosystem/go-grpc-middleware v1.2.2
	github.com/grpc-ecosystem/go-grpc-prometheus v1.2.0
	github.com/opentracing/opentracing-go v1.2.0
	github.com/prometheus/client_golang v1.10.0
	github.com/uber/jaeger-client-go v2.25.0+incompatible // indirect
	github.com/uber/jaeger-lib v2.4.0+incompatible // indirect
	// github.com/bphun/KubernetesAutoscaling/TransactionAPI v0.0.0-20210321215857-a570db2c029f
	go.mongodb.org/mongo-driver v1.5.0
	google.golang.org/grpc v1.36.0
	google.golang.org/protobuf v1.26.0
)

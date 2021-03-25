module github.com/bphun/KubernetesAutoscaling/CumSumApi

go 1.16

require (
	github.com/bphun/KubernetesAutoscaling/Tracing v0.0.0-20210324211738-07c8ff7bb59a // indirect
	github.com/bphun/KubernetesAutoscaling/TransactionAPI v0.0.0-20210324205732-3557535e53db
	// github.com/bphun/KubernetesAutoscaling/TransactionAPI v0.0.0-20210321215857-a570db2c029f
	github.com/cactus/go-statsd-client/v5 v5.0.0
	github.com/fasthttp/router v1.3.9
	github.com/grpc-ecosystem/go-grpc-middleware v1.2.2
	github.com/grpc-ecosystem/go-grpc-prometheus v1.2.0 // indirect
	github.com/opentracing/opentracing-go v1.2.0
	github.com/prometheus/client_golang v1.10.0 // indirect
	github.com/uber/jaeger-client-go v2.25.0+incompatible
	github.com/uber/jaeger-lib v2.4.0+incompatible
	github.com/valyala/fasthttp v1.22.0
	google.golang.org/grpc v1.36.0
)

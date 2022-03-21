module github.com/bphun/KubernetesAutoscaling/CumSumApi

go 1.16

require (
	github.com/HdrHistogram/hdrhistogram-go v1.1.2 // indirect
	github.com/bphun/KubernetesAutoscaling/Tracing v0.0.0-20210325065540-e66a2f7bfc2c
	github.com/bphun/KubernetesAutoscaling/TransactionAPI v0.0.0-20210324205732-3557535e53db
	// github.com/bphun/KubernetesAutoscaling/TransactionAPI v0.0.0-20210321215857-a570db2c029f
	github.com/cactus/go-statsd-client/v5 v5.0.0
	github.com/fasthttp/router v1.3.9
	github.com/grpc-ecosystem/go-grpc-middleware v1.2.2
	github.com/opentracing/opentracing-go v1.2.0
	github.com/valyala/fasthttp v1.34.0
	google.golang.org/grpc v1.36.0
)

module github.com/bphun/k8sAutoscaling/CumSumApi

go 1.16

require (
	// github.com/bphun/KubernetesAutoscaling/TransactionAPI v0.0.0-00010101000000-000000000000
	// bphun/k8sAutoscaling/TransactionAPI v0.0.0-00010101000000-000000000000
	github.com/cactus/go-statsd-client/v5 v5.0.0
	github.com/fasthttp/router v1.3.9
	github.com/valyala/fasthttp v1.22.0
	google.golang.org/grpc v1.36.0
)

// replace bphun/k8sAutoscaling/TransactionAPI => ../TransactionAPI

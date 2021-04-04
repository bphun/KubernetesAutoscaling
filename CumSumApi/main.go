package main

import (
	"context"
	"encoding/json"
	"flag"
	"log"
	"sync"
	"time"

	tracing "github.com/bphun/KubernetesAutoscaling/Tracing"
	"github.com/cactus/go-statsd-client/v5/statsd"
	"github.com/fasthttp/router"
	grpc_opentracing "github.com/grpc-ecosystem/go-grpc-middleware/tracing/opentracing"
	opentracing "github.com/opentracing/opentracing-go"
	"github.com/valyala/fasthttp"
	"google.golang.org/grpc"

	"github.com/valyala/fasthttp/reuseport"

	"github.com/grpc-ecosystem/go-grpc-middleware"

	pb "github.com/bphun/KubernetesAutoscaling/TransactionAPI/TransactionAPI"
	// _ "net/http/pprof"
)

type InputData struct {
	Arr  []int32 `json:"arr"`
	Mode string  `json:"mode"`
}

type OutputData struct {
	Message          string  `json:"message"`
	Data             []int32 `json:"data"`
	ExecutionTime    int64   `json:"exec_ns"`
	RequestStartTime uint32  `json:"start_time"`
}
type ClientStats struct {
	NumRequests   int64
	Duration      int64
	DurationCount int64
	DurationSum   int64
}

const (
	DEFAULT_HTTP_ADDR = ":8080"
	STATSD_ADDR       = "127.0.0.1:9125"
	DEFAULT_GRPC_ADDR = "localhost:5001"
)

var (
	httpAddr = flag.String("httpAddr", DEFAULT_HTTP_ADDR, "TCP address to listen to")

	statsdConfig = &statsd.ClientConfig{
		Address:       STATSD_ADDR,
		Prefix:        "api",
		ResInterval:   time.Minute,
		UseBuffered:   true,
		FlushInterval: 300 * time.Millisecond,
	}
	statsdClientInstance    statsd.Statter
	statsdClientInstanceErr error
	clientStats             ClientStats
	statsdOnce              sync.Once

	grpcAddr           = flag.String("grpcAddr", DEFAULT_GRPC_ADDR, "Address of Transaction gRPC server")
	grpcInstanceClient pb.TransactionAPIClient
	grpcInstanceConn   *grpc.ClientConn
	grpcInstanceError  error
	grpcOnce           sync.Once
)

func main() {
	flag.Parse()

	requestRouter := router.New()
	requestRouter.POST("/", postRequestHook(processCumSumRequest))

	listener, httpListenErr := reuseport.Listen("tcp4", *httpAddr)

	if httpListenErr != nil {
		log.Fatalf("Error in reuseport listener: %s", httpListenErr)
	}

	defer listener.Close()
	log.Printf("Listening on %s", *httpAddr)

	if err := fasthttp.Serve(listener, requestRouter.Handler); err != nil {
		log.Fatalf("Error in Serve: %s", err)
	}

	listener.Close()
	grpcInstanceConn.Close()
	statsdClientInstance.Close()
}

func getStatsdClient(ctx opentracing.SpanContext) (statsd.Statter, error) {
	statsdOnce.Do(func() {
		tracer := opentracing.GlobalTracer()
		getStatsdClientSpan := tracer.StartSpan("http.getStatsdClient", opentracing.ChildOf(ctx))
		log.Printf("Connecting to statsD server at %s", STATSD_ADDR)

		statsdClientInstance, statsdClientInstanceErr = statsd.NewClientWithConfig(statsdConfig)
		if statsdClientInstanceErr != nil {
			log.Fatalf("Unable to connect to statsD: %v", statsdClientInstanceErr)
		}
		log.Printf("Connected to statsD server at %s", STATSD_ADDR)
		getStatsdClientSpan.Finish()
	})

	return statsdClientInstance, statsdClientInstanceErr
}

func getGrpcClient() (pb.TransactionAPIClient, error) {
	grpcOnce.Do(func() {
		log.Printf("Connecting to gRPC server at %s", *grpcAddr)

		tracer, closer, err := tracing.NewTracer()
		defer closer.Close()
		if err != nil {
			log.Fatalf("Error: %v", err)
		}
		opentracing.SetGlobalTracer(tracer)

		grpcInstanceConn, grpcInstanceError = grpc.Dial(*grpcAddr,
			grpc.WithInsecure(),
			grpc.WithBlock(),
			grpc.WithStreamInterceptor(grpc_middleware.ChainStreamClient(
				grpc_opentracing.StreamClientInterceptor(grpc_opentracing.WithTracer(tracer)),
			)),
			grpc.WithUnaryInterceptor(grpc_middleware.ChainUnaryClient(
				grpc_opentracing.UnaryClientInterceptor(grpc_opentracing.WithTracer(tracer)),
			)),
		)
		if grpcInstanceError != nil {
			log.Fatalf("Error: %v", grpcInstanceError)
		}
		log.Printf("Connected to gRPC server at %s", *grpcAddr)

		grpcInstanceClient = pb.NewTransactionAPIClient(grpcInstanceConn)
	})

	return grpcInstanceClient, grpcInstanceError
}

func postRequestHook(h fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		tracer := opentracing.GlobalTracer()
		postRequestHookSpan := tracer.StartSpan("http.postRequestHook")
		statsdClient, _ := getStatsdClient(postRequestHookSpan.Context())

		t1 := time.Now()
		h(ctx)
		t2 := time.Now()
		duration := t2.Sub(t1)

		clientStats.NumRequests++

		statsdClient.TimingDuration("requests.duration", duration, 1.0)
		statsdClient.Inc("requests.count", clientStats.NumRequests, 1.0)
		postRequestHookSpan.Finish()
	}
}

func updateTransactionHistory(spanCtx opentracing.SpanContext, inArr []int32, outArr []int32, requestStartTime uint32, executionTime int64) {
	tracer := opentracing.GlobalTracer()
	updateTransactionHistorySpan := tracer.StartSpan("http.updateTransactionHistory", opentracing.ChildOf(spanCtx))
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	grpcClient, _ := getGrpcClient()
	r, err := grpcClient.SaveTransaction(ctx, &pb.TransactionRequest{InArr: inArr, OutArr: outArr, StartTime: requestStartTime, ExecTime: executionTime})

	if err != nil {
		log.Fatalf("Could not update transaction history: %v", err)
	}
	log.Printf("TransactionDB: %s", r.GetMessage())
	updateTransactionHistorySpan.Finish()
}

func processCumSumRequest(ctx *fasthttp.RequestCtx) {
	var postBody InputData
	var output OutputData
	var cumSumArr []int32
	var executionTime int64
	var requestStartTime uint32
	message := "success"
	tracer := opentracing.GlobalTracer()
	span := tracer.StartSpan("http.ProcessCumSumRequest")

	err := json.Unmarshal(ctx.PostBody(), &postBody)
	if err != nil {
		message = err.Error()
		cumSumArr = make([]int32, 0)
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
	} else {
		spanContext := span.Context()
		origArray := make([]int32, len(postBody.Arr))
		cumSumArr = postBody.Arr
		copy(origArray, cumSumArr)
		requestStartTime, executionTime = cumSum(spanContext, cumSumArr)

		go updateTransactionHistory(spanContext, origArray, cumSumArr, requestStartTime, executionTime)
	}

	output = OutputData{
		Message:          message,
		Data:             cumSumArr,
		ExecutionTime:    executionTime,
		RequestStartTime: requestStartTime,
	}

	marshalledJSON, _ := json.Marshal(output)
	ctx.SetContentType("application/json")
	ctx.SetBody(marshalledJSON)
	span.Finish()
}

func cumSum(ctx opentracing.SpanContext, arr []int32) (uint32, int64) {
	tracer := opentracing.GlobalTracer()
	span := tracer.StartSpan("http.cumSumAlgo", opentracing.ChildOf(ctx))
	defer span.Finish()

	start := time.Now()
	for i := 1; i < len(arr); i++ {
		arr[i] += arr[i-1]
	}
	duration := time.Since(start)

	return uint32(start.Unix()), duration.Nanoseconds()
}

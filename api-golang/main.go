package main

import (
	"context"
	"encoding/json"
	"flag"
	"log"
	"sync"
	"time"

	"github.com/cactus/go-statsd-client/v5/statsd"
	"github.com/fasthttp/router"
	"github.com/valyala/fasthttp"
	"google.golang.org/grpc"

	"github.com/valyala/fasthttp/reuseport"

	// pb "bphun/k8sAutoscaling/TransactionAPI/TransactionAPI"
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
	GRPC_ADDR         = "localhost:50051"
)

var (
	httpAddr = flag.String("addr", DEFAULT_HTTP_ADDR, "TCP address to listen to")

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

	if err := fasthttp.Serve(listener, requestRouter.Handler); err != nil {
		log.Fatalf("Error in Serve: %s", err)
	}

	listener.Close()
	grpcInstanceConn.Close()
	statsdClientInstance.Close()
}

func getStatsdClient() (statsd.Statter, error) {
	statsdOnce.Do(func() {
		statsdClientInstance, statsdClientInstanceErr = statsd.NewClientWithConfig(statsdConfig)
		if statsdClientInstanceErr != nil {
			log.Fatalf("Unable to connect to statsD: %v", statsdClientInstanceErr)
		}
		log.Printf("Connected to statsD server at %s", STATSD_ADDR)
	})

	return statsdClientInstance, statsdClientInstanceErr
}

func getGrpcClient() (pb.TransactionAPIClient, error) {
	grpcOnce.Do(func() {
		grpcInstanceConn, grpcInstanceError = grpc.Dial(GRPC_ADDR, grpc.WithInsecure(), grpc.WithBlock())
		if grpcInstanceError != nil {
			log.Fatalf("Error: %v", grpcInstanceError)
		}

		log.Printf("Connected to gRPC server at %s", GRPC_ADDR)

		grpcInstanceClient = pb.NewTransactionAPIClient(grpcInstanceConn)
	})

	return grpcInstanceClient, grpcInstanceError
}

func postRequestHook(h fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		statsdClient, _ := getStatsdClient()

		t1 := time.Now()
		h(ctx)
		t2 := time.Now()
		duration := t2.Sub(t1)

		clientStats.NumRequests++

		statsdClient.TimingDuration("requests.duration", duration, 1.0)
		statsdClient.Inc("requests.count", clientStats.NumRequests, 1.0)
	}
}

func updateTransactionHistory(inArr []int32, outArr []int32, requestStartTime uint32, executionTime int64) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	grpcClient, _ := getGrpcClient()
	r, err := grpcClient.SaveTransaction(ctx, &pb.TransactionRequest{InArr: inArr, OutArr: outArr, StartTime: requestStartTime, ExecTime: executionTime})

	if err != nil {
		log.Fatalf("Could not update transaction history: %v", err)
	}
	log.Printf("TransactionDB: %s", r.GetMessage())
}

func processCumSumRequest(ctx *fasthttp.RequestCtx) {
	var postBody InputData
	var output OutputData
	var cumSumArr []int32
	var executionTime int64
	var requestStartTime uint32
	message := "success"

	err := json.Unmarshal(ctx.PostBody(), &postBody)
	if err != nil {
		message = err.Error()
		cumSumArr = make([]int32, 0)
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
	} else {
		origArray := make([]int32, len(postBody.Arr))
		cumSumArr = postBody.Arr

		copy(origArray, cumSumArr)

		requestStartTime, executionTime = cumSum(cumSumArr)

		go updateTransactionHistory(origArray, cumSumArr, requestStartTime, executionTime)
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
}

func cumSum(arr []int32) (uint32, int64) {
	start := time.Now()
	for i := 1; i < len(arr); i++ {
		arr[i] += arr[i-1]
	}
	duration := time.Since(start)

	return uint32(start.Unix()), duration.Nanoseconds()
}

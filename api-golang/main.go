package main

import (
	"context"
	"encoding/json"
	"flag"
	"log"
	"time"

	"github.com/cactus/go-statsd-client/v5/statsd"
	"github.com/fasthttp/router"
	"github.com/valyala/fasthttp"

	"github.com/valyala/fasthttp/reuseport"

	pb "bphun/k8sAutoscaling/TransactionAPI/TransactionAPI"
	// _ "net/http/pprof"

	"google.golang.org/grpc"
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
	address    = "localhost:50051"
	httpAddr   = ":8080"
	statsdAddr = "127.0.0.1:9125"
)

var (
	addr         = flag.String("addr", httpAddr, "TCP address to listen to")
	grpcClient   pb.TransactionAPIClient
	statsdConfig = &statsd.ClientConfig{
		Address:       statsdAddr,
		Prefix:        "api",
		ResInterval:   time.Minute,
		UseBuffered:   true,
		FlushInterval: 300 * time.Millisecond,
	}
	statsdClient, statsdClientErr = statsd.NewClientWithConfig(statsdConfig)
	clientStats                   ClientStats
)

func main() {
	flag.Parse()

	if statsdClientErr != nil {
		log.Fatal("Unable to connect to statsD: %v", statsdClientErr)
	}
	log.Printf("Connected to statsd server at %s", statsdAddr)

	defer statsdClient.Close()

	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	log.Printf("Connected to gRPC server at %s", address)

	defer conn.Close()
	grpcClient = pb.NewTransactionAPIClient(conn)

	requestRouter := router.New()
	requestRouter.POST("/", PostRequestHook(processCumSumRequest))

	listener, httpListenErr := reuseport.Listen("tcp4", *addr)

	if httpListenErr != nil {
		log.Fatalf("Error in reuseport listener: %s", httpListenErr)
	}

	defer listener.Close()

	if err := fasthttp.Serve(listener, requestRouter.Handler); err != nil {
		log.Fatalf("Error in Serve: %s", err)
	}
}

func PostRequestHook(h fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		t1 := time.Now()
		h(ctx)
		t2 := time.Now()
		duration := t2.Sub(t1)

		clientStats.NumRequests++

		statsdClient.TimingDuration("requests.duration", duration, 1.0)
		statsdClient.Inc("requests.count", clientStats.NumRequests, 1.0)
	}
}

func processCumSumRequest(ctx *fasthttp.RequestCtx) {
	var postBody InputData
	var output OutputData
	err := json.Unmarshal(ctx.PostBody(), &postBody)

	if err != nil {
		output = OutputData{
			Message:       err.Error(),
			Data:          make([]int32, 0),
			ExecutionTime: 0,
		}
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
	} else {
		origArray := make([]int32, len(postBody.Arr))
		cumSumArr := postBody.Arr
		copy(origArray, postBody.Arr)
		requestStartTime, executionTime := cumSum(cumSumArr)

		output = OutputData{
			Message:          "success",
			Data:             cumSumArr,
			ExecutionTime:    executionTime,
			RequestStartTime: requestStartTime,
		}

		go updateTransactionHistory(origArray, cumSumArr, requestStartTime, executionTime)
	}

	marshalledJSON, _ := json.Marshal(output)
	ctx.SetContentType("application/json")
	ctx.SetBody(marshalledJSON)
}

func updateTransactionHistory(inArr []int32, outArr []int32, requestStartTime uint32, executionTime int64) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := grpcClient.SaveTransaction(ctx, &pb.TransactionRequest{InArr: inArr, OutArr: outArr, StartTime: requestStartTime, ExecTime: executionTime})
	if err != nil {
		log.Fatalf("Could not update transaction history: %v", err)
	}
	log.Printf("TransactionDB: %s", r.GetMessage())
}

func cumSum(arr []int32) (uint32, int64) {
	start := time.Now()
	for i := 1; i < len(arr); i++ {
		arr[i] += arr[i-1]
	}
	duration := time.Since(start)

	return uint32(start.Unix()), duration.Nanoseconds()
}

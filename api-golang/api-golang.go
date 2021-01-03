package main

import (
	"encoding/json"
	"flag"
	"github.com/cactus/go-statsd-client/statsd"
	"github.com/fasthttp/router"
	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttp/reuseport"
	"log"
	"time"
)

import _ "net/http/pprof"

type InputData struct {
	Arr  []int  `json:"arr"`
	Mode string `json:"mode"`
}

type OutputData struct {
	Message       string `json:"message"`
	Data          []int  `json:"data"`
	ExecutionTime int64  `json:"exec_ns"`
}
type ClientStats struct {
	NumRequests   int64
	Duration      int64
	DurationCount int64
	DurationSum   int64
}

var (
	addr = flag.String("addr", ":8080", "TCP address to listen to")
	// statsdClient, statsdClientErr = statsd.New(statsd.Address("127.0.0.1:9125"), statsd.Prefix("api"))
	statsdConfig = &statsd.ClientConfig{
		Address:       "127.0.0.1:9125",
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
		log.Fatal(statsdClientErr)
	}

	defer statsdClient.Close()

	requestRouter := router.New()
	requestRouter.POST("/", updateStatsDMetrics(processCumSumRequest))

	listener, httpListenErr := reuseport.Listen("tcp4", *addr)

	if httpListenErr != nil {
		log.Fatalf("Error in reuseport listener: %s", httpListenErr)
	}

	defer listener.Close()

	if err := fasthttp.Serve(listener, requestRouter.Handler); err != nil {
		log.Fatalf("Error in Serve: %s", err)
	}
}

func updateStatsDMetrics(h fasthttp.RequestHandler) fasthttp.RequestHandler {
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
			Data:          make([]int, 0),
			ExecutionTime: 0,
		}
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
	} else {
		cumSumArr := postBody.Arr
		executionTime := cumSum(cumSumArr)

		output = OutputData{
			Message:       "success",
			Data:          cumSumArr,
			ExecutionTime: executionTime,
		}
	}

	marshalledJSON, _ := json.Marshal(output)
	ctx.SetContentType("application/json")
	ctx.SetBody(marshalledJSON)
}

func cumSum(arr []int) int64 {
	start := time.Now()
	for i := 1; i < len(arr); i++ {
		arr[i] += arr[i-1]
	}
	duration := time.Since(start)

	return duration.Nanoseconds()
}

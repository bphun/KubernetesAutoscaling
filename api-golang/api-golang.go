package main

import (
	"encoding/json"
	"flag"
	"github.com/cactus/go-statsd-client/statsd"
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
	addr         = flag.String("addr", ":8080", "TCP address to listen to")
	statsdConfig = &statsd.ClientConfig{
		Address:       "127.0.0.1:9125",
		Prefix:        "api",
		ResInterval:   time.Minute,
		UseBuffered:   true,
		FlushInterval: 300 * time.Millisecond,
	}
	statsdClient, err = statsd.NewClientWithConfig(statsdConfig)
	clientStats       ClientStats
)

func main() {
	flag.Parse()

	h := requestHandler

	if err != nil {
		log.Fatal(err)
	}

	defer statsdClient.Close()

	listener, err := reuseport.Listen("tcp4", *addr)

	if err != nil {
		log.Fatalf("Error in reuseport listener: %s", err)
	}

	defer listener.Close()

	if err := fasthttp.Serve(listener, h); err != nil {
		log.Fatalf("Error in Serve: %s", err)
	}
}

func requestHandler(ctx *fasthttp.RequestCtx) {
	var requestDuration int64

	if string(ctx.Method()) == "POST" && string(ctx.RequestURI()) == "/" {
		start := time.Now()
		processCumSumRequest(ctx)
		requestDuration = time.Since(start).Microseconds()
	} else {
		ctx.SetStatusCode(fasthttp.StatusNotFound)
	}

	clientStats.NumRequests++
	clientStats.Duration += requestDuration

	updateStatsd(clientStats.NumRequests, clientStats.Duration)
}

func updateStatsd(numRequests int64, requestDuration int64) {
	err := statsdClient.Inc("requests.count", numRequests, 1.0)
	if err != nil {
		log.Printf("%s\n", err.Error())
	}
	err = statsdClient.Timing("requests.duration", requestDuration, 1.0)
	if err != nil {
		log.Printf("%s\n", err.Error())
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
		// var executionTime int64 = 0
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

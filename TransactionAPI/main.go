package main

import (
	"context"
	"flag"
	"log"
	"net"
	"net/http"
	"sync"

	pb "github.com/bphun/KubernetesAutoscaling/TransactionAPI/TransactionAPI"
	"github.com/grpc-ecosystem/go-grpc-prometheus"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedTransactionAPIServer
}

type Transaction struct {
	InArr     []int32 `bson:"in_arr"`
	OutArr    []int32 `bson:"out_arr"`
	ExecTime  int64   `bson:"exec_time"`
	StartTime uint32  `bson:"start_time"`
}

type ObjectID struct {
	Id string
}

const (
	GRPC_PORT                     = "0.0.0.0:5001"
	PROMETHEUS_PORT               = "0.0.0.0:9090"
	DEFAULT_MDB_CONNECTION_STRING = "mongodb://localhost:27017"
	DB                            = "db_transaction_manager"
	TRANSACTIONS                  = "col_transaction"
)

var (
	mdbConnectionString   = flag.String("mdbAddr", DEFAULT_MDB_CONNECTION_STRING, "MongoDB server address")
	dbClientInstance      *mongo.Client
	dbClientInstanceError error
	dbOnce                sync.Once

	// Create a metrics registry.
	reg = prometheus.NewRegistry()

	// Create some standard server metrics.
	grpcMetrics = grpc_prometheus.NewServerMetrics()

	// Create a customized counter metric.
	customizedCounterMetric = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "demo_server_say_hello_method_handle_count",
		Help: "Total number of RPCs handled on the server.",
	}, []string{"name"})
)

func (s *server) SaveTransaction(ctx context.Context, in *pb.TransactionRequest) (*pb.TransactionReply, error) {
	resultMessage := "Saved transaction"

	r, err := createTransaction(Transaction{InArr: in.GetInArr(), OutArr: in.GetOutArr(), ExecTime: in.GetExecTime(), StartTime: in.GetStartTime()})

	if err != nil {
		resultMessage = err.Error()
	} else {
		if oid, ok := r.InsertedID.(primitive.ObjectID); ok {
			resultMessage += " " + oid.String()
			log.Printf("Created transaction: %s", oid.String())
		}
	}

	return &pb.TransactionReply{Message: resultMessage}, nil
}

func getMongoClient() (*mongo.Client, error) {
	dbOnce.Do(func() {
		clientOptions := options.Client().ApplyURI(*mdbConnectionString)
		clientOptions = clientOptions.SetMinPoolSize(2)
		clientOptions = clientOptions.SetMaxPoolSize(20)

		log.Printf("Connecting to MongoDB at %s", *mdbConnectionString)

		client, err := mongo.Connect(context.Background(), clientOptions)
		if err != nil {
			dbClientInstanceError = err
		}

		log.Printf("Connected to MongoDB at %s", *mdbConnectionString)
		log.Printf("Pinging MongoDB cluster to test connection")

		err = client.Ping(context.Background(), nil)
		if err != nil {
			log.Fatalf("Failed to ping MongoDB cluster")
			dbClientInstanceError = err
		} else {
			log.Printf("Pinged MongoDB cluster")
		}

		dbClientInstance = client
	})
	return dbClientInstance, dbClientInstanceError
}

func createTransaction(task Transaction) (*mongo.InsertOneResult, error) {
	client, err := getMongoClient()
	if err != nil {
		return nil, err
	}
	collection := client.Database(DB).Collection(TRANSACTIONS)

	r, err := collection.InsertOne(context.Background(), task)
	if err != nil {
		return nil, err
	}
	return r, nil
}

func init() {
	// Register standard server metrics and customized metrics to registry.
	reg.MustRegister(grpcMetrics, customizedCounterMetric)
	customizedCounterMetric.WithLabelValues("Test")
}

func main() {
	flag.Parse()

	lis, err := net.Listen("tcp", GRPC_PORT)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	defer lis.Close()

	httpServer := &http.Server{Handler: promhttp.HandlerFor(reg, promhttp.HandlerOpts{}), Addr: PROMETHEUS_PORT}
	log.Printf("gRPC server listening on port %s", GRPC_PORT)

	grpc_prometheus.EnableHandlingTimeHistogram()
	grpcServer := grpc.NewServer(
		grpc.StreamInterceptor(grpcMetrics.StreamServerInterceptor()),
		grpc.UnaryInterceptor(grpcMetrics.UnaryServerInterceptor()),
	)
	pb.RegisterTransactionAPIServer(grpcServer, &server{})
	grpcMetrics.InitializeMetrics(grpcServer)

	go func() {
		log.Printf("Started prometheus server at %s", PROMETHEUS_PORT)
		if err := httpServer.ListenAndServe(); err != nil {
			log.Fatal("Unable to start a http server.")
		}
	}()

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}

package main

import (
	pb "bphun/k8sAutoscaling/TransactionAPI/TransactionAPI"
	"context"
	"log"
	"net"
	"sync"

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

const (
	GRPC_PORT         = ":50051"
	CONNECTION_STRING = "mongodb://localhost:27017"
	DB                = "db_transaction_manager"
	TRANSACTIONS      = "col_transaction"
)

var (
	dbClientInstance      *mongo.Client
	dbClientInstanceError error
	dbOnce                sync.Once
)

func (s *server) SaveTransaction(ctx context.Context, in *pb.TransactionRequest) (*pb.TransactionReply, error) {
	errorMessage := "Saved transaction"

	err := createTransaction(Transaction{InArr: in.GetInArr(), OutArr: in.GetOutArr(), ExecTime: in.GetExecTime(), StartTime: in.GetStartTime()})

	if err != nil {
		errorMessage = err.Error()
	}

	return &pb.TransactionReply{Message: errorMessage}, nil
}

func getMongoClient() (*mongo.Client, error) {
	dbOnce.Do(func() {
		clientOptions := options.Client().ApplyURI(CONNECTION_STRING)
		client, err := mongo.Connect(context.TODO(), clientOptions)
		if err != nil {
			dbClientInstanceError = err
		}

		err = client.Ping(context.TODO(), nil)
		if err != nil {
			dbClientInstanceError = err
		}
		dbClientInstance = client
	})
	return dbClientInstance, dbClientInstanceError
}

func createTransaction(task Transaction) error {
	client, err := getMongoClient()
	if err != nil {
		return err
	}
	collection := client.Database(DB).Collection(TRANSACTIONS)

	_, err = collection.InsertOne(context.Background(), task)
	if err != nil {
		return err
	}
	return nil
}

func main() {
	lis, err := net.Listen("tcp", GRPC_PORT)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	log.Printf("Listening on port %s", GRPC_PORT)

	s := grpc.NewServer()
	pb.RegisterTransactionAPIServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}

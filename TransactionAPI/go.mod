module github.com/bphun/k8sAutoscaling/TransactionAPI

go 1.16

require (
	go.mongodb.org/mongo-driver v1.5.0
	google.golang.org/grpc v1.36.0
	google.golang.org/protobuf v1.26.0
)

// replace bphun/k8sAutoscaling/TransactionAPI => ../TransactionAPI/TransactionAPI/

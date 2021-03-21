module bphun/k8sAutoscaling/TransactionAPI

go 1.16

require (
	google.golang.org/grpc v1.35.1
	google.golang.org/protobuf v1.25.0
  go.mongodb.org/mongo-driver v1.5.0
)

replace bphun/k8sAutoscaling/TransactionAPI => ../TransactionAPI/TransactionAPI/

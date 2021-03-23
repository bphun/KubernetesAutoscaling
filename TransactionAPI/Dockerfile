FROM golang:1.16.2

RUN adduser transaction-api-user

WORKDIR /home/transaction-api

COPY main.go main.go
COPY go.mod go.mod
COPY go.sum go.sum
COPY TransactionAPI TransactionAPI

RUN go mod vendor && \
    go build main.go && \
    chown -R transaction-api-user:transaction-api-user ./ 

USER transaction-api-user

# ENV GOMAXPROCS=1

EXPOSE 5001
ENTRYPOINT ["./main", "-mdbAddr=mongodb://transaction-db.transaction-db:27017"]
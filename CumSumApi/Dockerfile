FROM golang:1.16.2

RUN adduser cumsum-api-user

WORKDIR /home/cumsum-api

COPY main.go main.go
COPY go.mod go.mod
COPY go.sum go.sum

RUN go mod vendor && \
    go build main.go && \
    chown -R cumsum-api-user:cumsum-api-user ./ 

USER cumsum-api-user

# ENV GOMAXPROCS=1

EXPOSE 5000
ENTRYPOINT ["./main", "-httpAddr=0.0.0.0:5000", "-grpcAddr=transaction-api.transaction-api:5001"]
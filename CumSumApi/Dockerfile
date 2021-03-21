FROM golang:1.16.2

RUN adduser api-user

WORKDIR /home/api

RUN mkdir ../TransactionAPI

COPY main.go main.go
COPY go.mod go.mod
COPY ../TransactionAPI ../TransactionAPI

# RUN go get -u github.com/fasthttp/router && \
#     go get -u github.com/valyala/fasthttp && \
#     go get -u github.com/valyala/fasthttp/reuseport && \
#     go get -u github.com/cactus/go-statsd-client/statsd && \
#     go build && \
#     chown -R api-user:api-user ./ 
    
RUN go mod vendor \
    go build && \
    chown -R api-user:api-user ./ 

USER api-user

# ENV GOMAXPROCS=1

EXPOSE 5000
ENTRYPOINT ["./main", "-addr=0.0.0.0:5000"]
FROM golang:1.16.2

RUN adduser transaction-api-user

WORKDIR /home/api

COPY main.go main.go

RUN go get -u github.com/fasthttp/router && \
    go get -u github.com/valyala/fasthttp && \
    go get -u github.com/valyala/fasthttp/reuseport && \
    go get -u github.com/cactus/go-statsd-client/statsd && \
    go build && \
    chown -R transaction-api-user:transaction-api-user ./ 
    
USER transaction-api-user

# ENV GOMAXPROCS=1

EXPOSE 5000
ENTRYPOINT ["./main"]
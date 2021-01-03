FROM golang:1.15.6

RUN adduser api-user

WORKDIR /home/api

COPY api-golang.go api-golang.go

RUN go get -u github.com/fasthttp/router && \
    go get -u github.com/valyala/fasthttp && \
    go get -u github.com/valyala/fasthttp/reuseport && \
    go get -u github.com/cactus/go-statsd-client/statsd && \
    go build && \
    chown -R api-user:api-user ./ 
    
USER api-user

# ENV GOMAXPROCS=1

EXPOSE 5000
ENTRYPOINT ["./api", "-addr=0.0.0.0:5000"]
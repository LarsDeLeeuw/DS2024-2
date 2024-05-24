FROM golang:1.22.3 as builder

WORKDIR /go-microservice/

COPY go.mod .
COPY go.sum .
COPY cmd/gateway/v1/* cmd/gateway/
COPY pkg/gateway/v1/* pkg/gateway/v1/

RUN CGO_ENABLED=0 go build -o gateway /go-microservice/cmd/gateway/main.go

FROM alpine:latest

WORKDIR /go-microservice

COPY --from=builder /go-microservice/ /go-microservice/

CMD ./gateway
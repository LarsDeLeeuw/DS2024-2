FROM golang:1.22.3 as builder

WORKDIR /go-microservice/

COPY go.mod .
COPY go.sum .
COPY cmd/auth/v1/* cmd/auth/
COPY pkg/auth/v1/* pkg/auth/v1/

RUN CGO_ENABLED=0 go build -o auth /go-microservice/cmd/auth/main.go

FROM alpine:latest

WORKDIR /go-microservice

COPY --from=builder /go-microservice/ /go-microservice/

EXPOSE 3000

CMD ./auth
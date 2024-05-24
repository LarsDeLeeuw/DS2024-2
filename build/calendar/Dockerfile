FROM golang:1.22.3 as builder

WORKDIR /go-microservice/

COPY go.mod .
COPY go.sum .
COPY cmd/calendar/v1/* cmd/calendar/
COPY pkg/calendar/v1/* pkg/calendar/v1/

RUN CGO_ENABLED=0 go build -o calendar /go-microservice/cmd/calendar/main.go

FROM alpine:latest

WORKDIR /go-microservice

COPY --from=builder /go-microservice/ /go-microservice/

CMD ./calendar
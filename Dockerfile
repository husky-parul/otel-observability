FROM golang:alpine

WORKDIR /go/src/github.com/otel-go

COPY . .
RUN go install /go/src/github.com/otel-go/src/main.go

CMD ["/go/bin/main"]
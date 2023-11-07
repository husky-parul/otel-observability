FROM golang:alpine

WORKDIR /apps

COPY . .

# Build the Go application
RUN go get -d -v ./...
RUN CGO_ENABLED=0 GOOS=linux go build -o /otelobs


CMD ["/otelobs"]
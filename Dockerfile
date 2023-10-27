FROM golang:1.21-alpine

WORKDIR /usr/local/src

# dependencies
COPY ["go.mod", "go.sum", "./"]
RUN go mod download

# build
COPY ./ ./
RUN go build -o name-details-adder ./cmd/app/main.go
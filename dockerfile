# Builder image
FROM golang:alpine AS builder


WORKDIR /app

# Install dependencies
COPY go.mod go.sum ./
RUN go mod download

COPY *.go .

# Build
RUN go build -o /godork 

# Runtime image
FROM alpine AS app

COPY --from=builder /godork /godork 

ENTRYPOINT ["/godork"]
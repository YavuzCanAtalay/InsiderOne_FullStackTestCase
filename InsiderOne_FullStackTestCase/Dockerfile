# Build stage
FROM golang:1.22-alpine AS builder

WORKDIR /app

ENV GOPROXY=direct
ENV GONOSUMDB=*
ENV GOFLAGS=-mod=mod

RUN apk add --no-cache git

# Download dependencies first (cached layer)
COPY go.mod ./
RUN go mod download && go mod verify

# Copy source and build
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /football-api ./cmd/api

# Final stage — minimal image
FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /root/
COPY --from=builder /football-api .

EXPOSE 8080
CMD ["./football-api"]

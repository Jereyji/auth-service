FROM golang:1.22.3-alpine AS builder

WORKDIR /build

ADD go.mod go.sum ./

RUN go mod tidy

COPY . .

RUN go build -o main cmd/sender/main.go


FROM alpine

WORKDIR /app

COPY deployments/.env /app/deployments/.env

COPY config/sender_config.yaml /app/config/sender_config.yaml

COPY --from=builder /build/main /app/main
FROM golang:1.22.3-alpine AS builder

WORKDIR /build

ADD go.mod go.sum ./

RUN go mod tidy

COPY . .

RUN go build -o main cmd/auth/main.go


FROM alpine

WORKDIR /app

COPY deployments/.env /app/deployments/.env

COPY config/auth_config.yaml /app/config/auth_config.yaml

COPY --from=builder /build/main /app/main

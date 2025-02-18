FROM golang:1.22.3-alpine AS builder

WORKDIR /build

ADD go.mod go.sum ./

RUN go mod tidy

COPY . .

RUN go build -o main cmd/main.go


FROM alpine

WORKDIR /app

COPY --from=builder /build/main /app/main

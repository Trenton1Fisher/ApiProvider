FROM golang:1.22.1-alpine AS builder

WORKDIR /app

COPY api/go.mod api/go.sum ./
RUN go mod download 

COPY . .

RUN go build -o api/api api/server.go api/dbHandler.go api/redis.go

RUN ls -al api/api

EXPOSE 9010
CMD ["./api/api"]

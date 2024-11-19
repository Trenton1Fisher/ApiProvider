FROM golang:1.22.1-alpine AS builder

WORKDIR /app

COPY api/go.mod ./
RUN go mod tidy

COPY . .

RUN go build -o api/api api/server.go

RUN ls -al api/api

EXPOSE 9010
CMD ["./api/api"]

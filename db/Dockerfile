FROM golang:1.22.1-alpine AS builder

WORKDIR /app

COPY go.mod ./
RUN go mod tidy

COPY . .

RUN go build -o db/db db/seed.go

RUN ls -al /app/db

EXPOSE 5173
CMD ["./db/db"]

FROM golang:1.23 AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o crawler .

FROM alpine:latest

RUN apk --no-cache add ca-certificates

COPY --from=builder /app/crawler /usr/local/bin/crawler

ENTRYPOINT ["crawler"]
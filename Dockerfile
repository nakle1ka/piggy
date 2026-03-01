FROM golang:1.25-alpine AS builder
WORKDIR /app
COPY . .

RUN apk add --no-cache git
RUN go mod download
RUN go build -o ./bin/app ./cmd/app/main.go

FROM alpine:latest
WORKDIR /app

COPY --from=builder /app/bin/app .
COPY .env .

CMD ["./app"]
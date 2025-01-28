FROM golang:1.23-alpine AS builder

WORKDIR /app
COPY . .

RUN go mod download
RUN go build -o main .

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/main .
COPY --from=builder /app/docs ./docs
COPY --from=builder /app/migrations ./migrations
COPY .env .

CMD ["./main"]
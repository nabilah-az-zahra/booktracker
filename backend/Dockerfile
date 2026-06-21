FROM golang:1.26-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o main ./cmd/main.go

FROM alpine:latest

RUN addgroup -S appgroup && adduser -S appuser -G appgroup

RUN apk add --no-cache tzdata

WORKDIR /app

COPY --from=builder /app/main .

USER appuser

EXPOSE 8080
CMD ["./main"]
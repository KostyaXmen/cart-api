FROM golang:1.26-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o cart-api ./cmd/main.go

FROM alpine:3.19

WORKDIR /app

COPY --from=builder /app/cart-api .

COPY config/ config/
COPY db/ /app/db/

EXPOSE 8080

CMD ["./cart-api"]
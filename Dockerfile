FROM golang:1.23.1-alpine AS builder

RUN apk add --no-cache git

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -trimpath -ldflags="-s -w" -o notify ./main.go

FROM alpine:3.18
WORKDIR /app

COPY --from=builder /app/notify .

EXPOSE 8080

CMD ["./notify"]
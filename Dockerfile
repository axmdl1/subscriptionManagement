# ---------- BUILD STAGE ----------
FROM golang:1.24-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o app ./cmd/app


# ---------- RUNTIME STAGE ----------
FROM alpine:3.20

WORKDIR /app

RUN apk add --no-cache netcat-openbsd

COPY --from=builder /app/app .

COPY --from=builder /app/migrations ./migrations
COPY --from=builder /app/docs ./docs

COPY wait-for-postgres.sh .
RUN chmod +x wait-for-postgres.sh

EXPOSE 8080

CMD ["./wait-for-postgres.sh", "./app"]

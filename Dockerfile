FROM golang:1.23 AS builder

WORKDIR /app

COPY . .

ENV GOOS=linux
ENV GOARCH=amd64

RUN go build -o bin/server ./cmd/server

RUN [ -f config.yaml ] || touch config.yaml

FROM alpine:latest

RUN apk add --no-cache curl \
    && curl -L https://github.com/golang-migrate/migrate/releases/download/v4.17.1/migrate.linux-amd64.tar.gz | tar xvz \
    && mv migrate /usr/local/bin/migrate \
    && apk del curl

RUN addgroup -S appgroup && adduser -S appuser -G appgroup

WORKDIR /app

COPY --from=builder /app/bin/server /app/
COPY --from=builder /app/config.yaml /app/
COPY --from=builder /app/internal/infrastructure/database/mysql/migrations /app/migrations

RUN chmod +x /app/server

USER appuser

EXPOSE 8080

CMD migrate -path /app/migrations -database "mysql://$DB_USER:$DB_PASSWORD@tcp($DB_HOST:$DB_PORT)/$DB_NAME" up && /app/server

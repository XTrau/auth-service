FROM golang:1.27-alpine AS builder

WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build -ldflags "-s -w" -o ./bin/app -a ./cmd/api
RUN CGO_ENABLED=0 go build -ldflags "-s -w" -o ./bin/migrate -a ./cmd/migrate

FROM alpine:3.20

WORKDIR /app

COPY --from=builder /app/bin/app .
COPY --from=builder /app/bin/migrate .
COPY --from=builder /app/migrations ./migrations

EXPOSE 8080

ENTRYPOINT ["./app"]

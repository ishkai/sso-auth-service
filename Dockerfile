FROM golang:1.26-alpine AS builder

WORKDIR /src

RUN apk add --no-cache git gcc musl-dev

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=1 GOOS=linux go build -o /out/sso ./cmd/sso
RUN CGO_ENABLED=1 GOOS=linux go build -o /out/migator ./cmd/migator

FROM alpine:3.22

WORKDIR /app

RUN apk add --no-cache ca-certificates tzdata

COPY --from=builder /out/sso /usr/local/bin/sso
COPY --from=builder /out/migator /usr/local/bin/migator

COPY config ./config
COPY migrations ./migrations

EXPOSE 44044
EXPOSE 8080

CMD ["sso", "--config=/app/config/docker.yaml"]
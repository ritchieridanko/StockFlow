# ---------- BUILD STAGE ----------
FROM golang:1.26.0-alpine3.22 AS builder

RUN apk add --no-cache git

WORKDIR /app

COPY backend/go.mod backend/go.sum ./
RUN go mod download

COPY backend/cmd/app ./cmd/app
COPY backend/configs ./configs
COPY backend/internal ./internal

RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o bin/app cmd/app/main.go

# ---------- RUNTIME STAGE ----------
FROM alpine:3.22

RUN apk add --no-cache ca-certificates

WORKDIR /root

COPY --from=builder /app/bin ./bin
COPY --from=builder /app/configs ./configs

EXPOSE 8080

ENTRYPOINT ["./bin/app"]

FROM golang:1.23-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git ca-certificates tzdata

COPY backend/go.mod backend/go.sum ./
RUN go mod download

COPY backend/ .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o server ./cmd/server/main.go

FROM alpine:3.20

WORKDIR /app

RUN apk add --no-cache ca-certificates tzdata curl

COPY --from=builder /app/server /app/server

ENV PORT=8080
EXPOSE 8080

ENTRYPOINT ["/app/server"]

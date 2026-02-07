FROM golang:1.23-alpine AS builder
WORKDIR /app
COPY go.mod ./
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /telegram-pr-notify .

FROM gcr.io/distroless/static-debian12:nonroot
COPY --from=builder /telegram-pr-notify /telegram-pr-notify
ENTRYPOINT ["/telegram-pr-notify"]

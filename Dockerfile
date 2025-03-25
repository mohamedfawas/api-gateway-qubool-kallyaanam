FROM golang:1.23.4-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o api-gateway ./cmd/server

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/api-gateway .
# Copy .env file for development/testing containers
# (Comment this out for production builds)
COPY .env.docker ./.env
EXPOSE 8080
CMD ["./api-gateway"]
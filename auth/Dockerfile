# build stage
FROM golang:1.21-alpine AS builder
RUN apk add --no-cache git
WORKDIR /build
COPY src/authdb ./authdb
COPY src/main ./main
WORKDIR /build/main
RUN go mod download && go build -o /app ./...

# final stage
FROM alpine:3.19
RUN apk --no-cache add ca-certificates
COPY --from=builder /app /app
EXPOSE 8080
ENTRYPOINT ["/app"]

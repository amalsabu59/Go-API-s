# Build stage
FROM golang:1.22.4-alpine as build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o main ./cmd/main.go

# Final stage
FROM alpine:latest

WORKDIR /app
RUN apk --no-cache add ca-certificates

COPY --from=build /app/main /app/main

CMD ["/app/main"]

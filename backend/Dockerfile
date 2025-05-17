# Development-optimized backend setup
FROM golang:1.21-bullseye AS build
RUN apt-get update && apt-get install -y gcc libc6-dev
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=1 GOOS=linux go build -o payd ./cmd/main.go

# Lean production image
FROM debian:bullseye-slim
RUN apt-get update && apt-get install -y ca-certificates
COPY --from=build /app/payd /app/
WORKDIR /app
EXPOSE 8080
CMD ["./payd"]
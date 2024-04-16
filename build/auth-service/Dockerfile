FROM golang:1.22 AS build
WORKDIR /app

# Install required modules
COPY go.mod go.sum ./
RUN go mod download

# Build the service
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /bin/auth-service ./cmd/auth-service

# Runtime Environment
FROM gcr.io/distroless/static-debian12
COPY --from=build /bin/auth-service .
EXPOSE 90 9090
CMD ["./auth-service"]
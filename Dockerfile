FROM golang:1.24-alpine

WORKDIR /app

# Copy go modules files
COPY go.mod ./
COPY go.sum ./

RUN go mod download

# Copy all source code
COPY . .

# Install air for hot reload (use v1.52.3 compatible with Go 1.24)
RUN go install github.com/air-verse/air@v1.52.3

EXPOSE 8080

# Default command runs with air for hot reload
CMD ["air", "-c", ".air.toml"]

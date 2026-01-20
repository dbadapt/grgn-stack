FROM golang:1.22-alpine

WORKDIR /app

# Copy go modules files from backend
COPY backend/go.mod ./
COPY backend/go.sum ./

# Copy pkg module (required by go.mod replace directive)
COPY pkg ../pkg

RUN go mod download

# Copy backend source
COPY backend/ .

# Install air for hot reload (use v1.52.3 compatible with Go 1.24)
RUN go install github.com/air-verse/air@v1.52.3

EXPOSE 8080

CMD ["air", "-c", ".air.toml"]

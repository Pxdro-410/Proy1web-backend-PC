# Build stage
FROM golang:1.26-alpine AS builder

# Dependencias del sistema necesarias para compilar con cgo (lib/pq)
RUN apk add --no-cache gcc musl-dev

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o server .

# Runtime stage 
FROM alpine:latest

# Certificados TLS necesarios para conexiones HTTPS (ej. Railway externo)
RUN apk add --no-cache ca-certificates

WORKDIR /app

# Copiar solo el binario compilado
COPY --from=builder /app/server .

EXPOSE 8080

CMD ["./server"]

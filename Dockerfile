# Build stage
FROM golang:1.23-alpine AS builder
WORKDIR /build
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o h2pcontrol-manager .

# Runtime stage — needs Python for proto compilation via GetStub RPC
FROM python:3.12-slim
WORKDIR /app

RUN pip install --no-cache-dir grpcio-tools betterproto2_compiler==0.4.0

COPY --from=builder /build/h2pcontrol-manager .

EXPOSE 50051

ENTRYPOINT ["./h2pcontrol-manager"]
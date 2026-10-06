FROM golang:1.27.1-bookworm AS builder

RUN apt-get update && apt-get install -y --no-install-recommends \
    libx11-dev libgl1-mesa-dev libgles2-mesa-dev libegl1-mesa-dev \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /app
COPY go.mod go.sum ./
COPY . .
RUN CGO_ENABLED=1 GOOS=linux go build -o lx-client cmd/client/app/*
RUN CGO_ENABLED=1 GOOS=windows go build -o client.exe cmd/client/app/*

FROM scratch AS exporter
COPY --from=builder /app/lx-client .

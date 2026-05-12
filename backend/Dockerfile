## syntax=docker/dockerfile:1.7

# ---------- build stage ----------
FROM golang:1.26-alpine AS builder

WORKDIR /src

# Cache de dependências
COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download

# Código-fonte
COPY . .

# Build estático (sem CGO), com flags de redução de tamanho
ENV CGO_ENABLED=0 GOOS=linux GOFLAGS=-trimpath
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    go build -ldflags="-s -w" -o /out/smid10-server ./cmd/server

# ---------- runtime stage ----------
FROM gcr.io/distroless/static-debian12:nonroot

WORKDIR /app
COPY --from=builder /out/smid10-server /app/smid10-server

USER nonroot:nonroot
EXPOSE 8080

ENTRYPOINT ["/app/smid10-server"]

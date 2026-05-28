# syntax=docker/dockerfile:1

# ============================================
# Build arguments (override with --build-arg)
# ============================================
ARG GO_VERSION=1.25
ARG ALPINE_VERSION=3.23

# ============================================
# Stage 1: Build
# ============================================
FROM golang:${GO_VERSION}-alpine${ALPINE_VERSION} AS builder

RUN apk add --no-cache git ca-certificates tzdata

WORKDIR /src

# Dependency resolution first — maximizes layer cache reuse
COPY go.mod go.sum ./
RUN go mod download

# Copy source (excludes patterns in .dockerignore)
COPY . .

RUN CGO_ENABLED=0 go build \
    -ldflags="-s -w" \
    -trimpath \
    -o /app/deeplx-translategemma \
    .

# ============================================
# Stage 2: Runtime (minimal image)
# ============================================
FROM alpine:${ALPINE_VERSION}

# Metadata
LABEL org.opencontainers.image.title="DeepLX-TranslateGemma"
LABEL org.opencontainers.image.description="DeepLX-compatible translation API backed by OpenAI-compatible LLM"
LABEL org.opencontainers.image.source="https://github.com/mustang5910/deeplx-translategemma"

# Runtime dependencies: ca-certificates for TLS, tzdata for timezone
RUN apk add --no-cache ca-certificates tzdata

ENV TZ=Asia/Shanghai

# Create unprivileged user
RUN adduser -D -H -h /app appuser

WORKDIR /app

# Copy artifacts from builder
COPY --from=builder /app/deeplx-translategemma .
COPY docker-entrypoint.sh .
COPY etc/deeplx-api.yaml.example ./deeplx-api.yaml.example

RUN chmod +x docker-entrypoint.sh

EXPOSE 8888

# Health check — probes the translate endpoint to confirm the service is alive
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget -q -O /dev/null http://127.0.0.1:8888/translate || exit 1

USER appuser

ENTRYPOINT ["./docker-entrypoint.sh"]
CMD ["./deeplx-translategemma", "-f", "/app/deeplx-api.yaml"]

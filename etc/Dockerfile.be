# ── Production Backend ─────────────────────────────────────────────
# Multi-stage build: compile Go binaries, run FROM alpine (~20 MB)
# Serves both the API (cmd/server) and job worker (cmd/worker) images.
# ──────────────────────────────────────────────────────────────────

FROM golang:1.25-alpine AS builder

RUN apk add --no-cache git gcc musl-dev

WORKDIR /src

COPY be/go.mod be/go.sum ./
RUN go mod download

COPY be/ .

# -ldflags="-s -w"     strip debug info -> smaller binary
# -tags timetzdata     embed timezone DB into binary (no need OS tzdata)
RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags="-s -w" \
    -tags timetzdata \
    -o /opt/server ./cmd/api

RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags="-s -w" \
    -tags timetzdata \
    -o /opt/worker ./cmd/worker

RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags="-s -w" \
    -o /opt/genkey ./cmd/genkey

RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags="-s -w" \
    -tags timetzdata \
    -o /opt/migrate ./cmd/migration

# ── Runtime ───────────────────────────────────────────────────────
FROM alpine:3.21

RUN apk add --no-cache ca-certificates tzdata

COPY --from=builder /opt/server  /server
COPY --from=builder /opt/worker  /worker
COPY --from=builder /opt/genkey  /genkey
COPY --from=builder /opt/migrate /migrate
COPY etc/entrypoint.sh           /entrypoint.sh

RUN chmod +x /entrypoint.sh

WORKDIR /app
EXPOSE 8080

ENTRYPOINT ["/entrypoint.sh"]
CMD ["/server"]

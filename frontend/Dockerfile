# ── Build stage ───────────────────────────────────────────────────────────────
# Build context is the repository root (radarcdmx/) because the frontend module
# uses a local `replace` directive pointing at ../backend/rcapi.
FROM golang:1.26-alpine AS builder

RUN apk add --no-cache git

WORKDIR /src

# The replaced module must be present before `go mod download` / build.
COPY backend/rcapi/ ./backend/rcapi/

# Frontend module: copy manifests first to cache the dependency layer.
COPY frontend/go.mod frontend/go.sum ./frontend/
WORKDIR /src/frontend
RUN go mod download

# Now the rest of the frontend source.
COPY frontend/ ./
RUN go build -o /radarcdmx-web .

# ── Runtime stage ─────────────────────────────────────────────────────────────
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /app

COPY --from=builder /radarcdmx-web ./
# Static UI served directly from disk via os.DirFS("RadarMX-main").
COPY frontend/RadarMX-main/ ./RadarMX-main/

EXPOSE 4242

CMD ["./radarcdmx-web"]

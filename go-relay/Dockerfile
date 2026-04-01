# Build stage
FROM golang:1.22-alpine AS builder

RUN apk add --no-cache git ca-certificates

WORKDIR /app

# Copy dependency files first for caching
COPY go.mod go.sum ./
RUN go mod download

# Copy source
COPY . .

# Build binary
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /relay ./cmd/server/

# Runtime stage
FROM alpine:3.19

# Install Chromium and required fonts/libs for headless browser
RUN apk add --no-cache \
    ca-certificates \
    chromium \
    nss \
    freetype \
    harfbuzz \
    ttf-freefont \
    font-noto-emoji

WORKDIR /app

# Copy binary
COPY --from=builder /relay /app/relay

# Copy static assets (if any)
COPY --from=builder /app/static/ /app/static/ 2>/dev/null || true

# Chromium path for headless sessions (auto-detected, but set explicitly for reliability)
ENV PUPPETEER_EXECUTABLE_PATH=/usr/bin/chromium-browser
ENV PORT=3010
ENV DB_TYPE=sqlite

# Data volume
VOLUME /app/data

EXPOSE 3010

CMD ["/app/relay"]

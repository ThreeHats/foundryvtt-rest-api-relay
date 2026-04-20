# Go build stage
FROM golang:1.26-alpine AS builder

RUN apk add --no-cache git ca-certificates

WORKDIR /app

# Copy dependency files first for caching
COPY go-relay/go.mod go-relay/go.sum ./
RUN go mod download

# Copy source
COPY go-relay/ .

# Build binary
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /relay ./cmd/server/

# Frontend build stage
FROM node:22-alpine AS frontend

# Build Astro frontend
WORKDIR /frontend
COPY frontend/ .
RUN rm -rf node_modules && npm install && npx astro build

# Build docs
ARG ALGOLIA_APP_ID
ARG ALGOLIA_SEARCH_API_KEY
ARG ALGOLIA_INDEX_NAME
WORKDIR /docs
COPY docs/ .
RUN rm -rf node_modules && npm install && npx docusaurus build

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

# Run as non-root user
RUN adduser -D -u 1001 appuser

WORKDIR /app

# Copy binary
COPY --from=builder /relay /app/relay

# Copy API spec files (openapi.json, asyncapi.json, api-docs.json)
COPY public/*.json /app/static/public/

# Copy Astro frontend build output (astro outputs to ../public-dist relative to frontend/)
COPY --from=frontend /public-dist/ /app/static/public-dist/

# Copy docs build output
COPY --from=frontend /docs/build/ /app/static/docs/build/

# Chromium path for headless sessions
ENV PUPPETEER_EXECUTABLE_PATH=/usr/bin/chromium-browser
ENV PORT=3010
ENV DB_TYPE=sqlite
ENV STATIC_DIR=/app/static
ENV DATA_DIR=/app/data
# Self-hosted default: 0 = unlimited. Override in your compose file or fly.toml
# if you want to enforce a monthly request quota.
ENV MONTHLY_REQUEST_LIMIT=0

# Data volume
VOLUME /app/data

EXPOSE 3010

# Ensure data volume is writable
RUN mkdir -p /app/data && chown -R appuser:appuser /app /app/data

USER appuser

CMD ["/app/relay"]

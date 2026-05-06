# syntax=docker/dockerfile:1

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

# Frontend build stage (Astro)
FROM node:22-alpine AS frontend

WORKDIR /frontend
# Install deps before copying content so this layer is cached on content-only changes
COPY frontend/package.json frontend/pnpm-lock.yaml ./
RUN --mount=type=cache,target=/root/.cache/node/corepack \
    --mount=type=cache,target=/root/.local/share/pnpm/store \
    corepack enable && pnpm install --frozen-lockfile
COPY frontend/ .
RUN pnpm run build

# Docs build stage (Docusaurus)
FROM node:22-alpine AS docs

ARG ALGOLIA_APP_ID
ARG ALGOLIA_SEARCH_API_KEY
ARG ALGOLIA_INDEX_NAME

WORKDIR /docs
# Install deps before copying content so this layer is cached on content-only changes
COPY docs/package.json docs/pnpm-lock.yaml ./
RUN --mount=type=cache,target=/root/.cache/node/corepack \
    --mount=type=cache,target=/root/.local/share/pnpm/store \
    corepack enable && pnpm install --frozen-lockfile
COPY docs/ .
RUN pnpm run build

# Runtime stage: Ubuntu 24.04 (glibc) for NVIDIA GPU support.
# Alpine's musl libc cannot load NVIDIA's glibc-compiled EGL driver
# (libEGL_nvidia.so), causing Chrome WebGL initialization to fail.
# Ubuntu + glibc resolves this — Chrome can load the NVIDIA EGL ICD
# via GLVND and provide hardware-accelerated WebGL to Foundry VTT.
FROM ubuntu:24.04

ENV DEBIAN_FRONTEND=noninteractive

# Install Node.js 22 (NodeSource) and Google Chrome.
# Google Chrome is used over Chromium for better GPU driver integration.
# gosu: lightweight privilege-drop (replaces Alpine's su-exec).
# libvulkan1 + mesa-vulkan-drivers: Vulkan ICD (llvmpipe) for headless WebGL fallback.
# xvfb: virtual X11 framebuffer required by VirtualGL and the xvfb render mode.
# libgl1-mesa-dri + libegl-mesa0: Mesa OpenGL/EGL for non-NVIDIA GL rendering via Xvfb.
RUN apt-get update && apt-get install -y --no-install-recommends \
    ca-certificates curl gnupg \
    && curl -fsSL https://deb.nodesource.com/setup_22.x | bash - \
    && curl -fsSL https://dl.google.com/linux/linux_signing_key.pub \
       | gpg --dearmor -o /usr/share/keyrings/google-chrome-keyring.gpg \
    && echo "deb [arch=amd64 signed-by=/usr/share/keyrings/google-chrome-keyring.gpg] http://dl.google.com/linux/chrome/deb/ stable main" \
       > /etc/apt/sources.list.d/google-chrome.list \
    && apt-get update && apt-get install -y --no-install-recommends \
    nodejs \
    google-chrome-stable \
    fonts-freefont-ttf \
    fonts-noto-color-emoji \
    libvulkan1 \
    mesa-vulkan-drivers \
    xvfb \
    libgl1-mesa-dri \
    libegl-mesa0 \
    libegl1 \
    libglu1-mesa \
    gosu \
    && rm -rf /var/lib/apt/lists/* \
    && corepack enable && corepack prepare pnpm@latest --activate

# Create non-root user (entrypoint drops to this user via gosu)
RUN useradd -m -u 1001 appuser

WORKDIR /app

# Copy binary
COPY --from=builder /relay /app/relay

# Copy API spec files (openapi.json, asyncapi.json, api-docs.json)
COPY public/*.json /app/static/public/

# Copy Astro frontend build output (astro outputs to ../public-dist relative to frontend/)
COPY --from=frontend /public-dist/ /app/static/public-dist/

# Copy docs build output
COPY --from=docs /docs/build/ /app/static/docs/build/

# Chrome path for headless sessions (google-chrome-stable is found by findChromeBinary() too)
ENV PUPPETEER_EXECUTABLE_PATH=/usr/bin/google-chrome-stable
ENV PORT=3010
ENV DB_TYPE=sqlite
ENV STATIC_DIR=/app/static
ENV DATA_DIR=/app/data
# Self-hosted default: 0 = unlimited. Override in your compose file or fly.toml
# if you want to enforce a monthly request quota.
ENV MONTHLY_REQUEST_LIMIT=0
# NVIDIA container runtime: expose all GPUs and request graphics capabilities.
# NVIDIA_DRIVER_CAPABILITIES=graphics ensures libGLX_nvidia.so.0 (the NVIDIA
# Vulkan ICD) is injected with vk_icdGetInstanceProcAddr exported, which is
# required for ANGLE Vulkan hardware rendering via --use-vulkan=native.
# These ENV values are read by nvidia-container-toolkit at container start and
# have no effect when the NVIDIA runtime is not active.
ENV NVIDIA_VISIBLE_DEVICES=all
ENV NVIDIA_DRIVER_CAPABILITIES=graphics,utility

# Data volume
VOLUME /app/data

EXPOSE 3010

# Ensure data volume is writable
RUN mkdir -p /app/data && chown -R appuser:appuser /app /app/data

# Entrypoint resolves GPU device group GIDs at runtime and drops to appuser.
COPY go-relay/docker-entrypoint.sh /docker-entrypoint.sh
RUN chmod +x /docker-entrypoint.sh

ENTRYPOINT ["/docker-entrypoint.sh"]
CMD ["/app/relay"]

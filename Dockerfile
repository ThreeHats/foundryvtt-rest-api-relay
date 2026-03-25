FROM node:20-slim

WORKDIR /app

# Install required build dependencies for bcrypt, PostgreSQL, SQLite3, and Puppeteer (Chromium)
RUN apt-get update && apt-get install -y \
    python3 \
    make \
    g++ \
    build-essential \
    postgresql-client \
    sqlite3 \
    libsqlite3-dev \
    pkg-config \
    wget \
    ca-certificates \
    fonts-liberation \
    libappindicator3-1 \
    libasound2 \
    libatk-bridge2.0-0 \
    libatk1.0-0 \
    libcups2 \
    libdbus-1-3 \
    libdrm2 \
    libgbm1 \
    libnspr4 \
    libnss3 \
    libxcomposite1 \
    libxdamage1 \
    libxrandr2 \
    xdg-utils \
    git \
    # Add chromium-browser package
    chromium \
    --no-install-recommends && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/*

# Copy package.json (and lock files if they exist)
COPY package.json ./ 
COPY pnpm-lock.yaml* ./ 
COPY yarn.lock* ./ 
COPY package-lock.json* ./ 

# Install dependencies with --no-optional to avoid non-essential optional deps
RUN npm install -g pnpm && \
    if [ -f pnpm-lock.yaml ]; then \
      pnpm install --frozen-lockfile; \
    elif [ -f yarn.lock ]; then \
      yarn install --frozen-lockfile; \
    else \
      npm install; \
    fi

# Add node-fetch explicitly since it's needed for the forwarder
RUN pnpm add node-fetch @types/node-fetch

# Rebuild native modules for the current platform
RUN pnpm rebuild bcrypt
RUN pnpm rebuild sqlite3

# Copy source code
COPY src/ ./src/
COPY public/ ./public/
COPY frontend/ ./frontend/
COPY tsconfig.json ./

# Build the frontend (Astro + Svelte)
RUN if [ -d "frontend" ] && [ -f "frontend/package.json" ]; then \
      echo "Building frontend..." && \
      (cd frontend && pnpm install && pnpm build) || echo "Frontend build failed, continuing with legacy public/"; \
    fi

# Build the application
RUN pnpm build

# Try to build docs if possible, but don't fail the build if it doesn't work
COPY docs/ ./docs/
RUN if [ -d "docs" ] && [ -f "docs/package.json" ]; then \
      echo "Attempting to build documentation..." && \
      (cd docs && npm install && npm run build) || echo "Documentation build failed, continuing without docs"; \
    else \
      echo "Skipping documentation build"; \
    fi

# Set environment variables
ENV NODE_ENV=production \
    PUPPETEER_SKIP_CHROMIUM_DOWNLOAD=true \
    PUPPETEER_EXECUTABLE_PATH=/usr/bin/chromium \
    DB_TYPE=sqlite

# Create data directory for SQLite
RUN mkdir -p /app/data && chmod 755 /app/data

# Expose port
EXPOSE 3010

# Start the application
CMD ["node", "dist/index.js"]

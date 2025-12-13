# PostgreSQL Setup for Foundry REST API Relay

This setup provides persistent storage using PostgreSQL instead of the default SQLite or memory database.

## Quick Start

1. **Start the services:**
   ```bash
   docker-compose -f docker-compose.postgres.yml up -d
   ```

2. **Initialize the database:**
   ```bash
   # Wait for PostgreSQL to be ready (about 10-15 seconds), then run:
   docker-compose -f docker-compose.postgres.yml exec relay pnpm run db:reset
   ```

3. **Check the logs:**
   ```bash
   docker-compose -f docker-compose.postgres.yml logs relay
   ```

## Full Setup Steps

```bash
# 1. Start PostgreSQL first and wait for it to be ready
docker-compose -f docker-compose.postgres.yml up -d db

# 2. Wait for database to be ready (healthcheck will ensure this)
docker-compose -f docker-compose.postgres.yml ps

# 3. Initialize the database schema and create admin user
docker-compose -f docker-compose.postgres.yml exec relay pnpm run db:reset

# 4. Start the relay service
docker-compose -f docker-compose.postgres.yml up -d relay

# 5. View logs to confirm everything is working
docker-compose -f docker-compose.postgres.yml logs -f relay
```

## Default Database Credentials

- **Database:** `foundryvtt_relay`
- **Username:** `postgres`
- **Password:** `example`
- **Host:** `localhost` (from host machine) or `db` (from within Docker network)
- **Port:** `5432`

## Admin User Setup

To create an admin user during database reset, set environment variables:
```bash
ADMIN_EMAIL=your@email.com ADMIN_PASSWORD=your-secure-password pnpm run db:reset
```

If not provided, no admin user is created. Users can register at the website (e.g., `https://your-relay-server.com`) unless `DISABLE_REGISTRATION=true`.

## Environment Variables

You can customize the PostgreSQL setup by modifying the environment variables in `docker-compose.postgres.yml`:

```yaml
# Database configuration
- POSTGRES_DB=foundryvtt_relay        # Database name
- POSTGRES_USER=postgres              # Database user
- POSTGRES_PASSWORD=example           # Database password

# Relay configuration
- DATABASE_URL=postgres://postgres:example@db:5432/foundryvtt_relay
```

## Stopping and Cleanup

```bash
# Stop services
docker-compose -f docker-compose.postgres.yml down

# Stop and remove all data (WARNING: This deletes all data!)
docker-compose -f docker-compose.postgres.yml down -v
```

## Troubleshooting

### Database Connection Errors
If you see connection errors, make sure:
1. PostgreSQL is running: `docker-compose -f docker-compose.postgres.yml ps`
2. Database is healthy: Check the healthcheck status
3. Wait a bit longer for PostgreSQL to start up completely

### Reinitialize Database
If you need to reset the database:
```bash
# Stop services
docker-compose -f docker-compose.postgres.yml down

# Remove the database volume (WARNING: Deletes all data!)
docker volume rm foundryvtt-rest-api-relay_postgres_data

# Start fresh
docker-compose -f docker-compose.postgres.yml up -d db
docker-compose -f docker-compose.postgres.yml exec relay pnpm run db:reset
docker-compose -f docker-compose.postgres.yml up -d
```

## Data Persistence

All PostgreSQL data is stored in a Docker volume named `foundryvtt-rest-api-relay_postgres_data`. This ensures your data persists between container restarts and updates.

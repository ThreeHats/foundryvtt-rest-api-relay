---
id: update-docker-image
title: Update Docker Image
sidebar_position: 13
---

# Maintenance and Updates of the Docker Image

## 1. Stop the running container

```bash
docker compose -f docker-compose.local.yml down
```

(If you're using the Postgres setup, replace `docker-compose.local.yml` with `docker-compose.postgres.yml`.)

## 2. Pull the latest image

```bash
docker compose -f docker-compose.local.yml pull
```

## 3. Optional: remove old images

List all pulled images:
```bash
docker images -a
```

Remove an old image by ID:
```bash
docker rmi [image_ID]
```

Keeping the previous version around is a good idea — you can roll back by tagging or re-pulling a specific version if the update causes issues.

## 4. Restart with the latest image

```bash
docker compose -f docker-compose.local.yml up -d
```

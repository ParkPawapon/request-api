# request-api

Production-oriented Go API foundation for the Request system.

This repository is intentionally a foundation only. It prepares the backend shape for a future migration from the legacy Express monolith in `../request/backend` without changing the existing database schema or migrating business behavior yet.

## Stack

- Go
- Gin
- GORM
- PostgreSQL
- Redis
- Zap structured logging
- go-playground validator

## Current Endpoints

- `GET /v1/health`
- `GET /v1/health/live`
- `GET /v1/health/ready`

Readiness checks PostgreSQL and Redis. Liveness only reports that the HTTP process is alive.

## Local Commands

```sh
cp .env.example .env
make tidy
make test
make build
make run
```

Do not commit `.env` or real credentials.

## Migration Safety

The application does not run `AutoMigrate`, does not create tables, does not seed data, and does not modify the legacy database schema. Database compatibility work must be planned from the legacy schema first and approved explicitly before any real migration is added.

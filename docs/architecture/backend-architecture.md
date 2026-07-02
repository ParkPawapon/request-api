# Backend Architecture

## Overview

`request-api` is a Go API foundation for a staged migration from the legacy Express backend in `../request/backend`. The current implementation provides infrastructure, routing, health checks, normalized responses, error handling, configuration validation, logging, PostgreSQL and Redis connectivity, and test scaffolding. It does not migrate business logic yet.

The API version boundary is `/v1`. Legacy `/api` and `/api/v1` routes are documented for migration, but this service exposes new routes directly under `/v1`.

## Folder Responsibilities

| Path | Responsibility |
| --- | --- |
| `cmd/api` | Process entrypoint only. It loads config, builds infrastructure, and runs the HTTP server. |
| `internal/app` | Application composition root. Wires config, logger, DB, Redis, router, and server lifecycle. |
| `internal/config` | Environment loading and validation. It must fail fast on unsafe or missing required runtime config. |
| `internal/server` | HTTP server construction and graceful shutdown behavior. |
| `internal/transport/http` | Gin router, middleware, response helpers, and versioned HTTP handlers. |
| `internal/domain` | Future business entities, value objects, and domain contracts. No Gin, GORM, Redis, or transport imports. |
| `internal/usecase` | Future business workflows. No Gin imports. |
| `internal/repository` | Interfaces required by use cases. Implementations live under infrastructure. |
| `internal/infrastructure` | PostgreSQL, Redis, logger, clock, and external system implementations. |
| `internal/pkg` | Small reusable application primitives such as errors, validation, security utilities, and pagination. |
| `migrations` | Placeholder documentation only. No schema-changing migration has been added. |
| `tests` | Integration and fixture placeholders for future DB/container-backed tests. |

## Dependency Direction

Dependencies must point inward:

- `domain` has no dependency on transport, infrastructure, or frameworks.
- `usecase` may depend on `domain` and repository interfaces.
- `repository` defines interfaces; infrastructure implements them.
- `infrastructure` may depend on repository/domain contracts when implementing adapters.
- `transport/http` may call use cases and response helpers, but handlers must not talk to GORM directly.
- `cmd/api` and `internal/app` are composition boundaries and may wire everything together.

## HTTP Flow

1. `cmd/api/main.go` loads validated config.
2. Zap logger is created with environment-aware settings.
3. `internal/app` opens PostgreSQL, creates Redis, verifies connectivity, and builds the Gin router.
4. Gin middleware applies request ID, security headers, CORS, body limit, timeout, access logging, and panic recovery.
5. Versioned handlers return responses through `internal/transport/http/response`.
6. Errors are normalized through `internal/pkg/errors` and never expose stack traces to users.

## Response Convention

Success responses use:

```json
{
  "ok": true,
  "data": {}
}
```

Error responses use:

```json
{
  "ok": false,
  "error": {
    "code": "VALIDATION",
    "message": "Invalid request.",
    "requestId": "..."
  }
}
```

Validation errors include a `fields` array. Pagination metadata is reserved under `meta.pagination`.

## Error Handling Convention

Use `internal/pkg/errors.AppError` for user-facing application errors. Each error contains:

- stable error code
- HTTP status
- safe user-facing message
- wrapped internal error for logs

Handlers must use `response.AppError` or `response.ValidationError`. Raw database, Redis, or framework errors must not be sent directly to users.

## Configuration Convention

Configuration comes from environment variables only. Required values are validated during startup.

Required runtime values:

- `DATABASE_URL`
- `REDIS_ADDR`

Secrets must never be logged. `.env` files are ignored by git and only `.env.example` is committed.

## Logging Convention

Zap is the structured logger. Access logs intentionally avoid Authorization headers, cookies, passwords, tokens, and request bodies. Internal errors may be logged server-side, but user responses stay normalized and safe.

## Database Convention

PostgreSQL is accessed through GORM at the infrastructure boundary. The foundation deliberately does not call `AutoMigrate`, does not create tables, does not seed data, and does not mutate the legacy schema.

Repository implementations must keep GORM models inside infrastructure packages or tightly scoped adapters. Use cases should depend on repository interfaces, not `*gorm.DB`.

The first migrated business route, `GET /v1/petition-types`, follows this boundary:

`transport/http/v1/petitiontypes` -> `usecase/petitiontype` -> `repository.PetitionTypeRepository` -> `infrastructure/database/repository.PetitionTypeGormRepository`.

The GORM record type is private to the infrastructure repository and maps the legacy `"petitionType"` table explicitly before returning a domain entity.

## Rate Limit Convention

Routes that had legacy `express-rate-limit` protection should register a route-level limiter in HTTP transport. `GET /v1/petition-types` uses an in-memory 60 requests/minute limiter to match the legacy default memory-store behavior without adding cache semantics that legacy did not have.

## Redis Convention

Redis is initialized from `REDIS_ADDR`, `REDIS_PASSWORD`, and `REDIS_DB`. The current foundation only provides client creation, ping readiness, and close behavior. No cache behavior is enabled yet, so existing business behavior is not changed.

## Testing Convention

Unit tests should run without a real PostgreSQL or Redis instance unless they are explicitly marked as integration tests. Integration tests belong under `tests/integration` and must document their external dependencies.

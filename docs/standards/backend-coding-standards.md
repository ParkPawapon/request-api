# Backend Coding Standards

## Go Formatting

- Run `gofmt` on every Go file.
- Run `go vet ./...` before release.
- Keep packages small and named after responsibility, not implementation detail.

## Package Boundaries

- `cmd/api` remains an entrypoint only.
- Handlers may not import GORM directly.
- Use cases may not import Gin.
- Domain packages may not import Gin, GORM, Redis, or HTTP transport packages.
- Infrastructure implements external adapters and may depend on domain/repository contracts.

## Context

- Accept `context.Context` as the first parameter for IO, DB, Redis, and external calls.
- Do not store context on structs.
- Honor request cancellation and timeouts.

## Errors

- Use `internal/pkg/errors.AppError` for expected application errors.
- Wrap internal errors for logging but expose only safe messages to users.
- Do not return raw PostgreSQL, Redis, SSO, or filesystem errors in HTTP responses.
- Do not use panic as control flow.

## Validation

- Request DTOs should use struct tags where appropriate.
- Normalize validation failures before responding.
- Validate inputs before reaching repository or DB layers.

## Transactions

- Use repository/usecase-level transaction boundaries for multi-step writes.
- Avoid starting transactions in HTTP handlers.
- Keep transaction scope as small as possible.

## Logging

- Use structured Zap logging.
- Never log passwords, tokens, Authorization headers, cookies, session IDs, CSRF tokens, or raw credentials.
- Log request IDs with access logs and error logs.

## Security

- Do not commit `.env` or real secrets.
- Do not hardcode database URLs, Redis URLs, tokens, passwords, or SSO secrets.
- Do not fake authentication in production code.
- Keep CORS allowlists explicit.
- Keep request body limits enabled.
- Use normalized error responses.

## Testing

- Unit tests must not require live PostgreSQL or Redis.
- Integration tests that require external services belong under `tests/integration` and must document setup.
- Do not skip failing tests to make a build pass.
- Add behavior tests with each migrated endpoint.

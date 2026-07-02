# Production Readiness Checklist

Scope: `request-api` production readiness audit for the Go/Gin `/v1` API foundation.

| Checked item | Status | Evidence | File/route | Remediation | Remaining risk |
| --- | --- | --- | --- | --- | --- |
| Repository scope limited to `request-api` | Pass | Worktree diff only contains paths under this repository. | Repository root | Continue using explicit path staging. | None identified. |
| API version boundary | Pass | Router registers the API group under `/v1`; no legacy route group is exposed. | `internal/transport/http/router.go` | Keep new handlers under versioned packages. | Future routes need review. |
| Health and readiness endpoints | Pass | `/v1/health`, `/v1/health/live`, and `/v1/health/ready` are registered and covered by tests. | `internal/transport/http/v1/health` | Add dependency-specific checks as services are added. | Readiness depends on configured DB/Redis availability. |
| Config fail-fast behavior | Pass | Startup validation rejects invalid env, timeouts, DB pool settings, Redis DB, unsafe Redis prefix, and production/staging without CORS allowlist. | `internal/config/config.go` | Keep env validation in `internal/config` only. | Future env vars need tests. |
| Normalized error responses | Pass | `NoRoute` and handler errors use the response envelope and safe messages. | `internal/transport/http/response`, `internal/transport/http/router_test.go` | Route-specific validation should use `response.ValidationError`. | Future handlers need review. |
| Database schema safety | Pass | Runtime opens and pings PostgreSQL only; no schema mutation path is present. | `internal/infrastructure/database` | Require explicit migration approval for any schema change. | Business migration may require table-specific contract tests. |
| Redis/cache safety | Pass | Redis is initialized from env and used only for readiness; no cache behavior changes business results. | `internal/infrastructure/cache` | Apply `REDIS_KEY_PREFIX` when cache keys are introduced. | Prefix is validated but not yet consumed by a cache feature. |
| Observability baseline | Pass | Request ID, structured access logs, recovery logs, redacted startup config logging, and health endpoints exist. | `cmd/api/main.go`, `internal/transport/http/middleware` | Add metrics/tracing when deployment target is known. | No metrics exporter yet. |
| CI quality gates | Pass | GitHub Actions workflow runs tidy verification, gofmt check, vet, tests, and build. | `.github/workflows/request-api-ci.yml` | Enable required status check after GitHub observes the workflow. | Required status context may not exist until first run completes. |
| Release documentation | Pass | QA and release checklists are included in docs. | `docs/qa`, `docs/release` | Keep reports current per release. | Manual signoff still required before production deploy. |


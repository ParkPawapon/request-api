# Performance Report

Scope: runtime defaults, middleware, database access, and migrated route behavior.

| Checked item | Status | Evidence | File/route | Remediation | Remaining risk |
| --- | --- | --- | --- | --- | --- |
| Request timeout | Pass | Router applies request timeout middleware and config validates positive duration. | `internal/transport/http/middleware/timeout.go`, `internal/config/config.go` | Tune per deployment. | Long-running future endpoints need specific review. |
| Server timeouts | Pass | HTTP server uses read and write timeout from config. | `internal/server/http_server.go` | Tune with production traffic profile. | No load test has been run. |
| DB connection pool | Pass | Pool size and lifetime are configurable and validated. | `internal/infrastructure/database/postgres.go` | Load test before increasing limits. | Defaults may not match production capacity. |
| GORM query behavior | Pass | `GET /v1/petition-types` performs one ordered read query and does not preload or join. | `internal/infrastructure/database/repository/petition_type_gorm_repository.go` | Add pagination before migrating large lists. | Current migrated route is reference data only. |
| Redis behavior | Pass | Redis is used only for readiness, not hot-path caching. | `internal/infrastructure/cache/redis.go` | Add cache benchmarks when cache behavior is introduced. | No cache performance benefit exists yet. |
| Middleware order | Pass | Request ID, security headers, CORS, body limit, timeout, logger, and recovery are consistently applied before routes. | `internal/transport/http/router.go` | Keep route-specific middleware close to route registration. | Rate limiter is in-memory and per process. |
| Payload size | Pass | Current migrated route returns compact reference data. | `/v1/petition-types` | Add pagination metadata for large datasets. | Future list endpoints may need paging. |
| Bundle/runtime analysis | Blocked | No production load test, profiler, or benchmark suite exists yet. | Repository scan | Add k6/hey profiles or Go benchmarks after deployment shape is known. | Production capacity is not proven by unit tests alone. |


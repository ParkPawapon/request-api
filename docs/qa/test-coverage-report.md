# Test Coverage Report

Scope: local test and regression gates for the current backend foundation.

| Checked item | Status | Evidence | File/route | Remediation | Remaining risk |
| --- | --- | --- | --- | --- | --- |
| Config validation tests | Pass | Tests cover valid config, invalid env, missing database URL, invalid CORS, production CORS requirement, invalid pools/timeouts, and unsafe Redis prefix. | `internal/config/config_test.go` | Add tests for new env vars when introduced. | No table-driven exhaustive config matrix yet. |
| Health handler tests | Pass | Tests cover live success, ready success, and normalized readiness failure. | `internal/transport/http/v1/health/handler_test.go` | Add router-level health tests when more middleware is added. | DB/Redis checks are mocked, not integration-tested. |
| Petition type route tests | Pass | Tests cover success envelope, normalized usecase error, fail-closed missing usecase, and usecase error normalization. | `internal/transport/http/v1/petitiontypes`, `internal/usecase/petitiontype` | Add repository contract tests with isolated DB. | No real database test yet. |
| Middleware tests | Pass | Rate limiter and router `NoRoute` normalized error behavior are covered. | `internal/transport/http/middleware/rate_limit_test.go`, `internal/transport/http/router_test.go` | Add tests for CORS/security headers if behavior changes. | Request timeout/recovery tests are not exhaustive. |
| CI gates | Pass | CI runs module tidy verification, gofmt check, vet, tests, and build. | `.github/workflows/request-api-ci.yml` | Enable branch protection required check after first workflow run. | Security tools are recommendations until installed. |
| Race tests | Pass | `go test -race ./...` is part of local release verification for this hardening branch. | Local verification | Consider adding race gate to CI if runtime cost is acceptable. | CI workflow currently omits race testing for speed. |
| Vulnerability tests | Blocked | `govulncheck` and `gosec` are not installed locally. | Local toolchain | Install tools and optionally add a scheduled CI security workflow. | Known-vulnerability scan not completed in this gate. |


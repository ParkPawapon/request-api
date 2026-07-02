# Request API Release Checklist

Release target: `request-api` Go/Gin `/v1` backend foundation.

| Checked item | Status | Evidence | File/route | Remediation | Remaining risk |
| --- | --- | --- | --- | --- | --- |
| Source branch naming | Pass | Release branch is `security/request-api-production-hardening`. | Git branch | Continue avoiding forbidden prefixes. | None identified. |
| Conventional commits | Pass | Planned commits use `security`, `test`, `ci`, and `docs` Conventional Commit types. | Git history | Review commit log before merge. | None identified. |
| Quality gates | Pass | Local gates include tidy, gofmt, test, race test, vet, build, make test, make build, and make lint where available. | Local verification | Keep failed gates blocking merges. | `govulncheck` and `gosec` unavailable locally. |
| CI workflow | Pass | `Request API CI` workflow added for PRs and pushes to `main`. | `.github/workflows/request-api-ci.yml` | Require `quality-gates` in branch protection after first run appears. | GitHub must observe the new workflow context. |
| Database safety | Pass | No migration, schema mutation, seed, or runtime `AutoMigrate` is included. | `docs/qa/database-safety-report.md` | Keep DB changes in separate approved tasks. | Business logic migration may need explicit schema contract tests. |
| Security review | Pass | Security audit report documents pass/blocker status and secret-scan evidence. | `docs/qa/security-audit-report.md` | Add GitHub secret scanning/security tooling if not already enabled. | Local vulnerability scanners were not installed. |
| Observability review | Pass | Logging, request ID, readiness, and safe startup metadata are documented. | `docs/qa/observability-report.md` | Add metrics/tracing in platform integration. | SLO monitoring not complete. |
| Performance review | Pass | Timeout, pool, query, payload, and middleware checks documented. | `docs/qa/performance-report.md` | Run load tests before production traffic. | No benchmark/load report yet. |
| Branch protection | Pass | Existing branch protection should remain enabled with force push/deletion disabled and merge commits allowed. | GitHub branch protection | Add required status check when workflow context is available. | Explicit user bypass configuration may be unsupported on user-owned repo. |
| Deployment approval | Blocked | This prompt does not deploy and does not change infrastructure. | Release process | Run deployment-specific readiness after environment is chosen. | Production deploy remains out of scope. |


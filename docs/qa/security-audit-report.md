# Security Audit Report

Scope: source-level security hardening review for `request-api`.

| Checked item | Status | Evidence | File/route | Remediation | Remaining risk |
| --- | --- | --- | --- | --- | --- |
| Secrets and `.env` files | Pass | No `.env`, private key, certificate, or credential file is present in the repository scan; only `.env.example` is committed. | `.env.example`, `.gitignore` | Keep real secrets in deployment secret stores. | Secret scanning should run in GitHub as an additional control. |
| Hardcoded database or Redis credentials | Pass | Runtime uses `DATABASE_URL`, `REDIS_ADDR`, `REDIS_PASSWORD`, and `REDIS_DB`; `.env.example` contains placeholders only. | `internal/config/config.go` | Do not add literal production URLs or passwords. | Manual review remains necessary for future commits. |
| Public error exposure | Pass | App errors map internal errors to safe response messages; readiness failure exposes only `Service is not ready.` | `internal/pkg/errors`, `internal/transport/http/response` | Keep internal details in logs only. | New handlers must avoid returning raw errors. |
| Panic recovery | Pass | Recovery middleware logs panic server-side and returns `Internal Server Error`. | `internal/transport/http/middleware/recovery.go` | Avoid using panic for control flow. | Panic payload can contain sensitive internal detail in server logs; restrict log access. |
| Request logging redaction | Pass | Logger records request ID, method, route path, status, latency, and client IP only; it does not log body, cookies, Authorization header, or passwords. | `internal/transport/http/middleware/logger.go` | Keep request body logging disabled by default. | Downstream middleware must follow same rule. |
| Startup config logging | Pass | Startup log records safe metadata and a boolean for database URL presence, not secret values. | `cmd/api/main.go` | Do not log DSNs or Redis password. | Review added config fields before logging. |
| CORS policy | Pass | CORS uses an allowlist and production/staging config validation rejects empty origins. | `internal/config/config.go`, `internal/transport/http/middleware/cors.go` | Keep wildcard origins out of production. | Origin list must be managed per environment. |
| Request body size limit | Pass | Router applies body limit middleware and config rejects non-positive limits. | `internal/transport/http/router.go`, `internal/config/config.go` | Tune size per API feature. | File upload routes will need stricter validation. |
| Fake auth or fake success | Pass | Auth middleware placeholder fails closed if authenticator is nil; no fake authenticated principal is returned. | `internal/transport/http/middleware/auth.go` | Implement real auth facade before protected routes. | Auth is not yet migrated. |
| Open redirect and file upload risk | Pass | No redirect helper or upload route exists in current implementation. | Source scan | Add explicit validation when those features are migrated. | Future features may introduce risk. |
| Vulnerability tools | Blocked | `govulncheck` and `gosec` binaries are not installed locally. | Local toolchain | Install and run `govulncheck ./...` and `gosec ./...` in a follow-up gate. | Known-vulnerability scan was not claimed as passed. |


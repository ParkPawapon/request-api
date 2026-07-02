# Observability Report

Scope: logging, correlation, health, and operational diagnostics foundation.

| Checked item | Status | Evidence | File/route | Remediation | Remaining risk |
| --- | --- | --- | --- | --- | --- |
| Structured logging | Pass | Zap logger is used for application, request, recovery, and GORM warning logs. | `internal/infrastructure/logger`, `internal/infrastructure/database/gorm.go` | Add deployment log routing. | No centralized log backend configured in this repo. |
| Request correlation | Pass | Request ID middleware accepts or generates `X-Request-ID` and writes it to the response header. | `internal/transport/http/middleware/request_id.go` | Propagate request ID into downstream calls when added. | No distributed trace ID yet. |
| Access logs | Pass | Access logger records route path, method, status, latency, client IP, and request ID. | `internal/transport/http/middleware/logger.go` | Add status buckets/metrics when observability backend exists. | `FullPath()` may be empty for unmatched routes. |
| Panic visibility | Pass | Recovery middleware logs recovered panic and request ID. | `internal/transport/http/middleware/recovery.go` | Keep panic response generic. | Panic details require restricted log access. |
| Startup config visibility | Pass | Startup log records safe runtime metadata and avoids secret values. | `cmd/api/main.go` | Extend only with non-sensitive fields. | Operators must correlate this with deployment config. |
| Health and readiness | Pass | Liveness and readiness are separate; readiness checks DB and Redis. | `/v1/health/live`, `/v1/health/ready` | Add dependency-specific degraded states if needed. | Public readiness exposure should be controlled at ingress if required. |
| Metrics and tracing | Blocked | No metrics or tracing exporter exists in current foundation. | Repository scan | Add OpenTelemetry or Prometheus once platform target is selected. | Production SLO monitoring is not complete without metrics. |


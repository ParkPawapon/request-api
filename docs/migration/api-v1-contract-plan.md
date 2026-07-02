# API v1 Contract Plan

The new API version boundary is `/v1`. Legacy routes under `/api/v1` and `/api` should migrate to `/v1` equivalents. This document is a plan only; business handlers are not implemented yet except health.

| Legacy API path | New `/v1` path | Method | Request DTO | Response DTO | Auth required | Migration status | Notes |
| --- | --- | --- | --- | --- | --- | --- | --- |
| `/api/v1/healthz` | `/v1/health` | GET | None | `HealthStatus` | No | Implemented foundation | New route returns normalized envelope and checks DB/Redis. |
| `/api/v1/livez` | `/v1/health/live` | GET | None | `HealthStatus` | No | Implemented foundation | Liveness does not require DB/Redis. |
| `/api/v1/readyz` | `/v1/health/ready` | GET | None | `HealthStatus` | No | Implemented foundation | Readiness checks DB/Redis. Upload/session checks are not migrated yet. |
| `/api/v1/auth/csrf` | `/v1/auth/csrf` | GET | None | `CsrfTokenResponse` | Session | Planned | Requires session/CSRF strategy decision. |
| `/api/v1/auth/login` | `/v1/auth/login` | POST | `LoginRequest` | `CurrentUserResponse` or SSO redirect hint | Session/CSRF | Planned | Must preserve SSO redirect behavior. |
| `/api/v1/auth/dev/options` | `/v1/auth/dev/options` | GET | None | `DevLoginOptionsResponse` | Development only | Planned | Must be disabled in production. |
| `/api/v1/auth/dev/login` | `/v1/auth/dev/login` | POST | `DevLoginRequest` | `CurrentUserResponse` | Development only | Planned | Must not become fake production auth. |
| `/api/v1/auth/sso/login` | `/v1/auth/sso/login` | GET | Query params | Redirect or error page | Session | Planned | Needs openid-client equivalent or explicit SSO adapter. |
| `/api/v1/auth/sso/callback` | `/v1/auth/sso/callback` | GET | Query params | Redirect or error page | Session | Planned | Must validate state/nonce. |
| `/api/v1/auth/sso/claims` | `/v1/auth/sso/claims` | GET | None | HTML or JSON claims | SSO/session | Planned | Debug exposure controls must be preserved. |
| `/api/v1/auth/sso/debug` | `/v1/auth/sso/debug` | GET | None | HTML or JSON debug claims | SSO/session | Planned | Must remain disabled unless configured. |
| `/api/v1/auth/me` | `/v1/auth/me` | GET | None | `CurrentUserResponse` | Session | Planned | Must match user/role shape from legacy. |
| `/api/v1/auth/logout` | `/v1/auth/logout` | POST | None | `LogoutResponse` | Session/CSRF | Planned | Must clear same cookie/session state as legacy. |
| `/api/v1/petition-types` | `/v1/petition-types` | GET | None | `PetitionTypesResponse` | Verify | Planned | Read from legacy `petitionType` table. |
| `/api/v1/petitions/my` | `/v1/petitions/my` | GET | Query params | `PetitionListResponse` | User | Planned | Must preserve filtering and paging behavior. |
| `/api/v1/petitions` | `/v1/petitions` | POST | `CreatePetitionRequest` multipart | `PetitionResponse` | User | Planned | Requires upload adapter and DB transaction design. |
| `/api/v1/petitions/:id` | `/v1/petitions/:id` | GET | Path param | `PetitionResponse` | User | Planned | Must enforce submitter visibility. |
| `/api/v1/petitions/:id/attachments` | `/v1/petitions/:id/attachments` | GET | Path param | `AttachmentListResponse` | User | Planned | Must enforce authorization. |
| `/api/v1/petitions/:id/attachments/:attId/download` | `/v1/petitions/:id/attachments/:attId/download` | GET | Path params | File stream | User | Planned | Must preserve safe path resolution. |
| `/api/v1/petitions/:id/cancel` | `/v1/petitions/:id/cancel` | POST | Path param | `PetitionResponse` | User | Planned | Must preserve status transition rules. |
| `/api/v1/petitionsLecturers/my` | `/v1/petitionsLecturers/my` | GET | Query params | `PetitionListResponse` | Lecturer | Planned | Route casing kept for parity until frontend contract changes. |
| `/api/v1/petitionsLecturers` | `/v1/petitionsLecturers` | POST | `CreatePetitionRequest` multipart | `PetitionResponse` | Lecturer | Planned | Requires lecturer identity mapping. |
| `/api/v1/staff/requests` | `/v1/staff/requests` | GET | Query params | `StaffRequestListResponse` | Staff | Planned | Must preserve dashboard/list filters. |
| `/api/v1/staff/requests/:id` | `/v1/staff/requests/:id` | GET | Path param | `StaffRequestDetailResponse` | Staff | Planned | Must enforce staff role permissions. |
| `/api/v1/staff/requests/:id/status` | `/v1/staff/requests/:id/status` | PATCH | `UpdateRequestStatusRequest` | `StaffRequestDetailResponse` | Staff | Planned | Requires transaction and audit behavior review. |
| `/api/v1/staff-dashboard` | `/v1/staff-dashboard` | GET | Query params | `StaffDashboardResponse` | Staff | Planned | Aggregate query parity required. |

## Contract Rules

- Endpoint constants should be centralized when business handlers are added.
- Handlers must call use cases, not GORM directly.
- Response envelopes must go through `internal/transport/http/response`.
- Errors must use `internal/pkg/errors`.
- Production routes must not fake success when dependencies are not implemented.

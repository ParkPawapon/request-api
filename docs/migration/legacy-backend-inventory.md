# Legacy Backend Inventory

This inventory was created from read-only inspection of `../request/backend` and `../request/migrations`.

## Files and Folders Read

- `../request/backend/package.json`
- `../request/backend/src/app.js`
- `../request/backend/src/server.js`
- `../request/backend/src/config/api.js`
- `../request/backend/src/config/env.js`
- `../request/backend/src/config/sessionConfig.js`
- `../request/backend/src/config/ssoConfig.js`
- `../request/backend/src/db/index.js`
- `../request/backend/src/middlewares/csrf.js`
- `../request/backend/src/routes/auth.js`
- `../request/backend/src/routes/_auth.js`
- `../request/backend/src/routes/petitionTypes.js`
- `../request/backend/src/routes/petitions.js`
- `../request/backend/src/routes/petitionsLecturers.js`
- `../request/backend/src/routes/staffDashboard.js`
- `../request/backend/src/routes/staffManagePetitions.js`
- `../request/backend/src/routes/staffRequests.js`
- `../request/backend/src/controllers/ssoController.js`
- `../request/backend/src/services/authService.js`
- `../request/backend/src/services/petitionService.js`
- `../request/backend/src/services/ssoClient.js`
- `../request/backend/src/services/uploadStorage.js`
- `../request/backend/src/repositories/petitionRepository.js`
- `../request/backend/src/utils/claims.js`
- `../request/migrations/20250122010000_init_schema.sql`

## Legacy Runtime Shape

- Runtime: Node.js ESM
- HTTP framework: Express 4
- Database driver: `pg`
- Session store: `express-session` with `connect-pg-simple`
- Security middleware: `helmet`, CORS allowlist, cookie parser, same-origin and CSRF middleware
- Uploads: `multer`, disk staging, file type checks through upload storage service
- Validation: `zod`, `express-validator`, and route-level checks
- Auth: session-based local/dev auth plus KMUTT SSO via `openid-client`
- Logging: `morgan` plus console warnings/errors
- API prefixes: both `/api` and `/api/v1` are mounted in legacy

## Route Inventory

| Legacy route | Method | Purpose | Auth |
| --- | --- | --- | --- |
| `/api/v1/healthz`, `/api/healthz` | GET | Readiness alias | Public |
| `/api/v1/readyz`, `/api/readyz` | GET | DB, upload storage, session readiness | Public |
| `/api/v1/livez`, `/api/livez` | GET | Process liveness | Public |
| `/api/v1/auth/csrf` | GET | Issue CSRF token | Session |
| `/api/v1/auth/login` | POST | Local/dev login or SSO redirect hint | Session/CSRF |
| `/api/v1/auth/dev/options` | GET | Development login options | Development only |
| `/api/v1/auth/dev/login` | POST | Development login | Development only |
| `/api/v1/auth/sso/login` | GET | Start SSO authorization flow | Session |
| `/api/v1/auth/sso/callback` | GET | Complete SSO callback | Session |
| `/api/v1/auth/sso/claims` | GET | SSO claims page/API depending accept header | SSO/session |
| `/api/v1/auth/sso/debug` | GET | Debug SSO claims when enabled | SSO/session |
| `/api/v1/auth/sso/health` | GET | Development SSO health | Development only |
| `/api/v1/auth/sso/diagnose` | GET | Development SSO config diagnostics | Development only |
| `/api/v1/auth/me` | GET | Current user session | Session |
| `/api/v1/auth/logout` | POST | Destroy session and clear cookie | Session/CSRF |
| `/api/v1/petition-types` | GET | List petition types | Public/session behavior to verify |
| `/api/v1/petitions/my` | GET | List current student's petitions | Auth required |
| `/api/v1/petitions` | POST | Create student petition with attachments | Auth required |
| `/api/v1/petitions/:id` | GET | Get student petition detail | Auth required |
| `/api/v1/petitions/:id/attachments` | GET | List student petition attachments | Auth required |
| `/api/v1/petitions/:id/attachments/:attId/download` | GET | Download petition attachment | Auth required |
| `/api/v1/petitions/:id/cancel` | POST | Cancel student petition | Auth required |
| `/api/v1/petitionsLecturers/my` | GET | List lecturer petitions | Auth required |
| `/api/v1/petitionsLecturers` | POST | Create lecturer petition with attachments | Auth required |
| `/api/v1/petitionsLecturers/:id` | GET | Get lecturer petition detail | Auth required |
| `/api/v1/petitionsLecturers/:id/attachments` | GET | List lecturer petition attachments | Auth required |
| `/api/v1/petitionsLecturers/:id/attachments/:attId/download` | GET | Download lecturer attachment | Auth required |
| `/api/v1/petitionsLecturers/:id/cancel` | POST | Cancel lecturer petition | Auth required |
| `/api/v1/staff/requests` | GET | Staff request list | Staff required |
| `/api/v1/staff/requests/:id` | GET | Staff request detail | Staff required |
| `/api/v1/staff/requests/:id/attachments` | GET | Staff attachment list | Staff required |
| `/api/v1/staff/requests/:id/attachments/:attId/download` | GET | Staff attachment download | Staff required |
| `/api/v1/staff/requests/:id/status` | PATCH | Staff status update | Staff required |
| `/api/v1/staff/requests/:id/attachments` | POST | Staff add attachments | Staff required |
| `/api/v1/staff-dashboard` | GET | Staff dashboard aggregate | Staff required |

## Middleware Inventory

- `helmet` with development CSP disabled and production defaults
- `compression`
- CORS with `CORS_ORIGINS` or `CORS_ORIGIN`
- `morgan`
- `express.json` and `express.urlencoded` with 50 MB limit
- `cookie-parser`
- `express-session`
- same-origin check plus CSRF verification for non-health endpoints
- route-level rate limiters in auth, petitions, and staff routes
- multer upload handling for petition attachments

## Database Usage Inventory

Legacy DB access uses `pg.Pool` and repository/service functions. The initial schema migration defines these tables:

- `student`
- `lecturer`
- `staff`
- `lecturerRole`
- `staffRole`
- `petition`
- `petitionType`
- `petitionSubmitter`
- `petitionAttachment`
- `session`
- `thesystem`

The legacy server may create the `session` table through `connect-pg-simple` when `SESSION_TABLE_CREATE=true`. The Go service must not create or alter tables by default.

## Auth and Session Inventory

- Primary auth state is Express session cookie, default name `sid`.
- Production requires `SESSION_SECRET` of at least 16 characters.
- Session cookie options are controlled by `SESSION_COOKIE_SAMESITE`, `SESSION_COOKIE_SECURE`, `SESSION_MAX_AGE_MS`, and proxy settings.
- CSRF token is stored in session and verified on mutating requests.
- SSO config is controlled by `SSO_*` variables.
- Development-only local auth can be enabled through `LOCAL_AUTH_ENABLED`.

## Validation Inventory

- Auth routes use `express-validator` for login payloads.
- Petition routes use route-level validation and service-level checks.
- SSO config uses explicit parsing and normalization.
- Upload validation includes root path containment and file type/size handling.

## Environment Inventory

Important legacy variables include:

- `PORT`, `NODE_ENV`, `TRUST_PROXY`
- `CORS_ORIGIN`, `CORS_ORIGINS`
- `SESSION_SECRET`, `SESSION_NAME`, `SESSION_TABLE_CREATE`
- `SESSION_COOKIE_SAMESITE`, `SESSION_COOKIE_SECURE`, `SESSION_MAX_AGE_MS`
- `PGHOST`, `PGPORT`, `PGDATABASE`, `PGUSER`, `PGPASSWORD`
- `PG_POOL_MAX`, `PG_CONNECT_MS`, `PG_IDLE_MS`, `PG_STMT_MS`, `PG_IDLE_TX_MS`
- `UPLOAD_ROOT`, `UPLOAD_TMP_ROOT`
- `LOCAL_AUTH_ENABLED`, `ALLOW_PLAINTEXT_LOGIN`
- `SSO_ENABLED`, `SSO_ISSUER_URL`, `SSO_CLIENT_ID`, `SSO_CLIENT_SECRET`, `SSO_REDIRECT_URI`, `SSO_SCOPE`
- `SSO_USE_MANUAL_ENDPOINTS`, `SSO_AUTHORIZATION_ENDPOINT`, `SSO_TOKEN_ENDPOINT`, `SSO_USERINFO_ENDPOINT`, `SSO_JWKS_URI`
- `SSO_EXPOSE_CLAIMS`, `SSO_INCLUDE_CLAIMS_IN_RESPONSE`, `SSO_CLAIMS_MASK`, `SSO_DEBUG_CLAIMS`

## Risks and Unknowns

- Full API response contracts vary by route and need endpoint-by-endpoint parity tests before business migration.
- Session and CSRF behavior must be mapped carefully; the Go service currently has no production auth implementation.
- Upload storage behavior is significant for readiness and petition workflows but is intentionally not implemented in this foundation.
- Rate limiting exists in legacy routes and is not yet enabled in Go.
- Redis is part of the new foundation, but legacy backend did not show active Redis usage in inspected source.
- Legacy supports both `/api` and `/api/v1`; new API intentionally standardizes on `/v1`.

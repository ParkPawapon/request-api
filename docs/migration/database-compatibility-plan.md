# Database Compatibility Plan

## Current Rule

The legacy database schema is the source of truth. `request-api` must not change the schema until a separate explicit migration task is approved.

The current foundation:

- does not run `AutoMigrate`
- does not create tables
- does not alter tables
- does not drop tables
- does not seed data
- only opens PostgreSQL and performs readiness pings
- validates connection pool values before startup

## Legacy Schema Source

The inspected schema source is:

- `../request/migrations/20250122010000_init_schema.sql`

Primary tables found:

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

## Migrated Table: `petitionType`

Source:

- `../request/backend/src/routes/petitionTypes.js`
- `../request/migrations/20250122010000_init_schema.sql`

Schema facts from source:

| Column | Legacy SQL type | Nullability/default from source | Go mapping | Notes |
| --- | --- | --- | --- | --- |
| `"petitionTypeID"` | `integer` | `NOT NULL`, identity sequence added later | `int` | Primary key. |
| `"petitionTypeName"` | `character varying` | nullable; unique constraint | `string` after `IS NOT NULL` filter | Route filters null names before mapping. |

Constraints from source:

- Primary key: `"petitionType_pkey"` on `"petitionTypeID"`
- Unique key: `"petitionType_petitionTypeName_key"` on `"petitionTypeName"`
- Identity sequence: `"petitionType_petitionTypeID_seq"` on `"petitionTypeID"`

Runtime behavior for migrated route:

- Read-only query.
- Explicit table name: `"petitionType"` through `TableName()`.
- Explicit column tags: `column:petitionTypeID`, `column:petitionTypeName`.
- Filter: `"petitionTypeName" IS NOT NULL`.
- Ordering: `"petitionTypeName" ASC`.
- No soft delete behavior.
- No timestamp behavior.
- No pagination.
- No transaction required.
- No schema creation or alteration.

## Inspection Process Before Model Migration

Before adding Go models or repository implementations:

1. Export schema from a non-production database.
2. Compare the export to `../request/migrations/20250122010000_init_schema.sql`.
3. Identify table names, quoted identifiers, column casing, nullability, defaults, foreign keys, and indexes.
4. Map legacy query behavior from `../request/backend/src/repositories` and `../request/backend/src/services`.
5. Add GORM model tags only after confirming exact table and column names.
6. Add repository tests against an isolated test database or container.

## GORM Safety Rules

- Never call `AutoMigrate` in application startup.
- Never run schema-changing code from HTTP handlers.
- Never run `DropTable`, `CreateTable`, `AlterColumn`, or seed code from default runtime.
- Keep GORM models and DB-specific tags in infrastructure or persistence packages.
- Use explicit table names when legacy names are quoted or mixed case.
- Use transactions for multi-table petition and attachment writes.
- Normalize DB errors before sending HTTP responses.
- For `"petitionType"`, keep the GORM record model inside `internal/infrastructure/database/repository` and map out to the domain entity before returning to use cases.

## Release Audit Evidence

The production hardening gate records database safety status in `docs/qa/database-safety-report.md`. Each release must confirm that runtime code still avoids schema mutation and that any new repository maps legacy tables and columns explicitly.

## Migration Folder Policy

`migrations/` currently contains documentation only. Future migrations must include:

- migration ID
- purpose
- forward SQL
- rollback SQL where feasible
- risk assessment
- data backup expectation
- manual approval evidence

## Rollback Expectation

Business migration work should be reversible at the application layer first. Schema migrations require a tested rollback or documented manual recovery path before release.

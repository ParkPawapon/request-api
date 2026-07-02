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
- Keep GORM models and DB-specific tags in infrastructure or persistence packages.
- Use explicit table names when legacy names are quoted or mixed case.
- Use transactions for multi-table petition and attachment writes.
- Normalize DB errors before sending HTTP responses.

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

# Implementation Plan: Highlights and Establishment Dashboard

Branch: `005-highlights-dashboard`  
Spec: `specs/005-highlights-dashboard/spec.md`  
Date: 2026-05-13

## Summary

Implement Phase 5 as a backend-first vertical slice. Start with a dedicated
`highlights` table plus use cases and routes for create/delete highlight, then
extend establishment detail and add a minimal stats endpoint. Reuse existing
review data shape so the client can render highlights without inventing a second
review representation.

## Technical Context

- Language/Version: Go
- Storage: PostgreSQL
- Testing: `go test`
- Entry Points:
  - `internal/application/usecases`
  - `internal/infra/adapters/inbound/controller`
  - `internal/infra/adapters/outbound/repository/postgres`
  - `migrations/`
- Constraints:
  - preserve the credibility model already established for reviews
  - avoid large route or repository rewrites
  - keep MVP authorization explicit even if simplified

## Constitution Check

- Domain integrity preserved: yes, rules should live in use cases and domain
  objects before hitting persistence.
- HTTP contracts defined: yes, spec includes create/delete highlight and stats.
- DTO persistence boundaries respected: yes, new persistence code should follow
  the DTO patterns already aligned in review, checkin, like, and user flows.
- Incremental MVP slice: yes, start with minimal dashboard stats and highlight
  management only.
- Acceptance scenarios identified: yes, five scenarios captured in `spec.md`.

## Design

### Data Model

- Add `highlights` table with:
  - `establishment_id`
  - `review_id`
  - `created_by_user_id`
  - `created_at`
- Add uniqueness to prevent duplicate highlights.
- Add indexes to support listing highlights by establishment efficiently.

### Application Layer

- Add `Highlight` domain entity and repository contract.
- Add `CreateHighlightUseCase`.
- Add `DeleteHighlightUseCase`.
- Add `GetEstablishmentStatsUseCase`.
- Extend establishment detail retrieval to include highlighted reviews.
- Keep authorization policy behind a dedicated check so it can evolve later.

### Infrastructure Layer

- Add Postgres highlight repository.
- Extend Postgres establishment repository with:
  - `GetByID` highlight enrichment, or a dedicated detail query path
  - `Stats` query for aggregate metrics
- Add controller DTOs and handlers for highlight operations and stats.
- Wire routes and runner dependencies.

## Rollout Sequence

1. Add migration and repository interfaces.
2. Add domain entities and use case tests.
3. Implement Postgres repositories.
4. Add handlers and route wiring.
5. Validate targeted tests and manual API flows.

## Risks

- Authorization is underspecified for establishment operators.
  - Mitigation: isolate the check behind a clear use case dependency or explicit
    placeholder rule.
- Establishment detail payload may grow if review enrichment is duplicated.
  - Mitigation: reuse review DTO/domain mapping where possible.
- Stats query may duplicate existing review aggregation logic.
  - Mitigation: keep one aggregate query per endpoint and consolidate later only
    if duplication becomes real.

## Artifacts

- `specs/005-highlights-dashboard/spec.md`
- `specs/005-highlights-dashboard/tasks.md`
- Optional follow-up docs if needed:
  - `specs/005-highlights-dashboard/data-model.md`
  - `specs/005-highlights-dashboard/contracts/`

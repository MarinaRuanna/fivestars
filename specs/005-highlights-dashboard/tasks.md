# Tasks: Highlights and Establishment Dashboard

## Phase 1. Discovery and Schema

- [ ] T001 Define the `Highlight` and `EstablishmentStats` domain shapes.
- [ ] T002 Decide the temporary MVP authorization rule for establishment
  operators and record it in the spec or plan.
- [ ] T003 Add migration for the `highlights` table and supporting indexes.

## Phase 2. Repository Contracts

- [ ] T004 Add highlight repository interface(s) in `internal/domain`.
- [ ] T005 Extend establishment and/or review repository contracts for dashboard
  stats and highlighted review retrieval.

## Phase 3. Use Cases

- [ ] T006 Add failing unit tests for `CreateHighlightUseCase` covering
  same-establishment, duplicate, limit, and authorization rules.
- [ ] T007 Implement `CreateHighlightUseCase` until the tests pass.
- [ ] T008 Add failing unit tests for `DeleteHighlightUseCase` and
  `GetEstablishmentStatsUseCase`.
- [ ] T009 Implement `DeleteHighlightUseCase` and `GetEstablishmentStatsUseCase`
  until the tests pass.

## Phase 4. Persistence

- [ ] T010 Add failing tests or verification harnesses for Postgres highlight
  repository behavior where practical.
- [ ] T011 Implement Postgres highlight repository with DTO mapping and error
  translation, then extend establishment/review repositories to return
  highlighted reviews and aggregate stats.
- [ ] T012 Refactor persistence code and complete repository-level validation.

## Phase 5. HTTP Delivery

- [ ] T013 Add failing controller/HTTP tests for create/delete highlight and
  stats responses.
- [ ] T014 Add controller DTOs and handlers for:
  - `POST /establishments/:id/highlights`
  - `DELETE /establishments/:id/highlights/:reviewId`
  - `GET /establishments/:id/stats`
- [ ] T015 Extend establishment detail responses to include `highlights`.
- [ ] T016 Wire routes and dependencies in the runner.

## Phase 6. Validation

- [ ] T017 Run targeted `go test` for the affected use cases and repositories.
- [ ] T018 Manually validate the acceptance scenarios against the API.
- [ ] T019 Update `docs/EXECUTION_PLAN.md` or related product docs if scope
  decisions changed during implementation.

## Phase 7. Claim Security Refactor

- [ ] T020 Update `spec.md` and `data-model.md` to replace claim-by-`qr_code`
  with claim-by-`claim_code`, including one-time use, expiration, and `403`
  semantics for invalid or expired codes.
- [ ] T021 Add a new migration to extend `establishments` with
  `claim_code_hash` and `claim_code_expires_at`.
- [ ] T022 Extend the establishment domain model and repository contract to use
  `claimCode` instead of `claimQRCode`.
- [ ] T023 Add failing unit tests for `ClaimEstablishmentOwnershipUseCase`
  covering missing `claim_code`, invalid code, expired code, not found, already
  claimed, and success.
- [ ] T024 Implement `ClaimEstablishmentOwnershipUseCase` until the tests pass.
- [ ] T025 Add a small claim-code hashing/verifier utility using a dedicated
  application secret.
- [ ] T026 Extend application config to load a `CLAIM_CODE_SECRET`-style value
  for claim-code verification.
- [ ] T027 Extend the Postgres establishment DTO mapping with
  `claim_code_hash` and `claim_code_expires_at`.
- [ ] T028 Add failing repository tests or verification coverage for claim-code
  ownership transfer behavior, including invalid, expired, consumed, and
  already-claimed paths.
- [ ] T029 Implement transactional claim ownership in the Postgres repository,
  including hash comparison, expiration checks, and one-time invalidation of the
  claim code after success.
- [ ] T030 Refactor the inbound claim request DTO from `qr_code` to
  `claim_code`.
- [ ] T031 Add failing controller/HTTP tests for the new claim-code request
  contract and the `400`, `403`, `404`, and `409` response paths.
- [ ] T032 Implement the claim handler changes and route wiring needed to
  support the new request contract.
- [ ] T033 Run targeted and full Go test suites for the claim refactor:
  use cases, controller layer, and `./internal/... ./cmd/...`.
- [ ] T034 Refactor, clean up fixtures, and mark the claim-security tasks as
  complete once the branch is stable.

## Acceptance Scenario Mapping

- T006, T010, T014 support scenarios 1, 2, and 3.
- T007, T014 support scenario 4.
- T008, T011, T014 support scenario 5.
- T020 through T034 harden the ownership-claim path that supports operator
  authorization for scenarios 1 through 4.

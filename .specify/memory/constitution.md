# FiveStars Constitution

## Purpose

This repository uses Specification-Driven Development (SDD) to turn product intent
into implementation through stable artifacts. Specifications are the primary
source of truth. Plans and tasks must trace back to them.

## Product Context

FiveStars is a credibility-first reviews platform. Reviews exist only after
verified check-ins, and establishment-facing capabilities must preserve that
trust model.

## Core Principles

### 1. Domain Integrity First

Business rules must live in domain models and use cases before reaching
controllers or repositories. Persistence cannot be the first line of defense for
invalid state.

### 2. Contract-Driven HTTP APIs

Every externally visible behavior must be expressed through explicit request and
response contracts. New endpoints require documented payloads, validation rules,
and error cases.

### 3. DTO Boundaries in Persistence

Postgres adapters must translate domain entities through DTOs when persisting or
rehydrating state. Repository code should avoid hidden mutation unless it is
required by the contract, such as reading generated values from the database.

### 4. Incremental MVP Delivery

Features should be sliced so they provide useful behavior as early as possible.
Prefer a minimal vertical slice that can be tested end-to-end before adding
secondary enhancements.

### 5. Testable by Construction

Each feature spec must define acceptance scenarios before implementation. Use
cases and repositories should be shaped so they can be validated with automated
tests at the right level.

### 6. TDD by Default

Task implementation should follow Test-Driven Development by default. New
behavior should start with a failing test that expresses the intended outcome,
followed by the minimal implementation to make it pass, and then refactoring as
needed. If a task cannot reasonably start with a test, that exception should be
made explicit in the plan or task notes.

## Delivery Rules

- `spec.md` captures WHAT and WHY.
- `plan.md` captures HOW.
- `tasks.md` captures execution order and ownership of work.
- Implementation tasks should prefer a red-green-refactor flow.
- Migrations must be explicit and forward-only within a feature branch.
- New behavior in existing flows must preserve current contracts unless the spec
  explicitly documents a breaking change.

## Current Architecture Constraints

- Backend: Go
- Persistence: PostgreSQL
- Repository layout: `internal/infra/adapters/outbound/repository/postgres`
- HTTP layer: `internal/infra/adapters/inbound/controller`
- Use cases: `internal/application/usecases`
- Migrations: `migrations/`

## Feature-Specific Guidance for Phase 5

- Highlights must only reference reviews that belong to the same establishment.
- Establishment dashboard work should start with a minimal stats surface instead
  of a full management console.
- Authorization for establishment operators must be planned explicitly, even if
  implementation starts with a temporary simplification.

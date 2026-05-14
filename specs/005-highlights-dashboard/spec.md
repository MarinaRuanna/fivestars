# Feature Specification: Highlights and Establishment Dashboard

Feature Branch: `005-highlights-dashboard`  
Created: 2026-05-13  
Status: Draft

## Objective

Allow an establishment operator to select a small number of standout reviews and
expose a minimal dashboard summary so the establishment page can promote trusted
feedback without undermining the product's credibility model.

## Problem Statement

The current system supports reviews and likes, but establishments cannot surface
their best reviews or view a simple summary of review activity. This limits the
value of reviews for establishment-facing experiences and leaves Phase 5 of the
execution plan without a delivery path.

## User Stories

- As an establishment operator, I want to highlight selected reviews, so that my
  establishment page can feature the most useful social proof.
- As an app or web client, I want establishment details to include highlights and
  simple stats, so that I can render a richer establishment experience without
  assembling data from many endpoints.
- As a product team, I want highlight selection to obey strict rules, so that
  establishments cannot misrepresent reviews from other places or flood the page
  with too many featured items.

## Scope

### In Scope

- Add a highlight relationship between establishments and reviews.
- Allow an authorized establishment operator to create and remove highlights.
- Return highlights as part of establishment detail payloads.
- Expose a minimal stats endpoint for establishment review metrics.
- Enforce a maximum highlight limit per establishment.

### Out of Scope

- Full dashboard web UI.
- Establishment replies to reviews.
- Review moderation or denunciation workflow.
- Wallet rewards tied to highlights.

## Business Rules

- A highlight can reference only a review that belongs to the same establishment.
- A highlight can be created only by an authorized establishment operator.
- An establishment may have at most 3 active highlights.
- The same review cannot be highlighted twice for the same establishment.
- Removing a highlight must not alter or delete the underlying review.
- Highlights returned in establishment detail should be ordered by newest
  highlight first.
- Stats must reflect all reviews for the establishment, not only highlighted
  reviews.

## Authorization Decision

- Highlight management depends on an explicit authorization policy interface in
  the use case layer.
- The first concrete rule is ownership-based: only the authenticated
  `owner_id` of the establishment may create or remove highlights.
- To support preexisting establishments created before `owner_id` existed, the
  API exposes an ownership-claim endpoint so the operator path does not depend
  on manual database backfill.
- This keeps the endpoint contracts stable while preserving room to evolve
  later to `establishment_users`.

## API Contract

### Create Highlight

- Method: `POST`
- Path: `/establishments/:id/highlights`
- Auth: `establishment operator`

Request:

```json
{
  "review_id": "uuid"
}
```

Success response:

```json
{
  "establishment_id": "uuid",
  "review_id": "uuid",
  "created_by_user_id": "uuid",
  "created_at": "2026-05-13T12:00:00Z"
}
```

Error cases:

- `400` when `review_id` is missing or invalid.
- `401` when the caller is not authenticated.
- `403` when the caller is not allowed to manage the establishment.
- `404` when the establishment or review does not exist.
- `409` when the review already is highlighted or the limit is exceeded.

### Delete Highlight

- Method: `DELETE`
- Path: `/establishments/:id/highlights/:reviewId`
- Auth: `establishment operator`

Success response:

- `204 No Content`

Error cases:

- `401` when the caller is not authenticated.
- `403` when the caller is not allowed to manage the establishment.
- `404` when the establishment or highlight does not exist.

### Claim Establishment Ownership

- Method: `POST`
- Path: `/establishments/:id/claim`
- Auth: `authenticated user`

Request:

```json
{
  "qr_code": "secret"
}
```

Success response:

```json
{
  "establishment_id": "uuid",
  "owner_id": "uuid"
}
```

Error cases:

- `400` when `qr_code` is missing.
- `401` when the caller is not authenticated.
- `403` when the provided `qr_code` does not prove possession of the establishment.
- `404` when the establishment does not exist.
- `409` when the establishment is already claimed by another user.

### Get Establishment Stats

- Method: `GET`
- Path: `/establishments/:id/stats`
- Auth: `public` for now, unless the team decides to restrict it later

Success response:

```json
{
  "average_rating": 4.6,
  "total_reviews": 28,
  "total_likes": 91,
  "highlighted_reviews_count": 3
}
```

Error cases:

- `404` when the establishment does not exist.

### Get Establishment Detail

- Existing `GET /establishments/:id` response must include a `highlights` field
  with up to 3 highlighted reviews enriched with the same review payload shape
  already exposed elsewhere.

## Domain Impact

- Entities:
  - add `Highlight`
  - add `EstablishmentStats`
- Use cases:
  - create highlight
  - delete highlight
  - get establishment stats
  - extend establishment detail retrieval with highlights
- Repositories:
  - highlight repository
  - establishment repository stats/detail extension
  - possible review repository query reuse
- Controllers:
  - highlight management endpoints
  - stats endpoint
- Migrations:
  - create `highlights`
  - optional support indexes

## Acceptance Scenarios

1. Given an authorized operator and a review from the same establishment, when the
   operator creates a highlight, then the API persists it and the establishment
   detail response includes it in `highlights`.
2. Given a review from another establishment, when an operator attempts to
   highlight it, then the API rejects the request with a validation or conflict
   style error and no row is created.
3. Given an establishment that already has 3 highlights, when an operator tries
   to create a fourth one, then the API rejects the request with `409`.
4. Given an existing highlight, when an authorized operator removes it, then the
   highlight disappears from the establishment detail response and the review
   remains unchanged.
5. Given an establishment with reviews and likes, when a client fetches
   `/establishments/:id/stats`, then the response returns average rating, total
   reviews, total likes, and highlighted review count.

## Open Questions

- Should stats remain public, or should they require operator auth from day one?
- Should highlighted reviews be limited to positive ratings only, or can any
  valid review be highlighted in MVP?

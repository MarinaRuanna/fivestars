# Data Model: Highlights and Establishment Dashboard

## Overview

Phase 5 introduces a new highlight relationship owned by an establishment and a
small aggregate projection for dashboard stats.

The implementation should preserve the existing review model and avoid copying
review content into a second table. Highlights point to reviews; review data
continues to live in `reviews`.

## Entities

### Highlight

Represents a featured review selected by an establishment operator.

Fields:

- `establishment_id` (`uuid`, required)
- `review_id` (`uuid`, required)
- `created_by_user_id` (`uuid`, required)
- `created_at` (`timestamptz`, required)

Rules:

- `review_id` must reference an existing review.
- The referenced review must belong to `establishment_id`.
- Maximum of 3 highlights per establishment.
- Duplicate highlight pairs are forbidden.

Suggested constraints:

- primary key: optional surrogate `id`, or use composite key
- unique: `(establishment_id, review_id)`
- foreign key: `establishment_id -> establishments.id`
- foreign key: `review_id -> reviews.id`
- foreign key: `created_by_user_id -> users.id`

Suggested indexes:

- index on `(establishment_id, created_at desc)`
- unique index on `(establishment_id, review_id)`

### EstablishmentStats

Read model returned by the stats endpoint. This does not need a table in MVP.

Fields:

- `average_rating` (`numeric` or `float64`)
- `total_reviews` (`int`)
- `total_likes` (`int`)
- `highlighted_reviews_count` (`int`)

Source:

- computed on read from `reviews`, `review_likes`, and `highlights`

## Query Shapes

### Highlighted Reviews for Establishment Detail

Recommended query shape:

1. fetch highlights for the establishment ordered by `created_at desc`
2. join the associated reviews
3. enrich with existing review fields, including like count

Expected output:

- array of up to 3 highlighted reviews
- review payload shape should match existing review API shape

### Establishment Stats

Recommended aggregates:

- `average_rating`: `AVG(reviews.rating)`
- `total_reviews`: `COUNT(reviews.id)`
- `total_likes`: count of `review_likes` joined through reviews of the
  establishment
- `highlighted_reviews_count`: `COUNT(highlights.review_id)`

Null-handling:

- average should default to `0` if there are no reviews
- counts should default to `0`

## Authorization Placeholder

Phase 5 requires operator authorization, but the ownership model is still open.

Preferred long-term model:

- `establishment_users`
  - `establishment_id`
  - `user_id`
  - `role`
  - `created_at`

For MVP planning, keep the highlight use cases behind an explicit authorization
check so the implementation can start with a temporary rule and migrate later.

## Repository Impact

Likely additions:

- `HighlightRepository.Create`
- `HighlightRepository.Delete`
- `HighlightRepository.CountByEstablishment`
- `HighlightRepository.ListByEstablishment`
- `EstablishmentRepository.GetStats`

Potential reuse:

- existing review repository logic for mapping review rows and like counts

## Migration Notes

Migration should:

1. create `highlights`
2. add constraints and indexes
3. avoid backfill requirements for MVP

## Open Technical Decisions

- composite primary key vs surrogate `id`
- whether `GET /establishments/:id` should hydrate highlights inside the
  establishment repository or compose data in a use case
- whether stats should be computed through one SQL query or multiple simpler
  queries in MVP

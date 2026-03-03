-- Enforce one check-in per user+establishment per UTC day at database level.
CREATE UNIQUE INDEX IF NOT EXISTS ux_checkins_user_estab_day_utc
ON checkins (user_id, establishment_id, ((date_trunc('day', checked_at AT TIME ZONE 'utc'))));

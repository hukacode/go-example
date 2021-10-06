CREATE TABLE IF NOT EXISTS task (
  id bigserial PRIMARY KEY,
  created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
  content text NOT NULL,
  is_completed boolean NOT NULL DEFAULT false
)
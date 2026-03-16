CREATE TABLE capsules (
  id SERIAL PRIMARY KEY,
  owner_id TEXT NOT NULL,
  recipient_email TEXT NOT NULL DEFAULT '',
  message TEXT NOT NULL DEFAULT '',
  media_url TEXT NOT NULL DEFAULT '',
  unlock_at TIMESTAMPTZ NOT NULL,
  delivered_at TIMESTAMPTZ,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  deleted_at TIMESTAMPTZ
);

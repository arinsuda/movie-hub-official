package migrate

const WatchLogVisibilityV1 = `
-- Step 1
CREATE TABLE IF NOT EXISTS media_watch_logs (
  id BIGSERIAL PRIMARY KEY,
  user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  media_type VARCHAR(10) NOT NULL CHECK (media_type IN ('movie', 'tv')),
  media_id INTEGER NOT NULL CHECK (media_id > 0),
  watched_on DATE NOT NULL,
  visibility VARCHAR(20) NOT NULL DEFAULT 'public' CHECK (visibility IN ('public', 'followers', 'private')),
  created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX IF NOT EXISTS idx_watch_logs_user_media ON media_watch_logs(user_id, media_type, media_id, watched_on DESC) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_watch_logs_user_date ON media_watch_logs(user_id, watched_on DESC, id DESC) WHERE deleted_at IS NULL;

-- Step 2
ALTER TABLE reviews ADD COLUMN IF NOT EXISTS visibility VARCHAR(20);
UPDATE reviews SET visibility = CASE WHEN is_public = true THEN 'public' ELSE 'private' END WHERE visibility IS NULL;
ALTER TABLE reviews ALTER COLUMN visibility SET NOT NULL;
ALTER TABLE reviews ALTER COLUMN visibility SET DEFAULT 'public';
ALTER TABLE reviews DROP CONSTRAINT IF EXISTS chk_reviews_visibility;
ALTER TABLE reviews ADD CONSTRAINT chk_reviews_visibility CHECK (visibility IN ('public', 'followers', 'private'));

-- Step 3
ALTER TABLE activity_events ADD COLUMN IF NOT EXISTS watch_log_id BIGINT;
ALTER TABLE activity_events DROP CONSTRAINT IF EXISTS fk_activity_events_watch_log;
ALTER TABLE activity_events ADD CONSTRAINT fk_activity_events_watch_log
  FOREIGN KEY (watch_log_id) REFERENCES media_watch_logs(id) ON DELETE SET NULL;

-- Step 4
INSERT INTO media_watch_logs (user_id, media_type, media_id, watched_on, visibility, created_at, updated_at)
SELECT user_id, media_type, media_id,
  COALESCE(watched_at::date, created_at::date),
  'public',
  created_at, updated_at
FROM library_items
WHERE list_type = 'watched' AND deleted_at IS NULL
ON CONFLICT DO NOTHING;
`

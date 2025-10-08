-- ============================
-- 1. Drop triggers
-- ============================

DROP TRIGGER IF EXISTS update_events_updated_at ON events;
DROP TRIGGER IF EXISTS update_music_sheets_updated_at ON music_sheets;
DROP TRIGGER IF EXISTS update_musics_updated_at ON musics;
DROP TRIGGER IF EXISTS update_settings_updated_at ON settings;
DROP TRIGGER IF EXISTS update_users_updated_at ON users;

-- ============================
-- 2. Drop tables in reverse dependency order
-- ============================

-- Drop junction table first (depends on events and musics)
DROP TABLE IF EXISTS event_musics CASCADE;

-- Drop events (depends on users)
DROP TABLE IF EXISTS events CASCADE;

-- Drop music_sheets (depends on musics)
DROP TABLE IF EXISTS music_sheets CASCADE;

-- Drop musics (depends on users)
DROP TABLE IF EXISTS musics CASCADE;

-- Drop base tables with no dependencies
DROP TABLE IF EXISTS settings CASCADE;
DROP TABLE IF EXISTS users CASCADE;

-- ============================
-- 3. Drop indexes (optional, mostly handled by DROP TABLE)
-- ============================

-- Junction table indexes
DROP INDEX IF EXISTS idx_event_musics_event_id;
DROP INDEX IF EXISTS idx_event_musics_music_id;

-- Events indexes
DROP INDEX IF EXISTS idx_events_user_id;
DROP INDEX IF EXISTS idx_events_start_time;

-- Music sheets indexes
DROP INDEX IF EXISTS idx_music_sheets_title;
DROP INDEX IF EXISTS idx_music_sheets_music_id;

-- Musics indexes
DROP INDEX IF EXISTS idx_musics_user_id;

-- Users indexes
DROP INDEX IF EXISTS idx_users_email;

-- ============================
-- 4. Drop shared function
-- ============================

DROP FUNCTION IF EXISTS update_updated_at_column();
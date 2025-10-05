-- ============================
-- 1. Drop triggers
-- ============================

DROP TRIGGER IF EXISTS update_events_updated_at ON events;
DROP TRIGGER IF EXISTS update_musics_updated_at ON musics;
DROP TRIGGER IF EXISTS update_music_sheets_updated_at ON music_sheets;
DROP TRIGGER IF EXISTS update_users_updated_at ON users;
DROP TRIGGER IF EXISTS update_settings_updated_at ON settings;

-- ============================
-- 2. Drop tables in reverse dependency order
-- ============================

DROP TABLE IF EXISTS events CASCADE;
DROP TABLE IF EXISTS musics CASCADE;
DROP TABLE IF EXISTS music_sheets CASCADE;
DROP TABLE IF EXISTS users CASCADE;
DROP TABLE IF EXISTS settings CASCADE;

-- ============================
-- 3. Drop indexes (optional, mostly handled by DROP TABLE)
-- ============================

DROP INDEX IF EXISTS idx_events_user_id;
DROP INDEX IF EXISTS idx_events_start_time;

DROP INDEX IF EXISTS idx_musics_user_id;
DROP INDEX IF EXISTS idx_musics_sheet_id;

DROP INDEX IF EXISTS idx_music_sheets_title;

DROP INDEX IF EXISTS idx_users_email;

-- ============================
-- 4. Drop shared function
-- ============================

DROP FUNCTION IF EXISTS update_updated_at_column();

-- ============================
-- 1. Drop triggers
-- ============================

DROP TRIGGER IF EXISTS update_bookings_updated_at ON bookings;
DROP TRIGGER IF EXISTS update_events_updated_at ON events;
DROP TRIGGER IF EXISTS update_music_sheets_updated_at ON music_sheets;
DROP TRIGGER IF EXISTS update_musics_updated_at ON musics;
DROP TRIGGER IF EXISTS update_settings_updated_at ON settings;
DROP TRIGGER IF EXISTS update_tokens_updated_at ON tokens;
DROP TRIGGER IF EXISTS update_users_updated_at ON users;
DROP TRIGGER IF EXISTS update_roles_updated_at ON roles;

-- ============================
-- 2. Drop tables in reverse dependency order
-- ============================

-- Drop bookings (depends on events and users)
DROP TABLE IF EXISTS bookings CASCADE;

-- Drop junction table first (depends on events and musics)
DROP TABLE IF EXISTS event_musics CASCADE;

-- Drop events (depends on users)
DROP TABLE IF EXISTS events CASCADE;

-- Drop music_sheets (depends on musics)
DROP TABLE IF EXISTS music_sheets CASCADE;

-- Drop musics (depends on users)
DROP TABLE IF EXISTS musics CASCADE;

-- Drop tokens (depends on users)
DROP TABLE IF EXISTS tokens CASCADE;

-- Drop users (depends on roles)
DROP TABLE IF EXISTS users CASCADE;

-- Drop base tables with no dependencies
DROP TABLE IF EXISTS settings CASCADE;
DROP TABLE IF EXISTS roles CASCADE;

-- ============================
-- 3. Drop indexes (optional, mostly handled by DROP TABLE)
-- ============================

-- Bookings indexes
DROP INDEX IF EXISTS idx_bookings_event_id;
DROP INDEX IF EXISTS idx_bookings_user_id;
DROP INDEX IF EXISTS idx_bookings_status;

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

-- Tokens indexes
DROP INDEX IF EXISTS idx_tokens_user_id;
DROP INDEX IF EXISTS idx_tokens_token;
DROP INDEX IF EXISTS idx_tokens_is_active;
DROP INDEX IF EXISTS idx_tokens_expires_at;

-- Users indexes
DROP INDEX IF EXISTS idx_users_email;
DROP INDEX IF EXISTS idx_users_role_id;

-- Roles indexes
DROP INDEX IF EXISTS idx_roles_name;

-- ============================
-- 4. Drop shared function
-- ============================

DROP FUNCTION IF EXISTS update_updated_at_column();
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

-- Drop favorites (depends on users and musics)
DROP TABLE IF EXISTS favorites CASCADE;

-- Drop band_musics junction table (depends on bands and musics)
DROP TABLE IF EXISTS band_musics CASCADE;

-- Drop bands (depends on users)
DROP TABLE IF EXISTS bands CASCADE;

-- Drop bookings (depends on events and users)
DROP TABLE IF EXISTS bookings CASCADE;

-- Drop junction table first (depends on events and musics)
DROP TABLE IF EXISTS event_musics CASCADE;

-- Drop events (depends on users)
DROP TABLE IF EXISTS events CASCADE;

-- Drop music_sheets (depends on musics)
DROP TABLE IF EXISTS music_sheets CASCADE;

-- Drop music_audio (depends on musics)
DROP TABLE IF EXISTS music_audio CASCADE;

-- Drop musics (depends on users)
DROP TABLE IF EXISTS musics CASCADE;

-- Drop tokens (depends on users)
DROP TABLE IF EXISTS tokens CASCADE;

-- Remove foreign key constraint from users before dropping bands
ALTER TABLE users DROP CONSTRAINT IF EXISTS fk_users_band_id;

-- Drop users (depends on roles)
DROP TABLE IF EXISTS users CASCADE;

-- Drop base tables with no dependencies
DROP TABLE IF EXISTS settings CASCADE;
DROP TABLE IF EXISTS roles CASCADE;

-- ============================
-- 3. Drop indexes (optional, mostly handled by DROP TABLE)
-- ============================

-- Favorites indexes
DROP INDEX IF EXISTS idx_favorites_user_id;
DROP INDEX IF EXISTS idx_favorites_music_id;
DROP INDEX IF EXISTS idx_favorites_created_at;

-- Band musics indexes
DROP INDEX IF EXISTS idx_band_musics_band_id;
DROP INDEX IF EXISTS idx_band_musics_music_id;

-- Bands indexes
DROP INDEX IF EXISTS idx_bands_user_id;
DROP INDEX IF EXISTS idx_bands_name;
DROP INDEX IF EXISTS idx_bands_is_public;

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
DROP INDEX IF EXISTS idx_music_sheets_music_id;
DROP INDEX IF EXISTS idx_music_sheets_type;
DROP INDEX IF EXISTS idx_music_sheets_lang;
DROP INDEX IF EXISTS idx_music_sheets_difficulty;

-- Music audio indexes
DROP INDEX IF EXISTS idx_music_audio_music_id;
DROP INDEX IF EXISTS idx_music_audio_file_type;
DROP INDEX IF EXISTS idx_music_audio_is_primary;

-- Musics indexes
DROP INDEX IF EXISTS idx_musics_user_id;
DROP INDEX IF EXISTS idx_musics_title;
DROP INDEX IF EXISTS idx_musics_artist;
DROP INDEX IF EXISTS idx_musics_genre;

-- Tokens indexes
DROP INDEX IF EXISTS idx_tokens_user_id;
DROP INDEX IF EXISTS idx_tokens_token;
DROP INDEX IF EXISTS idx_tokens_is_active;
DROP INDEX IF EXISTS idx_tokens_expires_at;

-- Users indexes
DROP INDEX IF EXISTS idx_users_email;
DROP INDEX IF EXISTS idx_users_provider_id;
DROP INDEX IF EXISTS idx_users_band_id;

-- Roles indexes
DROP INDEX IF EXISTS idx_roles_name;

-- ============================
-- 4. Drop shared function
-- ============================

DROP FUNCTION IF EXISTS update_updated_at_column();
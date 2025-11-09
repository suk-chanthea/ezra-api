-- ============================
-- 1. Drop triggers
-- ============================

DROP TRIGGER IF EXISTS update_supporters_updated_at ON supporters;
DROP TRIGGER IF EXISTS update_donations_updated_at ON donations;
DROP TRIGGER IF EXISTS update_device_tokens_updated_at ON device_tokens;
DROP TRIGGER IF EXISTS update_notifications_updated_at ON notifications;
DROP TRIGGER IF EXISTS update_bookings_updated_at ON bookings;
DROP TRIGGER IF EXISTS update_events_updated_at ON events;
DROP TRIGGER IF EXISTS update_music_sheets_updated_at ON music_sheets;
DROP TRIGGER IF EXISTS update_music_audio_updated_at ON music_audio;
DROP TRIGGER IF EXISTS update_musics_updated_at ON musics;
DROP TRIGGER IF EXISTS update_settings_updated_at ON settings;
DROP TRIGGER IF EXISTS update_otps_updated_at ON otps;
DROP TRIGGER IF EXISTS update_tokens_updated_at ON tokens;
DROP TRIGGER IF EXISTS update_users_updated_at ON users;
DROP TRIGGER IF EXISTS create_settings_for_new_user ON users;
DROP TRIGGER IF EXISTS update_roles_updated_at ON roles;
DROP TRIGGER IF EXISTS update_bands_updated_at ON bands;
DROP TRIGGER IF EXISTS update_churches_updated_at ON churches;

-- ============================
-- 2. Drop foreign keys and columns from users (only if table exists)
-- ============================

-- Drop the constraints only if the users table exists
DO $$ 
BEGIN
    IF EXISTS (SELECT FROM pg_tables WHERE schemaname = 'public' AND tablename = 'users') THEN
        ALTER TABLE users DROP CONSTRAINT IF EXISTS fk_users_band_id;
        ALTER TABLE users DROP CONSTRAINT IF EXISTS fk_users_church_id;
    END IF;
END $$;

-- ============================
-- 3. Drop tables in reverse dependency order
-- ============================

-- Drop donations (depends on users, events, and supporters)
DROP TABLE IF EXISTS donations CASCADE;

-- Drop supporters (depends on users)
DROP TABLE IF EXISTS supporters CASCADE;

-- Drop churches (no dependencies, but users reference it)
DROP TABLE IF EXISTS churches CASCADE;

-- Drop device tokens (depends on users)
DROP TABLE IF EXISTS device_tokens CASCADE;

-- Drop notifications (depends on users)
DROP TABLE IF EXISTS notifications CASCADE;

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

-- Drop otps (no dependencies)
DROP TABLE IF EXISTS otps CASCADE;

-- Drop settings (depends on users)
DROP TABLE IF EXISTS settings CASCADE;

-- Drop users (depends on roles)
DROP TABLE IF EXISTS users CASCADE;

-- Drop base tables with no dependencies
DROP TABLE IF EXISTS roles CASCADE;

-- ============================
-- 4. Drop indexes (optional, mostly handled by DROP TABLE)
-- ============================

-- Supporters indexes
DROP INDEX IF EXISTS idx_supporters_email;
DROP INDEX IF EXISTS idx_supporters_user_id;
DROP INDEX IF EXISTS idx_supporters_type;
DROP INDEX IF EXISTS idx_supporters_name;

-- Donations indexes
DROP INDEX IF EXISTS idx_donations_type;
DROP INDEX IF EXISTS idx_donations_donor_type;
DROP INDEX IF EXISTS idx_donations_user_id;
DROP INDEX IF EXISTS idx_donations_supporter_id;
DROP INDEX IF EXISTS idx_donations_event_id;
DROP INDEX IF EXISTS idx_donations_status;
DROP INDEX IF EXISTS idx_donations_created_at;
DROP INDEX IF EXISTS idx_donations_transaction_id;
DROP INDEX IF EXISTS idx_donations_qr_expires_at;

-- Notifications indexes
DROP INDEX IF EXISTS idx_notifications_user_id;
DROP INDEX IF EXISTS idx_notifications_band_id;
DROP INDEX IF EXISTS idx_notifications_sender_id;
DROP INDEX IF EXISTS idx_notifications_recipient_type;
DROP INDEX IF EXISTS idx_notifications_is_read;
DROP INDEX IF EXISTS idx_notifications_created_at;
DROP INDEX IF EXISTS idx_notifications_type;

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

-- OTPs indexes
DROP INDEX IF EXISTS idx_otps_email;
DROP INDEX IF EXISTS idx_otps_purpose;
DROP INDEX IF EXISTS idx_otps_expires_at;
DROP INDEX IF EXISTS idx_otps_email_purpose;

-- Settings indexes
DROP INDEX IF EXISTS idx_settings_user_id;

-- Users indexes
DROP INDEX IF EXISTS idx_users_email;
DROP INDEX IF EXISTS idx_users_email_verified;
DROP INDEX IF EXISTS idx_users_phone;
DROP INDEX IF EXISTS idx_users_provider_id;
DROP INDEX IF EXISTS idx_users_band_id;
DROP INDEX IF EXISTS idx_users_church_id;
DROP INDEX IF EXISTS idx_users_church_status;
DROP INDEX IF EXISTS idx_users_birthday;

-- Churches indexes
DROP INDEX IF EXISTS idx_churches_fullname;
DROP INDEX IF EXISTS idx_churches_email;
DROP INDEX IF EXISTS idx_churches_denomination;
DROP INDEX IF EXISTS idx_churches_owner_id;

-- Roles indexes
DROP INDEX IF EXISTS idx_roles_name;

-- ============================
-- 5. Drop shared functions
-- ============================

DROP FUNCTION IF EXISTS create_user_settings();
DROP FUNCTION IF EXISTS update_updated_at_column();
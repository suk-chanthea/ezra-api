-- ============================
-- 1. Roles table (NEW - for flexible role management)
-- ============================
CREATE TABLE IF NOT EXISTS roles (
    id SERIAL PRIMARY KEY,
    name VARCHAR(50) UNIQUE NOT NULL,
    description TEXT,
    permissions JSONB DEFAULT '[]'::jsonb,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT now(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_roles_name ON roles(name);

-- Insert default roles
INSERT INTO roles (name, description, permissions) VALUES
('admin', 'Full system access', '["user.read", "user.write", "user.delete", "music.read", "music.write", "music.delete", "event.read", "event.write", "event.delete", "role.manage"]'::jsonb),
('moderator', 'Can manage content', '["user.read", "music.read", "music.write", "music.delete", "event.read", "event.write", "event.delete"]'::jsonb),
('user', 'Regular user access', '["music.read", "music.write", "event.read"]'::jsonb),
('guest', 'Read-only access', '["music.read", "event.read"]'::jsonb)
ON CONFLICT (name) DO NOTHING;

-- ============================
-- 2. Users table (with role_id FK and OAuth support)
-- ============================
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(100) NOT NULL,
    fullname VARCHAR(100) NOT NULL,
    profile VARCHAR(255) NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    password VARCHAR(255),
    role VARCHAR(20) DEFAULT 'user',
    token VARCHAR(255),
    provider VARCHAR(50) DEFAULT 'local',
    provider_id VARCHAR(255),
    band_id INTEGER DEFAULT NULL,  -- User's affiliated band/organization
    created_at TIMESTAMP WITH TIME ZONE DEFAULT now(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
CREATE INDEX IF NOT EXISTS idx_users_provider_id ON users(provider, provider_id);
CREATE INDEX IF NOT EXISTS idx_users_band_id ON users(band_id);

-- ============================
-- 3. Tokens table (for multi-device/session support)
-- ============================
CREATE TABLE IF NOT EXISTS tokens (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token VARCHAR(500) NOT NULL UNIQUE,
    device_name VARCHAR(100),
    device_type VARCHAR(50), -- 'web', 'mobile', 'desktop', 'api'
    ip_address INET,
    user_agent TEXT,
    is_active BOOLEAN DEFAULT true,
    last_used_at TIMESTAMP WITH TIME ZONE DEFAULT now(),
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT now(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_tokens_user_id ON tokens(user_id);
CREATE INDEX IF NOT EXISTS idx_tokens_token ON tokens(token);
CREATE INDEX IF NOT EXISTS idx_tokens_is_active ON tokens(is_active);
CREATE INDEX IF NOT EXISTS idx_tokens_expires_at ON tokens(expires_at);

-- ============================
-- 4. Settings table
-- ============================
CREATE TABLE IF NOT EXISTS settings (
    id SERIAL NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT now(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT now()
);

-- ============================
-- 5. Musics table (Core music metadata)
-- ============================
CREATE TABLE IF NOT EXISTS musics (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    artist VARCHAR(255),
    album VARCHAR(255),
    genre VARCHAR(100),
    duration INTEGER,  -- in seconds
    bpm INTEGER,       -- beats per minute
    key VARCHAR(10),   -- musical key (C, Am, G, etc.)
    cover VARCHAR(255),
    lyrics TEXT,
    description TEXT,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ DEFAULT now(),
    updated_at TIMESTAMPTZ DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_musics_user_id ON musics(user_id);
CREATE INDEX IF NOT EXISTS idx_musics_title ON musics(title);
CREATE INDEX IF NOT EXISTS idx_musics_artist ON musics(artist);
CREATE INDEX IF NOT EXISTS idx_musics_genre ON musics(genre);

-- ============================
-- 6. Music Audio table (Multiple audio files per music)
-- ============================
CREATE TABLE IF NOT EXISTS music_audio (
    id SERIAL PRIMARY KEY,
    music_id INTEGER NOT NULL REFERENCES musics(id) ON DELETE CASCADE,
    file_path VARCHAR(500) NOT NULL,
    file_type VARCHAR(50) NOT NULL DEFAULT 'Original',
    format VARCHAR(10) NOT NULL DEFAULT 'mp3',
    quality VARCHAR(50),
    size_bytes BIGINT,
    duration INTEGER,  -- in seconds
    is_primary BOOLEAN DEFAULT false,
    created_at TIMESTAMPTZ DEFAULT now(),
    updated_at TIMESTAMPTZ DEFAULT now(),
    CHECK (file_type IN ('Original', 'Instrumental', 'Acapella', 'Live', 'Acoustic', 'Remix', 'Cover')),
    CHECK (format IN ('mp3', 'wav', 'flac', 'aac', 'm4a', 'ogg'))
);

CREATE INDEX IF NOT EXISTS idx_music_audio_music_id ON music_audio(music_id);
CREATE INDEX IF NOT EXISTS idx_music_audio_file_type ON music_audio(file_type);
CREATE INDEX IF NOT EXISTS idx_music_audio_is_primary ON music_audio(is_primary);

-- ============================
-- 7. Music Sheets table (Sheet music files)
-- ============================
CREATE TABLE IF NOT EXISTS music_sheets (
    id SERIAL PRIMARY KEY,
    music_id INTEGER NOT NULL REFERENCES musics(id) ON DELETE CASCADE,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    file_path VARCHAR(500) NOT NULL,
    type VARCHAR(50) NOT NULL DEFAULT 'Lead Sheet',
    lang VARCHAR(10) NOT NULL DEFAULT 'kh',
    key VARCHAR(10),   -- musical key
    difficulty VARCHAR(20),  -- beginner, intermediate, advanced
    created_at TIMESTAMPTZ DEFAULT now(),
    updated_at TIMESTAMPTZ DEFAULT now(),
    CHECK (type IN ('Standard Notation', 'Lead Sheet', 'Chord Chart', 'Tablature', 'PVG Sheet', 'Orchestral Score', 'Drum Notation')),
    CHECK (lang IN ('kh', 'en', 'kr', 'cn')),
    CHECK (difficulty IN ('beginner', 'intermediate', 'advanced'))
);

CREATE INDEX IF NOT EXISTS idx_music_sheets_music_id ON music_sheets(music_id);
CREATE INDEX IF NOT EXISTS idx_music_sheets_type ON music_sheets(type);
CREATE INDEX IF NOT EXISTS idx_music_sheets_lang ON music_sheets(lang);
CREATE INDEX IF NOT EXISTS idx_music_sheets_difficulty ON music_sheets(difficulty);

-- ============================
-- 8. Events table
-- ============================
CREATE TABLE IF NOT EXISTS events (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    content TEXT,
    cover VARCHAR(255),
    location TEXT NOT NULL,
    start_time TIMESTAMPTZ NOT NULL,
    end_time TIMESTAMPTZ NOT NULL,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ DEFAULT now(),
    updated_at TIMESTAMPTZ DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_events_user_id ON events(user_id);
CREATE INDEX IF NOT EXISTS idx_events_start_time ON events(start_time);

-- ============================
-- 9. Event_Musics Junction Table (Many-to-Many)
-- ============================
CREATE TABLE IF NOT EXISTS event_musics (
    id SERIAL PRIMARY KEY,
    event_id INTEGER NOT NULL REFERENCES events(id) ON DELETE CASCADE,
    music_id INTEGER NOT NULL REFERENCES musics(id) ON DELETE CASCADE,
    display_order INTEGER DEFAULT 0,
    created_at TIMESTAMPTZ DEFAULT now(),
    UNIQUE(event_id, music_id)
);

CREATE INDEX IF NOT EXISTS idx_event_musics_event_id ON event_musics(event_id);
CREATE INDEX IF NOT EXISTS idx_event_musics_music_id ON event_musics(music_id);

-- ============================
-- 10. Bookings table (Event Registrations)
-- ============================
CREATE TABLE IF NOT EXISTS bookings (
    id SERIAL PRIMARY KEY,
    event_id INTEGER NOT NULL REFERENCES events(id) ON DELETE CASCADE,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    status VARCHAR(50) NOT NULL DEFAULT 'pending',
    notes TEXT,
    created_at TIMESTAMPTZ DEFAULT now(),
    updated_at TIMESTAMPTZ DEFAULT now(),
    UNIQUE(event_id, user_id),
    CHECK (status IN ('pending', 'confirmed', 'cancelled'))
);

CREATE INDEX IF NOT EXISTS idx_bookings_event_id ON bookings(event_id);
CREATE INDEX IF NOT EXISTS idx_bookings_user_id ON bookings(user_id);
CREATE INDEX IF NOT EXISTS idx_bookings_status ON bookings(status);

-- ============================
-- 11. Bands table (Music Collections/Libraries)
-- ============================
CREATE TABLE IF NOT EXISTS bands (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    cover VARCHAR(255),
    is_public BOOLEAN DEFAULT false,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ DEFAULT now(),
    updated_at TIMESTAMPTZ DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_bands_user_id ON bands(user_id);
CREATE INDEX IF NOT EXISTS idx_bands_name ON bands(name);
CREATE INDEX IF NOT EXISTS idx_bands_is_public ON bands(is_public);

-- Add foreign key constraint to users.band_id after bands table is created
ALTER TABLE users
ADD CONSTRAINT fk_users_band_id 
FOREIGN KEY (band_id) REFERENCES bands(id) ON DELETE SET NULL;

-- ============================
-- 12. Band_Musics Junction Table (Many-to-Many)
-- ============================
CREATE TABLE IF NOT EXISTS band_musics (
    id SERIAL PRIMARY KEY,
    band_id INTEGER NOT NULL REFERENCES bands(id) ON DELETE CASCADE,
    music_id INTEGER NOT NULL REFERENCES musics(id) ON DELETE CASCADE,
    display_order INTEGER DEFAULT 0,
    created_at TIMESTAMPTZ DEFAULT now(),
    UNIQUE(band_id, music_id)
);

CREATE INDEX IF NOT EXISTS idx_band_musics_band_id ON band_musics(band_id);
CREATE INDEX IF NOT EXISTS idx_band_musics_music_id ON band_musics(music_id);

-- ============================
-- 13. Favorites table (User Favorite Music)
-- ============================
CREATE TABLE IF NOT EXISTS favorites (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    music_id INTEGER NOT NULL REFERENCES musics(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ DEFAULT now(),
    UNIQUE(user_id, music_id)
);

CREATE INDEX IF NOT EXISTS idx_favorites_user_id ON favorites(user_id);
CREATE INDEX IF NOT EXISTS idx_favorites_music_id ON favorites(music_id);
CREATE INDEX IF NOT EXISTS idx_favorites_created_at ON favorites(created_at);

-- ============================
-- 14. Shared auto-update function
-- ============================
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = now();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- ============================
-- 15. Triggers for auto-updating
-- ============================

-- roles table
DROP TRIGGER IF EXISTS update_roles_updated_at ON roles;
CREATE TRIGGER update_roles_updated_at
BEFORE UPDATE ON roles
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();

-- users table
DROP TRIGGER IF EXISTS update_users_updated_at ON users;
CREATE TRIGGER update_users_updated_at
BEFORE UPDATE ON users
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();

-- tokens table
DROP TRIGGER IF EXISTS update_tokens_updated_at ON tokens;
CREATE TRIGGER update_tokens_updated_at
BEFORE UPDATE ON tokens
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();

-- settings table
DROP TRIGGER IF EXISTS update_settings_updated_at ON settings;
CREATE TRIGGER update_settings_updated_at
BEFORE UPDATE ON settings
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();

-- musics table
DROP TRIGGER IF EXISTS update_musics_updated_at ON musics;
CREATE TRIGGER update_musics_updated_at
BEFORE UPDATE ON musics
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();

-- music_audio table
DROP TRIGGER IF EXISTS update_music_audio_updated_at ON music_audio;
CREATE TRIGGER update_music_audio_updated_at
BEFORE UPDATE ON music_audio
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();

-- music_sheets table
DROP TRIGGER IF EXISTS update_music_sheets_updated_at ON music_sheets;
CREATE TRIGGER update_music_sheets_updated_at
BEFORE UPDATE ON music_sheets
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();

-- events table
DROP TRIGGER IF EXISTS update_events_updated_at ON events;
CREATE TRIGGER update_events_updated_at
BEFORE UPDATE ON events
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();

-- bookings table
DROP TRIGGER IF EXISTS update_bookings_updated_at ON bookings;
CREATE TRIGGER update_bookings_updated_at
BEFORE UPDATE ON bookings
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();

-- bands table
DROP TRIGGER IF EXISTS update_bands_updated_at ON bands;
CREATE TRIGGER update_bands_updated_at
BEFORE UPDATE ON bands
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();
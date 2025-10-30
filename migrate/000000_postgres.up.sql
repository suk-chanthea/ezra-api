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
-- 2. Users table (WITHOUT band_id initially)
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
    created_at TIMESTAMP WITH TIME ZONE DEFAULT now(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
CREATE INDEX IF NOT EXISTS idx_users_provider_id ON users(provider, provider_id);

-- ============================
-- 3. Churches table
-- ============================
CREATE TABLE IF NOT EXISTS churches (
    id SERIAL PRIMARY KEY,
    fullname VARCHAR(255) NOT NULL UNIQUE,
    address TEXT,
    phone VARCHAR(50),
    email VARCHAR(255),
    website VARCHAR(255),
    pastor_name VARCHAR(255),
    description TEXT,
    logo VARCHAR(255),
    established_date DATE,
    denomination VARCHAR(100),
    owner_id INTEGER REFERENCES users(id) ON DELETE SET NULL,  -- Church owner/admin
    created_at TIMESTAMPTZ DEFAULT now(),
    updated_at TIMESTAMPTZ DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_churches_fullname ON churches(fullname);
CREATE INDEX IF NOT EXISTS idx_churches_email ON churches(email);
CREATE INDEX IF NOT EXISTS idx_churches_denomination ON churches(denomination);
CREATE INDEX IF NOT EXISTS idx_churches_owner_id ON churches(owner_id);

-- ============================
-- 4. Tokens table (for multi-device/session support)
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
-- 5. Settings table
-- ============================
CREATE TABLE IF NOT EXISTS settings (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    language VARCHAR(10) DEFAULT 'en',
    theme VARCHAR(20) DEFAULT 'light',
    notify_on_booking BOOLEAN DEFAULT true,
    notify_on_music BOOLEAN DEFAULT false,
    notify_on_event BOOLEAN DEFAULT true,
    enable_push_notifications BOOLEAN DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT now(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT now(),
    UNIQUE(user_id),
    CHECK (theme IN ('light', 'dark', 'auto'))
);

CREATE INDEX IF NOT EXISTS idx_settings_user_id ON settings(user_id);

-- ============================
-- 6. Musics table (Core music metadata)
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

-- ============================
-- 11b. NOW add band_id column to users table
-- ============================
ALTER TABLE users ADD COLUMN IF NOT EXISTS band_id INTEGER;

-- Add foreign key constraint
ALTER TABLE users
ADD CONSTRAINT fk_users_band_id 
FOREIGN KEY (band_id) REFERENCES bands(id) ON DELETE SET NULL;

-- Add index
CREATE INDEX IF NOT EXISTS idx_users_band_id ON users(band_id);

-- ============================
-- 11c. Add birthday, church_id, and bio to users table
-- ============================
ALTER TABLE users ADD COLUMN IF NOT EXISTS birthday DATE;
ALTER TABLE users ADD COLUMN IF NOT EXISTS church_id INTEGER;
ALTER TABLE users ADD COLUMN IF NOT EXISTS bio TEXT;
ALTER TABLE users ADD COLUMN IF NOT EXISTS church_status VARCHAR(20) DEFAULT 'pending';  -- 'pending', 'approved', 'rejected'

-- Add foreign key constraint for church_id
ALTER TABLE users
ADD CONSTRAINT fk_users_church_id 
FOREIGN KEY (church_id) REFERENCES churches(id) ON DELETE SET NULL;

-- Add check constraint for church_status
ALTER TABLE users ADD CONSTRAINT users_church_status_check CHECK (church_status IN ('pending', 'approved', 'rejected'));

CREATE INDEX IF NOT EXISTS idx_users_church_id ON users(church_id);
CREATE INDEX IF NOT EXISTS idx_users_birthday ON users(birthday);
CREATE INDEX IF NOT EXISTS idx_users_church_status ON users(church_status);

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
-- 14b. Function to create default settings for new users
-- ============================
CREATE OR REPLACE FUNCTION create_user_settings()
RETURNS TRIGGER AS $$
BEGIN
    INSERT INTO settings (user_id, language, theme, notify_on_booking, notify_on_music, notify_on_event, enable_push_notifications)
    VALUES (NEW.id, 'en', 'light', true, false, true, true);
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

-- Trigger to create default settings for new users
DROP TRIGGER IF EXISTS create_settings_for_new_user ON users;
CREATE TRIGGER create_settings_for_new_user
AFTER INSERT ON users
FOR EACH ROW
EXECUTE FUNCTION create_user_settings();

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

-- churches table
DROP TRIGGER IF EXISTS update_churches_updated_at ON churches;
CREATE TRIGGER update_churches_updated_at
BEFORE UPDATE ON churches
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();

-- ============================
-- 16. Notifications table
-- ============================
CREATE TABLE IF NOT EXISTS notifications (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,  -- NULL for broadcast
    sender_id INTEGER REFERENCES users(id) ON DELETE SET NULL,  -- Who sent it
    band_id INTEGER REFERENCES bands(id) ON DELETE CASCADE,  -- For team notifications
    recipient_type VARCHAR(20) NOT NULL DEFAULT 'user',  -- 'user', 'band', 'all'
    title VARCHAR(255) NOT NULL,
    message TEXT NOT NULL,
    type VARCHAR(50) NOT NULL DEFAULT 'info',
    related_type VARCHAR(50),  -- 'music', 'event', 'booking', 'band'
    related_id INTEGER,
    is_read BOOLEAN DEFAULT false,
    read_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ DEFAULT now(),
    updated_at TIMESTAMPTZ DEFAULT now(),
    CHECK (type IN ('info', 'success', 'warning', 'error', 'booking', 'music', 'event')),
    CHECK (recipient_type IN ('user', 'band', 'all')),
    CHECK (
        (recipient_type = 'user' AND user_id IS NOT NULL AND band_id IS NULL) OR
        (recipient_type = 'band' AND band_id IS NOT NULL AND user_id IS NULL) OR
        (recipient_type = 'all' AND user_id IS NULL AND band_id IS NULL)
    )
);

-- Device tokens table for FCM (Firebase Cloud Messaging)
-- Stores device tokens for push notifications to mobile and web
CREATE TABLE IF NOT EXISTS device_tokens (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token TEXT NOT NULL,
    platform VARCHAR(20) NOT NULL,  -- 'ios', 'android', 'web'
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMPTZ DEFAULT now(),
    updated_at TIMESTAMPTZ DEFAULT now(),
    CONSTRAINT unique_user_token UNIQUE(user_id, token)
);

CREATE INDEX IF NOT EXISTS idx_notifications_user_id ON notifications(user_id);
CREATE INDEX IF NOT EXISTS idx_notifications_band_id ON notifications(band_id);
CREATE INDEX IF NOT EXISTS idx_notifications_sender_id ON notifications(sender_id);
CREATE INDEX IF NOT EXISTS idx_notifications_recipient_type ON notifications(recipient_type);
CREATE INDEX IF NOT EXISTS idx_notifications_is_read ON notifications(is_read);
CREATE INDEX IF NOT EXISTS idx_notifications_created_at ON notifications(created_at DESC);
CREATE INDEX IF NOT EXISTS idx_notifications_type ON notifications(type);

CREATE INDEX IF NOT EXISTS idx_device_tokens_user_id ON device_tokens(user_id);
CREATE INDEX IF NOT EXISTS idx_device_tokens_is_active ON device_tokens(is_active);
CREATE INDEX IF NOT EXISTS idx_device_tokens_token ON device_tokens(token);

-- notifications table trigger
DROP TRIGGER IF EXISTS update_notifications_updated_at ON notifications;
CREATE TRIGGER update_notifications_updated_at
BEFORE UPDATE ON notifications
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();

DROP TRIGGER IF EXISTS update_device_tokens_updated_at ON device_tokens;
CREATE TRIGGER update_device_tokens_updated_at
BEFORE UPDATE ON device_tokens
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();

-- ============================
-- 17. Donations table (Donations and Sponsorships)
-- ============================
CREATE TABLE IF NOT EXISTS donations (
    id SERIAL PRIMARY KEY,
    type VARCHAR(50) NOT NULL,  -- 'donate' or 'sponsor'
    donor_type VARCHAR(50) NOT NULL,  -- 'user' or 'company'
    user_id INTEGER REFERENCES users(id) ON DELETE SET NULL,
    company_name VARCHAR(255),
    company_email VARCHAR(255),
    company_phone VARCHAR(50),
    amount DECIMAL(12, 2) NOT NULL CHECK (amount > 0),
    currency VARCHAR(10) NOT NULL DEFAULT 'USD',
    message TEXT,
    status VARCHAR(50) NOT NULL DEFAULT 'pending',
    transaction_id VARCHAR(255),
    payment_method VARCHAR(100),
    qr_expires_at TIMESTAMPTZ,  -- QR code expiration (3 minutes for donate type)
    event_id INTEGER REFERENCES events(id) ON DELETE SET NULL,
    created_at TIMESTAMPTZ DEFAULT now(),
    updated_at TIMESTAMPTZ DEFAULT now(),
    CHECK (type IN ('donate', 'sponsor')),
    CHECK (donor_type IN ('user', 'company', 'organization', 'church')),
    CHECK (status IN ('pending', 'completed', 'failed', 'refunded')),
    CHECK (
        (donor_type = 'user' AND user_id IS NOT NULL) OR
        (donor_type IN ('company', 'organization', 'church') AND company_name IS NOT NULL AND company_email IS NOT NULL)
    )
);

CREATE INDEX IF NOT EXISTS idx_donations_type ON donations(type);
CREATE INDEX IF NOT EXISTS idx_donations_donor_type ON donations(donor_type);
CREATE INDEX IF NOT EXISTS idx_donations_user_id ON donations(user_id);
CREATE INDEX IF NOT EXISTS idx_donations_event_id ON donations(event_id);
CREATE INDEX IF NOT EXISTS idx_donations_status ON donations(status);
CREATE INDEX IF NOT EXISTS idx_donations_created_at ON donations(created_at DESC);
CREATE INDEX IF NOT EXISTS idx_donations_transaction_id ON donations(transaction_id);
CREATE INDEX IF NOT EXISTS idx_donations_qr_expires_at ON donations(qr_expires_at);

-- donations table trigger
DROP TRIGGER IF EXISTS update_donations_updated_at ON donations;
CREATE TRIGGER update_donations_updated_at
BEFORE UPDATE ON donations
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();

-- ============================
-- 18. Supporters table (Companies and Organizations)
-- ============================
CREATE TABLE IF NOT EXISTS supporters (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,  -- Company/Organization/Church name
    email VARCHAR(255) NOT NULL UNIQUE,
    phone VARCHAR(50),
    type VARCHAR(50) NOT NULL DEFAULT 'company',  -- 'company', 'organization', or 'church'
    website VARCHAR(255),
    address TEXT,
    logo VARCHAR(255),
    description TEXT,
    user_id INTEGER REFERENCES users(id) ON DELETE SET NULL,  -- Optional: user who manages this supporter
    created_at TIMESTAMPTZ DEFAULT now(),
    updated_at TIMESTAMPTZ DEFAULT now(),
    CHECK (type IN ('company', 'organization', 'church'))
);

CREATE INDEX IF NOT EXISTS idx_supporters_email ON supporters(email);
CREATE INDEX IF NOT EXISTS idx_supporters_user_id ON supporters(user_id);
CREATE INDEX IF NOT EXISTS idx_supporters_type ON supporters(type);
CREATE INDEX IF NOT EXISTS idx_supporters_name ON supporters(name);

-- supporters table trigger
DROP TRIGGER IF EXISTS update_supporters_updated_at ON supporters;
CREATE TRIGGER update_supporters_updated_at
BEFORE UPDATE ON supporters
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();

-- Add supporter_id to donations table
ALTER TABLE donations ADD COLUMN IF NOT EXISTS supporter_id INTEGER REFERENCES supporters(id) ON DELETE SET NULL;
CREATE INDEX IF NOT EXISTS idx_donations_supporter_id ON donations(supporter_id);

-- Update the constraint to allow supporter reference
ALTER TABLE donations DROP CONSTRAINT IF EXISTS donations_check;
ALTER TABLE donations ADD CONSTRAINT donations_check CHECK (
    (donor_type = 'user' AND user_id IS NOT NULL) OR
    (donor_type IN ('company', 'organization', 'church') AND (
        (company_name IS NOT NULL AND company_email IS NOT NULL) OR 
        supporter_id IS NOT NULL
    ))
);

-- ============================
-- NOTES FOR DOWN MIGRATION
-- ============================
-- The down migration should:
-- 1. First drop the constraint: ALTER TABLE users DROP CONSTRAINT IF EXISTS fk_users_band_id;
-- 2. Then drop the column: ALTER TABLE users DROP COLUMN IF EXISTS band_id;
-- 3. Continue with normal table drops
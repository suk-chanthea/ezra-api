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
-- 2. Users table (with role_id FK)
-- ============================
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(100) NOT NULL,
    fullname VARCHAR(100) NOT NULL,
    profile VARCHAR(255) NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    role_id INTEGER NOT NULL REFERENCES roles(id) ON DELETE RESTRICT DEFAULT 3, -- default to 'user' role
    created_at TIMESTAMP WITH TIME ZONE DEFAULT now(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
CREATE INDEX IF NOT EXISTS idx_users_role_id ON users(role_id);

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
-- 5. Musics table (MUST BE BEFORE music_sheets)
-- ============================
CREATE TABLE IF NOT EXISTS musics (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    cover VARCHAR(255) NOT NULL,
    audio VARCHAR(255),
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ DEFAULT now(),
    updated_at TIMESTAMPTZ DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_musics_user_id ON musics(user_id);

-- ============================
-- 6. Music Sheets table (AFTER musics)
-- ============================
CREATE TABLE IF NOT EXISTS music_sheets (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    sheet VARCHAR(255) NOT NULL,
    music_id INTEGER NOT NULL REFERENCES musics(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ DEFAULT now(),
    updated_at TIMESTAMPTZ DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_music_sheets_title ON music_sheets(title);
CREATE INDEX IF NOT EXISTS idx_music_sheets_music_id ON music_sheets(music_id);

-- ============================
-- 7. Events table
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
-- 8. Event_Musics Junction Table (Many-to-Many)
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
-- 9. Bookings table (Event Registrations)
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
-- 10. Shared auto-update function
-- ============================
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = now();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- ============================
-- 11. Triggers for auto-updating
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
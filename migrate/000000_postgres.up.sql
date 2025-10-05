-- ============================
-- 1. Users table (example)
-- ============================
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(100) NOT NULL,
    fullname VARCHAR(100) NOT NULL,
    profile VARCHAR(255) NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    role VARCHAR(20) NOT NULL DEFAULT 'user',
    token VARCHAR(255),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT now(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);

-- ============================
-- Settings table
-- ============================
CREATE TABLE IF NOT EXISTS settings (
    id SERIAL NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT now(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT now()
);

-- ============================
-- 2. Music Sheets table
-- ============================
CREATE TABLE IF NOT EXISTS music_sheets (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    khmer VARCHAR(255) NOT NULL,
    english VARCHAR(255),
    korean VARCHAR(255),
    chinese VARCHAR(255),
    created_at TIMESTAMPTZ DEFAULT now(),
    updated_at TIMESTAMPTZ DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_music_sheets_title ON music_sheets(title);

-- ============================
-- 3. Musics table
-- ============================
CREATE TABLE IF NOT EXISTS musics (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    cover VARCHAR(255) NOT NULL,
    audio VARCHAR(255),
    sheet_id INTEGER NOT NULL REFERENCES music_sheets(id) ON DELETE CASCADE,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ DEFAULT now(),
    updated_at TIMESTAMPTZ DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_musics_user_id ON musics(user_id);
CREATE INDEX IF NOT EXISTS idx_musics_sheet_id ON musics(sheet_id);

-- ============================
-- 4. Events table
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
-- 5. Shared auto-update function
-- ============================
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = now();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- ============================
-- 6. Triggers for auto-updating
-- ============================

-- users table
DROP TRIGGER IF EXISTS update_users_updated_at ON users;
CREATE TRIGGER update_users_updated_at
BEFORE UPDATE ON users
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();

-- music_sheets table
DROP TRIGGER IF EXISTS update_music_sheets_updated_at ON music_sheets;
CREATE TRIGGER update_music_sheets_updated_at
BEFORE UPDATE ON music_sheets
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();

-- musics table
DROP TRIGGER IF EXISTS update_musics_updated_at ON musics;
CREATE TRIGGER update_musics_updated_at
BEFORE UPDATE ON musics
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();

-- events table
DROP TRIGGER IF EXISTS update_events_updated_at ON events;
CREATE TRIGGER update_events_updated_at
BEFORE UPDATE ON events
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();

DROP TRIGGER IF EXISTS update_settings_updated_at ON settings;
CREATE TRIGGER update_settings_updated_at
BEFORE UPDATE ON settings
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();

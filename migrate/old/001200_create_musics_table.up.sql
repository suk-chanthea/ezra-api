CREATE TABLE IF NOT EXISTS musics (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    cover VARCHAR(255) NOT NULL,
    audio VARCHAR(255),
    band VARCHAR(30),
    sheet_id INTEGER NOT NULL REFERENCES music_sheets(id) ON DELETE CASCADE,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ DEFAULT now(),
    updated_at TIMESTAMPTZ DEFAULT now()
);

--2. auto-updating
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
   NEW.updated_at = now();
   RETURN NEW;
END;
$$ LANGUAGE plpgsql;

--3. index
-- musics table indexes
CREATE INDEX IF NOT EXISTS idx_musics_user_id ON musics(user_id);
CREATE INDEX IF NOT EXISTS idx_musics_sheet_id ON musics(sheet_id);

--4. triggers
-- musics table trigger
DROP TRIGGER IF EXISTS update_musics_updated_at ON musics;
CREATE TRIGGER update_musics_updated_at
BEFORE UPDATE ON musics
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();

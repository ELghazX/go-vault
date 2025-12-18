CREATE TABLE IF NOT EXISTS users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    username TEXT UNIQUE NOT NULL,
    password_hash TEXT NOT NULL,
    created_at DATETIME NOT NULL
);

CREATE TABLE IF NOT EXISTS files (
    uuid TEXT PRIMARY KEY,
    filename TEXT NOT NULL,
    filepath TEXT NOT NULL,
    filesize INTEGER DEFAULT 0,
    content_type TEXT NOT NULL,
    owner_id INTEGER NOT NULL,
    is_onetime BOOLEAN DEFAULT FALSE,
    expires_at DATETIME NOT NULL,
    download_count INTEGER DEFAULT 0,
    created_at DATETIME NOT NULL,
    FOREIGN KEY (owner_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_files_owner_id ON files(owner_id);
CREATE INDEX IF NOT EXISTS idx_files_expires_at ON files(expires_at);

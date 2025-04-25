CREATE TABLE IF NOT EXISTS users (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  name TEXT NOT NULL,
  email TEXT UNIQUE NOT NULL,
  phone TEXT NOT NULL,
  password_hash BLOB NOT NULL,
  session_token_hash TEXT UNIQUE,
  session_expires_at DATETIME
);

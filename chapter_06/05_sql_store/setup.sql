CREATE TABLE IF NOT EXISTS posts (
  id        INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
  content   TEXT,
  author    VARCHAR(255)
);
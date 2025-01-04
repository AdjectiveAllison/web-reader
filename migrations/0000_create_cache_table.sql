-- Migration number: 0000
CREATE TABLE page_cache (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    url TEXT NOT NULL,
    title TEXT,
    markdown TEXT NOT NULL,
    fetched_at INTEGER NOT NULL,
    UNIQUE(url)
);
CREATE INDEX idx_page_cache_url ON page_cache(url);
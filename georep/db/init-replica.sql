CREATE TABLE IF NOT EXISTS articles (
    id SERIAL PRIMARY KEY,
    title   TEXT NOT NULL,
    content TEXT NOT NULL,
    author  TEXT NOT NULL,
    region  TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);



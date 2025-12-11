-- ✅ Yeni sürüm: tüm replikalar için güncel tablo yapısı
CREATE TABLE IF NOT EXISTS articles (
    id SERIAL PRIMARY KEY,
    title TEXT NOT NULL,
    summary TEXT NOT NULL,
    content_long TEXT NOT NULL,
    author TEXT NOT NULL,
    region TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

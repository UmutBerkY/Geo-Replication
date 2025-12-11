-- Initial schema for the master and replica databases.
-- We keep all three independent and let the Go service simulate replication.

CREATE TABLE IF NOT EXISTS articles (
    id SERIAL PRIMARY KEY,
    title   TEXT NOT NULL,
    content TEXT NOT NULL,
    author  TEXT NOT NULL,
    region  TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Başlangıçta 12 makale yüklü olsun.
INSERT INTO articles (title, content, author, region)
VALUES
  ('EU Master''ın Gücü', 'Tek yazma noktası, dağıtık okuma. Bu makale EU master üzerinden seed edildi.', 'Seed Bot', 'eu'),
  ('Replikalara Yakın Okuma', 'Okuyucular en yakın replikadan okuduğunda gecikme düşer.', 'Seed Bot', 'eu'),
  ('Eventual Consistency Nedir?', 'Yazı önce master''a düşer, replikalar birkaç saniye gecikmeyle güncellenir.', 'Seed Bot', 'eu'),
  ('US Bölgesi için Performans', 'US replikası, Amerika kıtasındaki okuyuculara hızlı yanıt verir.', 'Seed Bot', 'eu'),
  ('ASIA Replikası', 'Asya kıtasındaki okuyucular master yerine ASIA replikasından okur.', 'Seed Bot', 'eu'),
  ('TR Replikası', 'Türkiye yakınındaki kullanıcılar TR replikasından okur.', 'Seed Bot', 'eu'),
  ('SA Replikası', 'Güney Amerika replikası, oradaki okuyucular için düşük gecikme sağlar.', 'Seed Bot', 'eu'),
  ('AFRICA Replikası', 'Afrika replikası bölgedeki okuyuculara hizmet eder.', 'Seed Bot', 'eu'),
  ('Yazma Akışı', 'Yazar nerede olursa olsun, istek EU master üzerinden geçer.', 'Seed Bot', 'eu'),
  ('Okuma Akışı', 'Okuma isteği bölgeye göre en yakın replika veya master''a yönlenir.', 'Seed Bot', 'eu'),
  ('Gecikme Penceresi', 'Replikalar birkaç saniye geriden gelebilir; buna rağmen okuma hızlıdır.', 'Seed Bot', 'eu'),
  ('Test Senaryosu', 'Bu makale test için eklendi.', 'Seed Bot', 'eu')
ON CONFLICT DO NOTHING;



## Geo-Replicated Okuma/Yazma Servisi (Yakın okuma, Master yazma, Eventual Consistency)

Tek EU master’a yazma, bölgeye göre en yakın replikadan okuma, 5 replikaya gecikmeli kopyalama. Sahte login ile rol (okuyucu/yazar) ve bölge seçilip akış deneyimlenir.

### Yapı
- `backend/` Go (Gin) API
  - `/api/login` (sahte giriş, rol & bölge döner)
  - `/api/articles` GET (region param, en yakın node’dan okur)
  - `/api/articles` POST (yalnızca EU master’a yazar, replikalara gecikmeli kopyalar)
  - `/api/replication-status` (replikaların durumu)
- `frontend/` React (Vite) SPA
  - LoginPage → ReaderPage → WriterPage
- `db/init-master.sql` 12 hazır makale
- `db/init-replica.sql` şema
- `deploy/docker-compose.yml` master + 5 replika + backend + frontend

### Çalıştırma
```bash
cd deploy
docker compose up --build
```
- Frontend: http://localhost:5173
- API: http://localhost:8080/api

### Servisler
- postgres-master (EU)        : 5432
- postgres-replica1 (US)      : 5433
- postgres-replica2 (ASIA)    : 5434
- postgres-replica3 (TR)      : 5435
- postgres-replica4 (SA)      : 5436
- postgres-replica5 (AFRICA)  : 5437
- backend (Go)                : 8080
- frontend (React)            : 5173

### Ortam Değişkenleri (backend)
```
MASTER_DSN=postgres://admin:pass@postgres-master:5432/articles?sslmode=disable
REPLICA1_DSN=postgres://admin:pass@postgres-replica1:5432/articles?sslmode=disable
REPLICA2_DSN=postgres://admin:pass@postgres-replica2:5432/articles?sslmode=disable
REPLICA3_DSN=postgres://admin:pass@postgres-replica3:5432/articles?sslmode=disable
REPLICA4_DSN=postgres://admin:pass@postgres-replica4:5432/articles?sslmode=disable
REPLICA5_DSN=postgres://admin:pass@postgres-replica5:5432/articles?sslmode=disable
API_PORT=8080
```
Frontend: `VITE_API_BASE=http://localhost:8080/api`

### API Örnekleri
```bash
# Login (sahte)
curl -X POST http://localhost:8080/api/login \
  -H "Content-Type: application/json" \
  -d '{"username":"demo","role":"writer","region":"eu"}'

# Master’a yaz (EU)
curl -X POST http://localhost:8080/api/articles \
  -H "Content-Type: application/json" \
  -d '{"title":"Deneme","content":"İçerik","author":"demo"}'

# Bölgeye göre oku (ör. TR replikası)
curl "http://localhost:8080/api/articles?region=tr"

# Replikasyon durumu
curl http://localhost:8080/api/replication-status
```

### Eventual Consistency
- Yazı EU master’a düşer.
- 5 replikaya 2–3 sn gecikmeyle kopyalanır (kod içinde goroutine + gecikme).
- Okuyucu seçtiği bölgedeki replikadan okur; son yazıyı görmek için kısa süre bekleyebilir.

### Özet Akış (UI)
1. Giriş yap → Rol (okuyucu/yazar) + Bölge seç.
2. Okuyucu: Seçilen bölgeden `/articles?region=XX` ile okur, replikasyon durumunu görür.
3. Yazar: Her zaman EU master’a yazar; replikalar gecikmeyle güncellenir.

### Test İçin Farklı Bölgeleri Deneme

**Yöntem 1: Query Parameter (Hızlı Test)**
API çağrılarında `?region=` parametresi ile manuel bölge seçebilirsiniz:
```bash
# US bölgesini test et
curl "http://localhost:8080/api/articles?region=us"

# Asia bölgesini test et
curl "http://localhost:8080/api/articles?region=asia"

# SA bölgesini test et
curl "http://localhost:8080/api/articles?region=sa"

# Africa bölgesini test et
curl "http://localhost:8080/api/articles?region=africa"
```

**Yöntem 2: VPN ile Gerçek IP Testi**
Farklı bölgeleri test etmek için VPN kullanabilirsiniz:

**Ücretsiz VPN Seçenekleri:**
- **ProtonVPN** (ücretsiz plan mevcut)
- **Windscribe** (aylık 10GB ücretsiz)
- **TunnelBear** (aylık 500MB ücretsiz)
- Tarayıcı eklentileri (Location Guard, etc.)

**Test Senaryoları:**
- **US bölgesi**: ABD sunucusuna bağlanın
- **Asia bölgesi**: Japonya, Çin, Hindistan sunucularına bağlanın
- **SA bölgesi**: Brezilya, Arjantin sunucularına bağlanın
- **Africa bölgesi**: Güney Afrika sunucularına bağlanın
- **EU bölgesi**: Türkiye, Almanya, Fransa (varsayılan)

**Not:** 
- Türkiye artık EU bölgesi olarak tanımlanmaktadır.
- Docker içinden test ederken IP `172.19.0.1` gibi private IP'ler görünebilir. Bu durumda otomatik olarak EU bölgesi seçilir. Test için query parameter kullanabilirsiniz.




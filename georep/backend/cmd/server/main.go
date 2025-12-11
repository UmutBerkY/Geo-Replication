package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"geo-repl-demo/internal/article"
	"geo-repl-demo/internal/auth"
	"geo-repl-demo/internal/config"
	"geo-repl-demo/internal/db"
	"geo-repl-demo/internal/geoip"
	"geo-repl-demo/internal/middleware"
	"geo-repl-demo/internal/replication"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("❌ load config: %v", err)
	}

	// 🌍 GeoIP veritabanı yükle
	if err := geoip.Init("/app/GeoLite2-Country.mmdb"); err != nil {
		log.Printf("⚠️ GeoIP DB yüklenemedi: %v", err)
	} else {
		log.Println("🌍 GeoIP veritabanı başarıyla yüklendi")
	}

	// 🔗 Veritabanı bağlantıları
	masterDB := mustConnectMaster(cfg.MasterDSN)
	defer masterDB.Close()

	replicas := mustConnectReplicas(cfg.ReplicaDSNs())
	defer replicas.Close()

	// 🧱 Şema ve varsayılan veriler
	mustEnsureSchema(masterDB)
	addDefaultArticles(masterDB)

	repo := article.NewRepository(masterDB, replicas)
	replicator := replication.NewReplicator(masterDB, replicas)
	svc := article.NewService(repo, replicator)

	log.Println("🔁 İlk replikasyon başlatılıyor...")
	replicator.FullSync()

	// ⏰ Periyodik tam senkronizasyon
	go func() {
		for range time.Tick(10 * time.Second) {
			replicator.FullSync()
		}
	}()

	// 🌐 HTTP Sunucu
	r := gin.Default()
	r.SetTrustedProxies(nil)
	r.ForwardedByClientIP = true
	r.Use(cors.Default())
	r.Use(middleware.RegionMiddleware())

	authHandler := auth.NewHandler()
	articleHandler := article.NewHandler(svc)
	auth.RegisterRoutes(r, authHandler)
	article.RegisterRoutes(r, articleHandler)

	// 🌍 IP tabanlı bölge tespiti
	r.GET("/api/region", func(c *gin.Context) {
		clientIP := c.ClientIP()
		region := "unknown"
		if rVal, ok := c.Get("region"); ok {
			if val, ok := rVal.(string); ok {
				region = val
			}
		}
		c.JSON(200, gin.H{
			"region": region,
			"ip":     clientIP,
		})
	})

	// ⚡ Gecikme ölçüm endpoint’i
	r.GET("/api/latency", func(c *gin.Context) {
		region := c.Query("region")
		if region == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "region parametresi gerekli"})
			return
		}

		result, err := svc.MeasureLatency(c, region)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// ❗ Burada sadece string döndürüyoruz, frontend bu string’i olduğu gibi gösteriyor
		c.JSON(http.StatusOK, gin.H{
			"region":   region,
			"latency":  result,
			"measured": time.Now().Format(time.RFC3339),
		})
	})

	addr := ":" + cfg.APIPort
	log.Printf("✅ API listening on %s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("server error: %v", err)
	}
}

// ==========================================================
// 🔧 Yardımcı Fonksiyonlar
// ==========================================================

func mustConnectMaster(dsn string) *db.Master {
	for i := 0; i < 10; i++ {
		m, err := db.NewMaster(dsn)
		if err == nil {
			return m
		}
		log.Printf("master db bağlantısı başarısız (deneme %d/10): %v", i+1, err)
		time.Sleep(2 * time.Second)
	}
	log.Fatalf("connect master db başarısız")
	return nil
}

func mustConnectReplicas(dsns []string) *db.ReplicaSet {
	for i := 0; i < 10; i++ {
		r, err := db.NewReplicas(dsns)
		if err == nil {
			return r
		}
		log.Printf("replica db bağlantısı başarısız (deneme %d/10): %v", i+1, err)
		time.Sleep(2 * time.Second)
	}
	log.Fatalf("connect replicas başarısız")
	return nil
}

func mustEnsureSchema(master *db.Master) {
	for i := 0; i < 10; i++ {
		if err := db.EnsureSchema(master); err != nil {
			log.Printf("schema oluşturma başarısız (deneme %d/10): %v", i+1, err)
			time.Sleep(2 * time.Second)
			continue
		}
		return
	}
	log.Fatalf("schema oluşturma başarısız")
}

func addDefaultArticles(master *db.Master) {
	log.Println("🪶 Varsayılan makaleler ekleniyor (sadece master)")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := master.Pool.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS articles (
			id SERIAL PRIMARY KEY,
			title TEXT NOT NULL,
			content TEXT NOT NULL,
			author TEXT NOT NULL,
			region TEXT NOT NULL,
			created_at TIMESTAMP DEFAULT NOW()
		);

		INSERT INTO articles (title, content, author, region, created_at)
		SELECT 'Yeni Nesil Replikasyon Sistemi', 
		       'Veriler artık dünya genelinde anlık olarak kopyalanıyor.', 
		       'Admin', 'eu', NOW()
		WHERE NOT EXISTS (SELECT 1 FROM articles WHERE title='Yeni Nesil Replikasyon Sistemi');

		INSERT INTO articles (title, content, author, region, created_at)
		SELECT 'Küresel Veri Senkronizasyonu', 
		       'Asia ve US bölgelerindeki veriler saniyeler içinde güncelleniyor.', 
		       'Sistem', 'eu', NOW()
		WHERE NOT EXISTS (SELECT 1 FROM articles WHERE title='Küresel Veri Senkronizasyonu');

		INSERT INTO articles (title, content, author, region, created_at)
		SELECT 'Go ve PostgreSQL ile Dağıtık Sistemler', 
		       'Bu demo, Go dilinde PostgreSQL replikasyonunu gösteriyor.', 
		       'DemoUser', 'eu', NOW()
		WHERE NOT EXISTS (SELECT 1 FROM articles WHERE title='Go ve PostgreSQL ile Dağıtık Sistemler');

		INSERT INTO articles (title, content, author, region, created_at)
		SELECT 'Veri Tutarlılığı Testi', 
		       'Eventual consistency mekanizması başarıyla çalışıyor.', 
		       'Tester', 'eu', NOW()
		WHERE NOT EXISTS (SELECT 1 FROM articles WHERE title='Veri Tutarlılığı Testi');

		INSERT INTO articles (title, content, author, region, created_at)
		SELECT 'TR Bölgesi için İlk Veri', 
		       'Türkiye replikası 2 saniye gecikmeyle senkronize olur.', 
		       'Umut', 'eu', NOW()
		WHERE NOT EXISTS (SELECT 1 FROM articles WHERE title='TR Bölgesi için İlk Veri');
	`)
	if err != nil {
		log.Printf("⚠️ Varsayılan makale ekleme hatası: %v", err)
	} else {
		log.Println("✅ Varsayılan makaleler eklendi")
	}
}

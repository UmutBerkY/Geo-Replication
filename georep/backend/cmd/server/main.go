package main

import (
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

	// 🧱 Şema kontrolü (addDefaultArticles artık yok)
	mustEnsureSchema(masterDB)

	// 🧩 Repository & Servis
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

		reg, result := svc.MeasureLatency(region)

		c.JSON(http.StatusOK, gin.H{
			"region":   reg,
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

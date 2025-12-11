package geoip

import (
	"log"
	"net"
	"strings"

	"github.com/oschwald/geoip2-golang"
)

var db *geoip2.Reader

// 🌍 GeoLite2 veritabanını başlat
func Init(path string) error {
	var err error
	db, err = geoip2.Open(path)
	if err != nil {
		return err
	}
	return nil
}

// 🌐 IP adresinden otomatik bölge belirle
func RegionFromIP(ip string) string {
	if db == nil {
		log.Printf("⚠️ GeoIP veritabanı yüklenmemiş, varsayılan bölge 'eu' kullanılıyor")
		return "eu" // varsayılan
	}

	// IP adresini temizle (port varsa kaldır)
	if idx := strings.Index(ip, ":"); idx != -1 {
		ip = ip[:idx]
	}

	parsed := net.ParseIP(ip)
	if parsed == nil {
		log.Printf("⚠️ Geçersiz IP adresi: %s, varsayılan bölge 'eu' kullanılıyor", ip)
		return "eu"
	}

	// Private IP'ler için GeoIP lookup yapma (spam'i önle)
	if parsed.IsPrivate() || parsed.IsLoopback() || parsed.IsLinkLocalUnicast() {
		return "eu" // Sessizce varsayılan dön
	}

	record, err := db.Country(parsed)
	if err != nil {
		// Sadece gerçek hatalarda log (private IP'ler için değil)
		log.Printf("⚠️ GeoIP lookup hatası (IP: %s): %v", ip, err)
		return "eu"
	}

	country := strings.ToUpper(record.Country.IsoCode)
	if country == "" {
		log.Printf("⚠️ Ülke kodu bulunamadı (IP: %s)", ip)
		return "eu"
	}

	log.Printf("🌍 IP: %s → Ülke: %s", ip, country)

	// 🌎 Ülke koduna göre replikasyon bölgesi
	switch country {
	case "US", "CA", "MX":
		return "us"

	case "CN", "JP", "KR", "IN", "ID", "SG", "PH", "TH", "VN", "MY", "TW", "HK":
		return "asia"

	case "TR", "DE", "FR", "IT", "ES", "GB", "NL", "PL", "SE", "NO", "FI", "DK", 
		 "BE", "AT", "CH", "PT", "GR", "CZ", "HU", "RO", "BG", "HR", "SK", "SI",
		 "IE", "IS", "LU", "EE", "LV", "LT", "MT", "CY", "RS", "BA", "MK", "AL",
		 "ME", "XK", "MD", "UA", "BY", "RU", "GE", "AM", "AZ", "KZ", "UZ", "KG",
		 "TJ", "TM", "AE", "SA", "IL", "QA", "KW", "BH", "OM", "YE", "JO", "LB",
		 "IQ", "IR", "PS", "SY":
		return "eu"

	case "BR", "AR", "CL", "CO", "PE", "VE", "EC", "BO", "PY", "UY", "GY", "SR":
		return "sa"

	case "ZA", "NG", "EG", "KE", "ET", "GH", "TZ", "UG", "DZ", "MA", "TN", "LY",
		 "SD", "SS", "CM", "CI", "SN", "BF", "ML", "NE", "TD", "MR", "DJ", "SO",
		 "ER", "RW", "BI", "MW", "ZM", "ZW", "BW", "NA", "LS", "SZ", "MG", "MU",
		 "SC", "KM", "AO", "MZ", "CD", "CF", "CG", "GA", "GQ", "ST", "CV", "GW",
		 "GN", "SL", "LR", "TG", "BJ":
		return "africa"

	default:
		return "eu"
	}
}

// 🌐 Bölgeye göre replika adresi
func ReplicaForRegion(region string) string {
	switch region {
	case "us":
		return "postgres-replica1:5432"
	case "eu":
		return "postgres-replica2:5432"
	case "asia":
		return "postgres-replica3:5432"
	case "sa":
		return "postgres-replica4:5432"
	case "africa":
		return "postgres-replica5:5432"
	default:
		return "postgres-master:5432"
	}
}

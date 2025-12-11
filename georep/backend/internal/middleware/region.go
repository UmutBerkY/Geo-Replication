package middleware

import (
	"geo-repl-demo/internal/geoip"
	"log"
	"net"
	"strings"

	"github.com/gin-gonic/gin"
)

func RegionMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Query parameter ile manuel bölge override (test için)
		if regionParam := c.Query("region"); regionParam != "" {
			validRegions := map[string]bool{"us": true, "eu": true, "asia": true, "sa": true, "africa": true}
			if validRegions[strings.ToLower(regionParam)] {
				c.Set("region", strings.ToLower(regionParam))
				c.Set("client_ip", "test-override")
				log.Printf("🌍 Test modu: Manuel bölge seçildi → %s", strings.ToLower(regionParam))
				c.Next()
				return
			}
		}

		// IP adresini al - önce X-Forwarded-For, sonra X-Real-IP, son olarak ClientIP
		clientIP := c.GetHeader("X-Forwarded-For")
		if clientIP == "" {
			clientIP = c.GetHeader("X-Real-IP")
		}
		if clientIP == "" {
			clientIP = c.ClientIP()
		}
		
		// X-Forwarded-For birden fazla IP içerebilir (proxy chain), ilkini al
		if idx := strings.Index(clientIP, ","); idx != -1 {
			clientIP = strings.TrimSpace(clientIP[:idx])
		}
		
		// Private IP kontrolü (Docker network, localhost, vb.)
		parsedIP := net.ParseIP(clientIP)
		isPrivate := false
		if parsedIP != nil {
			isPrivate = parsedIP.IsLoopback() || parsedIP.IsPrivate() || parsedIP.IsLinkLocalUnicast()
		} else {
			isPrivate = clientIP == "::1" || clientIP == "127.0.0.1" || clientIP == "localhost"
		}

		// Private IP ise sessizce varsayılan bölge kullan (log spam'ini önle)
		if isPrivate {
			c.Set("region", "eu")
			c.Set("client_ip", clientIP)
			c.Next()
			return
		}

		// Public IP için GeoIP lookup yap
		region := geoip.RegionFromIP(clientIP)
		
		// Sadece başarılı lookup'larda log (spam'i azalt)
		if region != "eu" || !isPrivate {
			log.Printf("🌍 Client IP: %s → Bölge: %s", clientIP, region)
		}

		c.Set("region", region)
		c.Set("client_ip", clientIP)
		c.Next()
	}
}

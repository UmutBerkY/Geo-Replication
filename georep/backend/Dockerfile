# 1️⃣ Build aşaması
FROM golang:1.23-alpine AS build

WORKDIR /app

# go.mod ve go.sum dosyalarını kopyala
COPY go.mod go.sum ./
RUN go mod download

# Kaynak kodları kopyala
COPY . .

# Derleme
RUN go build -o /app/server ./cmd/server

# 2️⃣ Final aşaması
FROM alpine:3.20

WORKDIR /app

# Derlenmiş binary'yi kopyala
COPY --from=build /app/server /app/server

# GeoIP veritabanını kopyala
COPY GeoLite2-Country.mmdb /app/GeoLite2-Country.mmdb

# Port ve çalıştırma komutu
ENV API_PORT=8080
EXPOSE 8080

CMD ["/app/server"]

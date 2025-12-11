FROM golang:1.23-alpine AS build

WORKDIR /app/backend

# Gerekli araçları kur ve modülleri indir
RUN apk add --no-cache git
COPY backend/go.mod backend/go.sum ./
RUN go env -w GOPROXY=direct && go mod download

# Uygulama kodunu kopyala
COPY backend/ ./
RUN go build -o /app/server ./cmd/server

# ---- Runtime Stage ----
FROM alpine:3.20
WORKDIR /app

# Binary ve GeoIP veritabanını ekle
COPY --from=build /app/server /app/server
COPY backend/GeoLite2-Country.mmdb /app/GeoLite2-Country.mmdb

ENV API_PORT=8080
EXPOSE 8080

CMD ["/app/server"]

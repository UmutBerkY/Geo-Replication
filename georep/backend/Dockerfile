# ---- Build Stage ----
FROM golang:1.23-alpine AS build
WORKDIR /app

# Bağımlılıklar
COPY georep/backend/go.mod georep/backend/go.sum ./
RUN go env -w GOPROXY=direct && go mod download

# Kaynak kodları kopyala
COPY georep/backend/ ./backend/

WORKDIR /app/backend
RUN go build -o /app/server ./cmd/server

# ---- Runtime Stage ----
FROM alpine:3.20
WORKDIR /app

# Binary ve GeoIP veritabanını ekle
COPY --from=build /app/server /app/server
COPY georep/backend/GeoLite2-Country.mmdb /app/GeoLite2-Country.mmdb

ENV API_PORT=8080
EXPOSE 8080

CMD ["/app/server"]

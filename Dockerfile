# ---- Build Stage ----
FROM golang:1.23-alpine AS build
WORKDIR /app

# 🧩 Git eklendi
RUN apk add --no-cache git

COPY georep/backend/go.mod georep/backend/go.sum ./
RUN go env -w GOPROXY=direct && go mod download

COPY georep/backend ./backend
WORKDIR /app/backend
RUN go build -o /app/server ./cmd/server

# ---- Runtime Stage ----
FROM alpine:3.20
WORKDIR /app

COPY --from=build /app/server /app/server
COPY georep/backend/GeoLite2-Country.mmdb /app/GeoLite2-Country.mmdb

ENV API_PORT=8080
EXPOSE 8080

CMD ["/app/server"]

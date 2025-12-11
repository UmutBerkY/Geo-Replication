package replication

import (
    "context"
    "log"
    "time"

    "geo-repl-demo/internal/db"
    "geo-repl-demo/internal/model"

    "github.com/jackc/pgx/v5/pgxpool"
)

type Replicator struct {
    master   *db.Master
    replicas *db.ReplicaSet
}

// Constructor
func NewReplicator(master *db.Master, replicas *db.ReplicaSet) *Replicator {
    return &Replicator{master: master, replicas: replicas}
}

// Yeni makale eklendiğinde çağrılır
func (r *Replicator) Schedule(a model.Article) {
    if r.replicas == nil {
        return
    }

    for i, pool := range r.replicas.Pools {
        go func(idx int, pool *pgxpool.Pool) {
            time.Sleep(2 * time.Second) // eventual consistency gecikmesi
            ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
            defer cancel()

            // Her replikada tabloyu garanti altına al
            _, _ = pool.Exec(ctx, `
                CREATE TABLE IF NOT EXISTS articles (
                    id SERIAL PRIMARY KEY,
                    title TEXT NOT NULL,
                    content TEXT NOT NULL,
                    author TEXT NOT NULL,
                    region TEXT NOT NULL,
                    created_at TIMESTAMP DEFAULT NOW()
                );
            `)

            _, err := pool.Exec(ctx, `
                INSERT INTO articles (id, title, content, author, region, created_at)
                VALUES ($1,$2,$3,$4,$5,$6)
                ON CONFLICT (id) DO UPDATE
                SET title=EXCLUDED.title,
                    content=EXCLUDED.content,
                    author=EXCLUDED.author,
                    region=EXCLUDED.region,
                    created_at=EXCLUDED.created_at;
            `, a.ID, a.Title, a.Content, a.Author, a.Region, a.CreatedAt)

            if err != nil {
                log.Printf("❌ Replikasyon hatası (replica %d): %v", idx+1, err)
            } else {
                log.Printf("✅ Article %d kopyalandı → replica %d", a.ID, idx+1)
            }
        }(i, pool)
    }
}

// Periyodik tam senkronizasyon (Master → tüm replikalar)
func (r *Replicator) FullSync() {
    if r.replicas == nil {
        return
    }

    ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
    defer cancel()

    rows, err := r.master.Pool.Query(ctx, `
        SELECT id, title, content, author, region, created_at
        FROM articles
    `)
    if err != nil {
        log.Printf("⚠️ Master verilerini okuma hatası: %v", err)
        return
    }
    defer rows.Close()

    var articles []model.Article
    for rows.Next() {
        var a model.Article
        if err := rows.Scan(&a.ID, &a.Title, &a.Content, &a.Author, &a.Region, &a.CreatedAt); err == nil {
            articles = append(articles, a)
        }
    }

    for i, pool := range r.replicas.Pools {
        // Her replikada tabloyu garanti et
        _, _ = pool.Exec(ctx, `
            CREATE TABLE IF NOT EXISTS articles (
                id SERIAL PRIMARY KEY,
                title TEXT NOT NULL,
                content TEXT NOT NULL,
                author TEXT NOT NULL,
                region TEXT NOT NULL,
                created_at TIMESTAMP DEFAULT NOW()
            );
        `)

        for _, a := range articles {
            _, err := pool.Exec(ctx, `
                INSERT INTO articles (id, title, content, author, region, created_at)
                VALUES ($1,$2,$3,$4,$5,$6)
                ON CONFLICT (id) DO UPDATE SET
                    title=EXCLUDED.title,
                    content=EXCLUDED.content,
                    author=EXCLUDED.author,
                    region=EXCLUDED.region,
                    created_at=EXCLUDED.created_at
            `, a.ID, a.Title, a.Content, a.Author, a.Region, a.CreatedAt)
            if err != nil {
                log.Printf("⚠️ FullSync hata (replica %d): %v", i+1, err)
            }
        }
        log.Printf("✅ FullSync: replica %d güncellendi (%d makale)", i+1, len(articles))
    }
}

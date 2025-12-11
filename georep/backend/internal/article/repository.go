package article

import (
    "context"
    "fmt"

    "github.com/jackc/pgx/v5/pgxpool"
    "geo-repl-demo/internal/db"
    "geo-repl-demo/internal/model"
)

type Repository struct {
    master   *db.Master
    replicas *db.ReplicaSet
}

func NewRepository(master *db.Master, replicas *db.ReplicaSet) *Repository {
    return &Repository{master: master, replicas: replicas}
}

// Master’a yazar
func (r *Repository) InsertMaster(ctx context.Context, in model.CreateArticleInput, region string) (model.Article, error) {
    var a model.Article
    err := r.master.Pool.QueryRow(ctx, `
        INSERT INTO articles (title, content, author, region)
        VALUES ($1, $2, $3, $4)
        RETURNING id, title, content, author, region, created_at
    `, in.Title, in.Content, in.Author, region).
        Scan(&a.ID, &a.Title, &a.Content, &a.Author, &a.Region, &a.CreatedAt)
    if err != nil {
        return model.Article{}, fmt.Errorf("insert master: %w", err)
    }
    return a, nil
}

// Replikaya kopyalar
func (r *Repository) CopyToReplica(ctx context.Context, idx int, a model.Article) error {
    if idx < 0 || r.replicas == nil || idx >= len(r.replicas.Pools) {
        return fmt.Errorf("invalid replica index %d", idx)
    }
    pool := r.replicas.Pools[idx]
    _, err := pool.Exec(ctx, `
        INSERT INTO articles (id, title, content, author, region, created_at)
        VALUES ($1,$2,$3,$4,$5,$6)
        ON CONFLICT (id) DO UPDATE
        SET title=EXCLUDED.title,
            content=EXCLUDED.content,
            author=EXCLUDED.author,
            region=EXCLUDED.region,
            created_at=EXCLUDED.created_at
    `, a.ID, a.Title, a.Content, a.Author, a.Region, a.CreatedAt)
    return err
}

// Silme işlemleri
func (r *Repository) DeleteFromMaster(ctx context.Context, id int64) error {
    _, err := r.master.Pool.Exec(ctx, `DELETE FROM articles WHERE id=$1`, id)
    return err
}

func (r *Repository) DeleteFromReplica(ctx context.Context, idx int, id int64) error {
    if idx < 0 || r.replicas == nil || idx >= len(r.replicas.Pools) {
        return fmt.Errorf("invalid replica index %d", idx)
    }
    pool := r.replicas.Pools[idx]
    _, err := pool.Exec(ctx, `DELETE FROM articles WHERE id=$1`, id)
    return err
}

// 🔍 Bölgeye göre okuma
func (r *Repository) ListByRegion(ctx context.Context, region string) ([]model.Article, error) {
    pool := r.poolForRegion(region)

    rows, err := pool.Query(ctx, `
        SELECT id, title, content, author, region, created_at
        FROM articles
        ORDER BY created_at DESC
    `)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var res []model.Article
    for rows.Next() {
        var a model.Article
        if err := rows.Scan(&a.ID, &a.Title, &a.Content, &a.Author, &a.Region, &a.CreatedAt); err != nil {
            return nil, err
        }
        res = append(res, a)
    }
    return res, nil
}

// 🧭 Replika seçimi
func (r *Repository) poolForRegion(region string) *pgxpool.Pool {
    if r.replicas == nil || len(r.replicas.Pools) == 0 {
        return r.master.Pool
    }

    switch region {
    case "us":
        if len(r.replicas.Pools) > 0 {
            return r.replicas.Pools[0]
        }
    case "eu":
        return r.master.Pool // ✅ EU her zaman master’dan gelir
    case "asia":
        if len(r.replicas.Pools) > 1 {
            return r.replicas.Pools[1]
        }
    case "sa":
        if len(r.replicas.Pools) > 2 {
            return r.replicas.Pools[2]
        }
    case "africa":
        if len(r.replicas.Pools) > 3 {
            return r.replicas.Pools[3]
        }
    }

    return r.master.Pool
}

func (r *Repository) NumReplicas() int {
    if r.replicas == nil {
        return 0
    }
    return len(r.replicas.Pools)
}

// ✅ RegionToIndex: MeasureLatency fonksiyonu için gerekli
func (r *Repository) RegionToIndex(region string) int {
    switch region {
    case "us":
        return 0
    case "asia":
        return 1
    case "sa":
        return 2
    case "africa":
        return 3
    default:
        return -1
    }
}

package article

import (
    "context"
    "fmt"
    "math/rand"
    "strings"
    "time"

    "geo-repl-demo/internal/model"
)

// Service katmanı
type Service struct {
    repo       *Repository
    replicator interface {
        Schedule(a model.Article)
        FullSync()
    }
}

func NewService(repo *Repository, replicator interface {
    Schedule(a model.Article)
    FullSync()
}) *Service {
    return &Service{repo: repo, replicator: replicator}
}

// ✅ Master’a yaz, ardından replikalara asenkron gönder
func (s *Service) Create(ctx context.Context, in model.CreateArticleInput) (model.Article, error) {
    a, err := s.repo.InsertMaster(ctx, in, "eu")
    if err != nil {
        return model.Article{}, err
    }
    s.replicator.Schedule(a)
    return a, nil
}

// ✅ Bölgeye göre oku
func (s *Service) ListByRegion(ctx context.Context, region string) ([]model.Article, error) {
    return s.repo.ListByRegion(ctx, region)
}

// ✅ Silme (master + tüm replikalar)
func (s *Service) Delete(ctx context.Context, id int64) error {
    if err := s.repo.DeleteFromMaster(ctx, id); err != nil {
        return err
    }
    for i := range s.repo.replicas.Pools {
        _ = s.repo.DeleteFromReplica(ctx, i, id)
    }
    return nil
}

// ✅ Replikasyon durumlarını oku
func (s *Service) ReplicationStatus(ctx context.Context) ([]model.ReplicationStatus, error) {
    var masterCount int
    err := s.repo.master.Pool.QueryRow(ctx, `SELECT COUNT(*) FROM articles`).Scan(&masterCount)
    if err != nil {
        return nil, err
    }

    replicaNames := []string{"US", "Asia", "SA", "Africa"}
    statuses := make([]model.ReplicationStatus, 0, len(replicaNames))

    for i, pool := range s.repo.replicas.Pools {
        if i >= len(replicaNames) {
            break
        }
        var replicaCount int
        var lastUpdate *time.Time
        err := pool.QueryRow(ctx, `SELECT COUNT(*), MAX(created_at) FROM articles`).Scan(&replicaCount, &lastUpdate)

        state := "syncing"
        if err != nil {
            state = "error"
        } else if replicaCount == masterCount && masterCount > 0 {
            state = "ok"
        }

        lastAt := time.Time{}
        if lastUpdate != nil {
            lastAt = *lastUpdate
        }

        statuses = append(statuses, model.ReplicationStatus{
            Replica: replicaNames[i],
            Status:  state,
            LastAt:  lastAt,
        })
    }
    return statuses, nil
}

func (s *Service) MeasureLatency(ctx context.Context, region string) (string, error) {
    rand.Seed(time.Now().UnixNano())

    // 🌍 Ortalama gecikme değerleri (ms)
    simulatedReplicaDelay := map[string]int64{
        "eu":     0,   // Master kendi
        "us":     80,
        "asia":   120,
        "sa":     100,
        "africa": 130,
    }

    simulatedMasterDelay := map[string]int64{
        "eu":     5,   // Master'a doğrudan erişim
        "us":     220,
        "asia":   270,
        "sa":     250,
        "africa": 260,
    }

    replicaDelay, ok1 := simulatedReplicaDelay[region]
    masterDelay, ok2 := simulatedMasterDelay[region]
    if !ok1 || !ok2 {
        replicaDelay = 100
        masterDelay = 250
    }

    // 🎲 ±10 ms jitter
    jitter := func(base int64) int64 {
        j := base + int64(rand.Intn(21)-10)
        if j < 0 {
            j = 0
        }
        return j
    }

    masterDelay = jitter(masterDelay)
    replicaDelay = jitter(replicaDelay)

    // 🧭 Eğer EU bölgesindeyse → doğrudan master’dan oku
    if strings.ToLower(region) == "eu" {
        replicaDelay = masterDelay
    }

    // 🕒 Simülasyon
    time.Sleep(time.Duration(masterDelay) * time.Millisecond)
    masterLatency := masterDelay
    time.Sleep(time.Duration(replicaDelay) * time.Millisecond)
    replicaLatency := replicaDelay

    // 🔢 Kazanç (Master - Replica)
    gain := masterLatency - replicaLatency
    if gain < 0 {
        gain = 0
    }

    // ✅ Formatlı çıktı
    result := fmt.Sprintf(
        "⏱ Master’a göre gecikme kazancı: %d ms (%s=%d ms, Master=%d ms)",
        gain,
        strings.ToUpper(region),
        replicaLatency,
        masterLatency,
    )

    return result, nil
}


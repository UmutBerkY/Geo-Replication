package article

import (
	"context"
	"fmt"
	"math/rand"
	"strings"
	"sync"
	"time"

	"geo-repl-demo/internal/model"
	"geo-repl-demo/internal/replication"
)

// Service iş katmanı (Repository + Replicator’ı birleştiriyor)
type Service struct {
	repo       *Repository
	replicator *replication.Replicator

	mu               sync.Mutex
	lastReplicaWrite []time.Time // her replikaya son yazma/silme zamanı (syncing göstermek için)
}

// Yeni servis oluşturur
func NewService(repo *Repository, replicator *replication.Replicator) *Service {
	return &Service{
		repo:       repo,
		replicator: replicator,
	}
}

// 🔹 Makaleleri bölgeye göre getir
func (s *Service) ListByRegion(ctx context.Context, region string) ([]model.Article, error) {
	return s.repo.ListByRegion(ctx, region)
}

// 🔹 Yeni makale ekle (master’a)
func (s *Service) Create(ctx context.Context, in model.CreateArticleInput) (*model.Article, error) {
	// Her zaman EU master’a yazıyoruz
	a, err := s.repo.InsertMaster(ctx, in, "eu")
	if err != nil {
		return nil, err
	}

	// Replikasyon başlat (eventual consistency)
	if s.replicator != nil {
		go s.replicator.Schedule(a)
	}

	// Replikasyon durumu için kısa süre "syncing" göster
	s.markReplicasSyncing()

	return &a, nil
}

// 🔹 Makale sil – master + tüm replikalardan
func (s *Service) Delete(ctx context.Context, id int64) error {
	// Önce master’dan sil
	if err := s.repo.DeleteFromMaster(ctx, id); err != nil {
		return err
	}

	// Tüm replikalardan da sil
	n := s.repo.NumReplicas()
	for i := 0; i < n; i++ {
		_ = s.repo.DeleteFromReplica(ctx, i, id)
	}

	// Silme işlemi de bir “replikasyon olayı” – kısa süre syncing gösterelim
	s.markReplicasSyncing()

	return nil
}

// 🔹 Replikasyon durumu (US/ASIA/SA/TR/AFRICA + syncing/ok)
func (s *Service) ReplicationStatus(ctx context.Context) ([]model.ReplicationStatus, error) {
	count := s.repo.NumReplicas()
	regionLabels := []string{"US", "ASIA", "SA", "TR", "AFRICA"} // EU yok, çünkü master

	statuses := make([]model.ReplicationStatus, 0, count)

	// lastReplicaWrite için thread-safe snapshot al
	s.mu.Lock()
	snapshot := make([]time.Time, len(s.lastReplicaWrite))
	copy(snapshot, s.lastReplicaWrite)
	s.mu.Unlock()

	now := time.Now()

	for i := 0; i < count; i++ {
		label := fmt.Sprintf("Replica %d", i+1)
		if i < len(regionLabels) {
			label = regionLabels[i]
		}

		status := "ok"
		if i < len(snapshot) && !snapshot[i].IsZero() {
			// Son yazma/silme olayı üzerinden 3 sn’den az geçmişse "syncing"
			if now.Sub(snapshot[i]) < 3*time.Second {
				status = "syncing"
			}
		}

		statuses = append(statuses, model.ReplicationStatus{
			Replica: label,
			Status:  status,
			LastAt:  now,
		})
	}

	return statuses, nil
}

// ⏱ Master’a göre gecikme kazancı ölçümü (frontend için)
// NOT: Handler sadece dönen stringi gösteriyor.
func (s *Service) MeasureLatency(region string) (string, string) {
	rand.Seed(time.Now().UnixNano())
	r := strings.ToLower(region)

	var latencyRegion, latencyMaster int

	switch r {
	case "eu":
		// EU: master ile aynı – kazanç yok, ikisi de aynı değer
		base := rand.Intn(20) + 25 // 25–44 ms
		latencyRegion = base
		latencyMaster = base
	case "us":
		// US: master (EU) uzak, replikaya yakın
		latencyRegion = rand.Intn(40) + 40  // 40–79 ms (US replikası)
		latencyMaster = rand.Intn(80) + 160 // 160–239 ms (EU master)
	case "asia":
		// ASIA: master daha uzak
		latencyRegion = rand.Intn(50) + 50  // 50–99 ms
		latencyMaster = rand.Intn(90) + 180 // 180–269 ms
	case "sa":
		// South America
		latencyRegion = rand.Intn(45) + 45  // 45–89 ms
		latencyMaster = rand.Intn(90) + 170 // 170–259 ms
	case "africa", "tr":
		// TR ve AFRICA – master’a görece uzak, kendi replikası daha yakın
		latencyRegion = rand.Intn(35) + 35  // 35–69 ms
		latencyMaster = rand.Intn(80) + 150 // 150–229 ms
	default:
		// Bilinmeyen bölge -> hafif fark
		latencyRegion = rand.Intn(60) + 60  // 60–119 ms
		latencyMaster = rand.Intn(40) + 120 // 120–159 ms
	}

	diff := latencyMaster - latencyRegion
	if diff < 0 {
		diff = 0
	}

	result := fmt.Sprintf(
		"⏱ Master’a göre gecikme kazancı: %d ms (%s=%d ms, Master=%d ms)",
		diff, strings.ToUpper(r), latencyRegion, latencyMaster,
	)

	return r, result
}

// ------------------------------------------------------
//  Yardımcı: Replikaları kısa süre “syncing” durumuna al
// ------------------------------------------------------
func (s *Service) markReplicasSyncing() {
	s.mu.Lock()
	defer s.mu.Unlock()

	count := s.repo.NumReplicas()
	if len(s.lastReplicaWrite) < count {
		s.lastReplicaWrite = make([]time.Time, count)
	}

	now := time.Now()
	for i := 0; i < count; i++ {
		s.lastReplicaWrite[i] = now
	}
}

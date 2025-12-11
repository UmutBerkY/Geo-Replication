package location

import (
	"context"
	"math/rand"
	"time"
)

// Service implements business logic around locations and replication.
type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

// CreateLocation writes to master and asynchronously propagates to replicas.
// Replication is intentionally delayed and jittered to show eventual consistency.
func (s *Service) CreateLocation(ctx context.Context, in CreateLocationInput) (Location, error) {
	loc, err := s.repo.InsertToMaster(ctx, in)
	if err != nil {
		return Location{}, err
	}

	// Fire-and-forget goroutines simulating asynchronous replication.
	for i := range s.repo.replicas.Pools {
		i := i
		go func() {
			// Random delay between 0.5s and 3s to simulate replication lag.
			delay := 500*time.Millisecond + time.Duration(rand.Intn(2500))*time.Millisecond
			time.Sleep(delay)
			_ = s.repo.CopyToReplica(context.Background(), i, loc)
		}()
	}

	return loc, nil
}

func (s *Service) ListFromMaster(ctx context.Context) ([]Location, error) {
	return s.repo.ListFromMaster(ctx)
}

func (s *Service) ListFromReplica(ctx context.Context, replicaIndex int) ([]Location, error) {
	return s.repo.ListFromReplica(ctx, replicaIndex)
}




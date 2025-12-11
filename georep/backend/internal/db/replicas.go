package db

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

// ReplicaSet holds connections to all replica databases.
// We simulate replication in application code, not at the Postgres level.
type ReplicaSet struct {
	Pools []*pgxpool.Pool
}

// NewReplicas connects to the given replica DSNs.
func NewReplicas(dsns []string) (*ReplicaSet, error) {
	pools := make([]*pgxpool.Pool, 0, len(dsns))
	for i, dsn := range dsns {
		cfg, err := pgxpool.ParseConfig(dsn)
		if err != nil {
			return nil, fmt.Errorf("parse replica %d dsn: %w", i+1, err)
		}
		pool, err := pgxpool.NewWithConfig(context.Background(), cfg)
		if err != nil {
			return nil, fmt.Errorf("connect replica %d: %w", i+1, err)
		}
		pools = append(pools, pool)
	}
	return &ReplicaSet{Pools: pools}, nil
}

// Close closes all replica pools.
func (r *ReplicaSet) Close() {
	for _, p := range r.Pools {
		p.Close()
	}
}




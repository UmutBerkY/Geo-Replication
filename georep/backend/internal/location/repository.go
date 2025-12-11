package location

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"

	"geo-repl-demo/internal/db"
)

// Repository coordinates reads and writes across master and replicas.
type Repository struct {
	master   *db.Master
	replicas *db.ReplicaSet
}

func NewRepository(master *db.Master, replicas *db.ReplicaSet) *Repository {
	return &Repository{
		master:   master,
		replicas: replicas,
	}
}

// InsertToMaster writes a new location to the master database.
func (r *Repository) InsertToMaster(ctx context.Context, in CreateLocationInput) (Location, error) {
	var loc Location
	err := r.master.Pool.QueryRow(ctx,
		`INSERT INTO locations (city, lat, lon, updated_at)
		 VALUES ($1, $2, $3, NOW())
		 RETURNING id, city, lat, lon, updated_at`,
		in.City, in.Lat, in.Lon,
	).Scan(&loc.ID, &loc.City, &loc.Lat, &loc.Lon, &loc.UpdatedAt)
	if err != nil {
		return Location{}, fmt.Errorf("insert master: %w", err)
	}
	return loc, nil
}

// CopyToReplica inserts a location row into a specific replica.
// This is called asynchronously to simulate eventual consistency.
func (r *Repository) CopyToReplica(ctx context.Context, replicaIndex int, loc Location) error {
	if replicaIndex < 0 || replicaIndex >= len(r.replicas.Pools) {
		return fmt.Errorf("invalid replica index %d", replicaIndex)
	}
	pool := r.replicas.Pools[replicaIndex]

	_, err := pool.Exec(ctx,
		`INSERT INTO locations (id, city, lat, lon, updated_at)
		 VALUES ($1, $2, $3, $4, $5)
		 ON CONFLICT (id) DO UPDATE
		 SET city = EXCLUDED.city,
		     lat = EXCLUDED.lat,
		     lon = EXCLUDED.lon,
		     updated_at = EXCLUDED.updated_at`,
		loc.ID, loc.City, loc.Lat, loc.Lon, loc.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("copy to replica %d: %w", replicaIndex+1, err)
	}
	return nil
}

// ListFromMaster returns all locations from the master.
func (r *Repository) ListFromMaster(ctx context.Context) ([]Location, error) {
	return queryLocations(ctx, r.master.Pool)
}

// ListFromReplica returns all locations from a given replica.
func (r *Repository) ListFromReplica(ctx context.Context, replicaIndex int) ([]Location, error) {
	if replicaIndex < 0 || replicaIndex >= len(r.replicas.Pools) {
		return nil, fmt.Errorf("invalid replica index %d", replicaIndex)
	}
	return queryLocations(ctx, r.replicas.Pools[replicaIndex])
}

// queryLocations is a helper shared by master/replica queries.
func queryLocations(ctx context.Context, pool *pgxpool.Pool) ([]Location, error) {
	rows, err := pool.Query(ctx,
		`SELECT id, city, lat, lon, updated_at
		 FROM locations
		 ORDER BY id`,
	)
	if err != nil {
		return nil, fmt.Errorf("query locations: %w", err)
	}
	defer rows.Close()

	var result []Location
	for rows.Next() {
		var loc Location
		if err := rows.Scan(&loc.ID, &loc.City, &loc.Lat, &loc.Lon, &loc.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan location: %w", err)
		}
		result = append(result, loc)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}
	return result, nil
}




package replication

import (
	"context"
	"geo-repl-demo/internal/model"
)

// Article replikalara yazılırken kullanılacak arabirim.
// Böylece replication paketi article paketine bağımlı olmaz.
type ReplicaWriter interface {
	CopyToReplica(ctx context.Context, replicaIndex int, a model.Article) error
}

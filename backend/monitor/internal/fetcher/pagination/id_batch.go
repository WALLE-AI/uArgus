package pagination

import (
	"context"

	"github.com/WALL-AI/uArgus/backend/monitor/internal/fetcher/pool"
)

// IDBatchFetcher splits a list of IDs into batches and fetches them concurrently.
type IDBatchFetcher[T any] struct {
	BatchSize   int
	Concurrency int
}

// NewIDBatchFetcher creates an IDBatchFetcher with the given batch size and concurrency.
func NewIDBatchFetcher[T any](batchSize, concurrency int) *IDBatchFetcher[T] {
	if batchSize <= 0 {
		batchSize = 50
	}
	if concurrency <= 0 {
		concurrency = 5
	}
	return &IDBatchFetcher[T]{BatchSize: batchSize, Concurrency: concurrency}
}

// FetchByIDs splits ids into batches and calls fn for each batch concurrently.
func (f *IDBatchFetcher[T]) FetchByIDs(ctx context.Context, ids []string, fn func(ctx context.Context, batch []string) ([]T, error)) ([]T, error) {
	batches := splitBatches(ids, f.BatchSize)

	results := pool.BoundedPool(ctx, f.Concurrency, batches, func(ctx context.Context, batch []string) ([]T, error) {
		return fn(ctx, batch)
	})

	var all []T
	for _, r := range results {
		if r.Err != nil {
			continue // skip failed batches (allSettled)
		}
		all = append(all, r.Value...)
	}
	return all, nil
}

func splitBatches(ids []string, size int) [][]string {
	var batches [][]string
	for i := 0; i < len(ids); i += size {
		end := i + size
		if end > len(ids) {
			end = len(ids)
		}
		batches = append(batches, ids[i:end])
	}
	return batches
}

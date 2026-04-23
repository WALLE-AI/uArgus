package tracking

import "context"

// TrackInfo holds the tracking state for a single item.
type TrackInfo struct {
	Hash            string
	FirstSeen       int64
	LastSeen        int64
	Count           int
	Stage           string // "BREAKING" | "DEVELOPING" | "SUSTAINED" | "FADING"
	Corroboration   int
}

// Trackable is implemented by items that can be tracked.
type Trackable interface {
	TrackHash() string
}

// StatefulTracker reads and writes tracking state to Redis.
type StatefulTracker interface {
	Read(ctx context.Context, keys []string) (map[string]TrackInfo, error)
	Write(ctx context.Context, items []Trackable) error
}

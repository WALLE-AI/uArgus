package tracking

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/WALL-AI/uArgus/backend/monitor/internal/cache"
)

const (
	storyKeyPrefix = "story-track:"
	storyTTL       = 172800 * time.Second // 48h
)

// StoryTracker tracks news story lifecycle: BREAKING → DEVELOPING → SUSTAINED → FADING.
// Uses Redis Pipeline for batch reads and writes to minimise round-trips.
type StoryTracker struct {
	rdb cache.Client
}

// NewStoryTracker creates a StoryTracker.
func NewStoryTracker(rdb cache.Client) *StoryTracker {
	return &StoryTracker{rdb: rdb}
}

// Read retrieves tracking info for a set of hashes using a single Pipeline GET.
func (t *StoryTracker) Read(ctx context.Context, hashes []string) (map[string]TrackInfo, error) {
	if len(hashes) == 0 {
		return nil, nil
	}
	keys := make([]string, len(hashes))
	for i, h := range hashes {
		keys[i] = storyKeyPrefix + h
	}
	raw, err := cache.PipelineGet(ctx, t.rdb, keys)
	if err != nil {
		return nil, fmt.Errorf("story tracker read: %w", err)
	}
	result := make(map[string]TrackInfo, len(raw))
	for i, h := range hashes {
		data, ok := raw[keys[i]]
		if !ok {
			continue
		}
		var info TrackInfo
		if err := json.Unmarshal(unquoteJSON(data), &info); err == nil {
			result[h] = info
		}
	}
	return result, nil
}

// Write updates tracking state for items using batch Pipeline read + write.
func (t *StoryTracker) Write(ctx context.Context, items []Trackable) error {
	if len(items) == 0 {
		return nil
	}
	now := time.Now().UnixMilli()

	// 1. batch read existing state
	hashes := make([]string, len(items))
	keys := make([]string, len(items))
	for i, item := range items {
		hashes[i] = item.TrackHash()
		keys[i] = storyKeyPrefix + hashes[i]
	}
	existing, err := cache.PipelineGet(ctx, t.rdb, keys)
	if err != nil {
		return fmt.Errorf("story tracker write (read phase): %w", err)
	}

	// 2. compute updated state
	writes := make(map[string][]byte, len(items))
	for i, item := range items {
		h := item.TrackHash()
		key := keys[i]

		var info TrackInfo
		if raw, ok := existing[key]; ok {
			_ = json.Unmarshal(unquoteJSON(raw), &info)
		}

		if info.Hash == "" {
			info = TrackInfo{
				Hash:      h,
				FirstSeen: now,
				LastSeen:  now,
				Count:     1,
				Stage:     "BREAKING",
			}
		} else {
			info.LastSeen = now
			info.Count++
			info.Stage = deriveStage(info)
		}

		data, err := json.Marshal(info)
		if err != nil {
			return fmt.Errorf("story tracker write: marshal %s: %w", h, err)
		}
		writes[key] = data
	}

	// 3. batch write
	return cache.PipelineSet(ctx, t.rdb, writes, storyTTL)
}

// unquoteJSON strips a surrounding JSON string wrapper that Upstash REST returns.
func unquoteJSON(b []byte) []byte {
	if len(b) >= 2 && b[0] == '"' {
		var s string
		if err := json.Unmarshal(b, &s); err == nil {
			return []byte(s)
		}
	}
	return b
}

// ComputeCorroboration counts how many items share the same hash.
func (t *StoryTracker) ComputeCorroboration(items []Trackable) map[string]int {
	counts := make(map[string]int)
	for _, item := range items {
		counts[item.TrackHash()]++
	}
	return counts
}

func deriveStage(info TrackInfo) string {
	age := time.Duration(info.LastSeen-info.FirstSeen) * time.Millisecond
	switch {
	case info.Count <= 1:
		return "BREAKING"
	case age < 2*time.Hour:
		return "DEVELOPING"
	case age < 24*time.Hour:
		return "SUSTAINED"
	default:
		return "FADING"
	}
}

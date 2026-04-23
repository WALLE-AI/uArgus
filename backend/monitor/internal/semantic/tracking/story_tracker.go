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
type StoryTracker struct {
	rdb cache.Client
}

// NewStoryTracker creates a StoryTracker.
func NewStoryTracker(rdb cache.Client) *StoryTracker {
	return &StoryTracker{rdb: rdb}
}

// Read retrieves tracking info for a set of hashes.
func (t *StoryTracker) Read(ctx context.Context, hashes []string) (map[string]TrackInfo, error) {
	result := make(map[string]TrackInfo, len(hashes))
	for _, h := range hashes {
		key := storyKeyPrefix + h
		data, err := t.rdb.Get(ctx, key)
		if err != nil || data == nil {
			continue
		}
		var info TrackInfo
		if err := json.Unmarshal(data, &info); err == nil {
			result[h] = info
		}
	}
	return result, nil
}

// Write updates tracking state for items.
func (t *StoryTracker) Write(ctx context.Context, items []Trackable) error {
	now := time.Now().UnixMilli()
	for _, item := range items {
		h := item.TrackHash()
		key := storyKeyPrefix + h

		existing, err := t.rdb.Get(ctx, key)
		var info TrackInfo
		if err == nil && existing != nil {
			_ = json.Unmarshal(existing, &info)
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
			return fmt.Errorf("story tracker write: %w", err)
		}
		if err := t.rdb.Set(ctx, key, data, storyTTL); err != nil {
			return fmt.Errorf("story tracker write: %w", err)
		}
	}
	return nil
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

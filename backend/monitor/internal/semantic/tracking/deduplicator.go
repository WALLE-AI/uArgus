package tracking

import (
	"crypto/sha256"
	"fmt"
)

// HashDedup provides generic hash-based deduplication.
type HashDedup struct{}

// Dedup removes duplicate items based on TrackHash().
func (HashDedup) Dedup(items []Trackable) []Trackable {
	seen := make(map[string]bool)
	var result []Trackable
	for _, item := range items {
		h := item.TrackHash()
		if seen[h] {
			continue
		}
		seen[h] = true
		result = append(result, item)
	}
	return result
}

// Hash computes a SHA-256 hex hash for dedup purposes.
func Hash(s string) string {
	h := sha256.Sum256([]byte(s))
	return fmt.Sprintf("%x", h[:16]) // 32 hex chars
}

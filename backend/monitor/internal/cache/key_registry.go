package cache

import (
	"fmt"
	"time"

	"github.com/WALL-AI/uArgus/backend/monitor/internal/registry"
)

// KeyEntry holds the full metadata for one Redis key, derived from SourceSpec.
type KeyEntry struct {
	Key         string
	SeedMetaKey string        // "seed-meta:{domain}:{resource}"
	DataTTL     time.Duration // from SourceSpec.DataTTL
	MaxStale    time.Duration // from SourceSpec.MaxStaleDuration
	Tier        string        // "fast" | "slow"
	Category    string        // "bootstrap" | "standalone" | "on-demand" | "derived"
}

// KeyRegistry centralises all Redis key definitions.
type KeyRegistry struct {
	entries map[string]KeyEntry
}

// NewKeyRegistryFromSpecs builds a KeyRegistry from a slice of SourceSpecs.
func NewKeyRegistryFromSpecs(specs []registry.SourceSpec) *KeyRegistry {
	r := &KeyRegistry{entries: make(map[string]KeyEntry, len(specs))}
	for _, sp := range specs {
		category := deriveCategory(sp)
		tier := "slow"
		if sp.DataTTL <= 10*time.Minute {
			tier = "fast"
		}

		entry := KeyEntry{
			Key:         sp.CanonicalKey,
			SeedMetaKey: fmt.Sprintf("seed-meta:%s:%s", sp.Domain, sp.Resource),
			DataTTL:     sp.DataTTL,
			MaxStale:    sp.MaxStaleDuration,
			Tier:        tier,
			Category:    category,
		}
		r.entries[sp.CanonicalKey] = entry

		// also register extra (bootstrap mirror) keys
		for _, ek := range sp.ExtraKeys {
			mirror := entry
			mirror.Key = ek
			mirror.Category = "bootstrap"
			r.entries[ek] = mirror
		}
	}
	return r
}

// Get retrieves a KeyEntry by canonical key.
func (r *KeyRegistry) Get(key string) (KeyEntry, bool) {
	e, ok := r.entries[key]
	return e, ok
}

// BootstrapKeys returns all keys categorised as "bootstrap" — replaces v1 BOOTSTRAP_CACHE_KEYS.
func (r *KeyRegistry) BootstrapKeys() map[string]string {
	out := make(map[string]string)
	for k, e := range r.entries {
		if e.Category == "bootstrap" || e.Category == "standalone" {
			out[k] = e.Tier
		}
	}
	return out
}

// SeedMetaEntries returns all seed-meta entries — replaces v1 SEED_META.
func (r *KeyRegistry) SeedMetaEntries() map[string]KeyEntry {
	out := make(map[string]KeyEntry)
	for _, e := range r.entries {
		if e.SeedMetaKey != "" {
			out[e.SeedMetaKey] = e
		}
	}
	return out
}

// Validate checks TTL consistency and returns the first error found.
func (r *KeyRegistry) Validate() error {
	for k, e := range r.entries {
		if e.DataTTL <= 0 {
			return fmt.Errorf("key_registry: %s has non-positive DataTTL", k)
		}
		if e.MaxStale > 0 && e.MaxStale < e.DataTTL {
			return fmt.Errorf("key_registry: %s MaxStale (%v) < DataTTL (%v)", k, e.MaxStale, e.DataTTL)
		}
	}
	return nil
}

// ── helpers ─────────────────────────────────────────────────

func deriveCategory(sp registry.SourceSpec) string {
	switch sp.Schedule.(type) {
	case registry.OnDemandSchedule:
		return "on-demand"
	default:
		if len(sp.ExtraKeys) > 0 {
			return "bootstrap"
		}
		return "standalone"
	}
}

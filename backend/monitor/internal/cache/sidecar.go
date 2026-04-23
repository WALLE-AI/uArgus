package cache

import (
	"sync"
	"time"
)

// SidecarCache is an in-process LRU+TTL cache for Tauri sidecar mode.
type SidecarCache struct {
	mu         sync.Mutex
	maxEntries int
	maxBytes   int64
	sweepEvery time.Duration

	entries map[string]*sidecarEntry
	order   []string // oldest first
	size    int64
	stopCh  chan struct{}
}

type sidecarEntry struct {
	value     []byte
	expiresAt time.Time
}

// SidecarStats reports cache statistics.
type SidecarStats struct {
	Entries  int
	Bytes    int64
	Capacity int
}

// NewSidecarCache creates a SidecarCache and starts the sweep goroutine.
func NewSidecarCache(maxEntries int, maxBytes int64, sweepEvery time.Duration) *SidecarCache {
	if maxEntries <= 0 {
		maxEntries = 500
	}
	if maxBytes <= 0 {
		maxBytes = 50 * 1024 * 1024
	}
	if sweepEvery <= 0 {
		sweepEvery = 60 * time.Second
	}

	sc := &SidecarCache{
		maxEntries: maxEntries,
		maxBytes:   maxBytes,
		sweepEvery: sweepEvery,
		entries:    make(map[string]*sidecarEntry),
		stopCh:     make(chan struct{}),
	}
	go sc.sweepLoop()
	return sc
}

// Get returns the value for key, or nil if missing/expired.
func (sc *SidecarCache) Get(key string) ([]byte, bool) {
	sc.mu.Lock()
	defer sc.mu.Unlock()

	e, ok := sc.entries[key]
	if !ok {
		return nil, false
	}
	if time.Now().After(e.expiresAt) {
		sc.removeLocked(key)
		return nil, false
	}
	return e.value, true
}

// Set stores a value with the given TTL.
func (sc *SidecarCache) Set(key string, value []byte, ttl time.Duration) {
	sc.mu.Lock()
	defer sc.mu.Unlock()

	// remove old if exists
	if old, ok := sc.entries[key]; ok {
		sc.size -= int64(len(old.value))
		delete(sc.entries, key)
	} else {
		sc.order = append(sc.order, key)
	}

	sc.entries[key] = &sidecarEntry{
		value:     value,
		expiresAt: time.Now().Add(ttl),
	}
	sc.size += int64(len(value))

	// evict if over capacity
	for len(sc.entries) > sc.maxEntries || sc.size > sc.maxBytes {
		if len(sc.order) == 0 {
			break
		}
		sc.removeLocked(sc.order[0])
	}
}

// Stats returns current cache statistics.
func (sc *SidecarCache) Stats() SidecarStats {
	sc.mu.Lock()
	defer sc.mu.Unlock()
	return SidecarStats{
		Entries:  len(sc.entries),
		Bytes:    sc.size,
		Capacity: sc.maxEntries,
	}
}

// Stop halts the sweep goroutine.
func (sc *SidecarCache) Stop() {
	close(sc.stopCh)
}

func (sc *SidecarCache) removeLocked(key string) {
	if e, ok := sc.entries[key]; ok {
		sc.size -= int64(len(e.value))
		delete(sc.entries, key)
	}
	// remove from order
	for i, k := range sc.order {
		if k == key {
			sc.order = append(sc.order[:i], sc.order[i+1:]...)
			break
		}
	}
}

func (sc *SidecarCache) sweepLoop() {
	ticker := time.NewTicker(sc.sweepEvery)
	defer ticker.Stop()
	for {
		select {
		case <-sc.stopCh:
			return
		case <-ticker.C:
			sc.sweep()
		}
	}
}

func (sc *SidecarCache) sweep() {
	sc.mu.Lock()
	defer sc.mu.Unlock()
	now := time.Now()
	for key, e := range sc.entries {
		if now.After(e.expiresAt) {
			sc.removeLocked(key)
		}
	}
}

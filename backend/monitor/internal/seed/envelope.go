package seed

import (
	"encoding/json"
	"fmt"
)

// SeedMeta is the single-source-of-truth for seed envelope metadata.
// Replaces v1's three-file manual sync (_seed-envelope-source.mjs + seed-envelope.ts + api/_seed-envelope.js).
type SeedMeta struct {
	FetchedAt      int64    `json:"fetchedAt"`
	RecordCount    int      `json:"recordCount"`
	SourceVersion  string   `json:"sourceVersion"`
	SchemaVersion  int      `json:"schemaVersion"`
	State          string   `json:"state"`                    // "OK" | "OK_ZERO" | "ERROR"
	FailedDatasets []string `json:"failedDatasets,omitempty"`
	ErrorReason    string   `json:"errorReason,omitempty"`
	GroupID        string   `json:"groupId,omitempty"`
}

// SeedEnvelope wraps data with _seed metadata.
type SeedEnvelope struct {
	Seed *SeedMeta `json:"_seed"`
	Data any       `json:"data"`
}

// Unwrap extracts metadata and data from raw bytes.
// Supports both legacy (no _seed field) and contract (has _seed) formats.
func Unwrap(raw []byte) (*SeedMeta, json.RawMessage, error) {
	var probe struct {
		Seed *SeedMeta       `json:"_seed"`
		Data json.RawMessage `json:"data"`
	}
	if err := json.Unmarshal(raw, &probe); err != nil {
		return nil, nil, fmt.Errorf("seed: unwrap: %w", err)
	}

	// contract format: has _seed field
	if probe.Seed != nil {
		if len(probe.Data) == 0 {
			return probe.Seed, nil, nil
		}
		return probe.Seed, probe.Data, nil
	}

	// legacy format: entire payload is data (no _seed wrapper)
	return nil, raw, nil
}

// Build constructs a SeedEnvelope and marshals it to JSON.
func Build(meta SeedMeta, data any) ([]byte, error) {
	env := SeedEnvelope{Seed: &meta, Data: data}
	b, err := json.Marshal(env)
	if err != nil {
		return nil, fmt.Errorf("seed: build envelope: %w", err)
	}
	return b, nil
}

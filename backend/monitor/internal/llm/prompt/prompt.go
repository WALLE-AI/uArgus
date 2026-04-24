package prompt

// SystemPromptBlock is a cache-aware system prompt segment.
type SystemPromptBlock struct {
	Text       string `json:"text"`
	CacheScope string `json:"cache_scope,omitempty"` // "" | "ephemeral"
}

// CacheControl is the Anthropic cache control hint.
type CacheControl struct {
	Type string `json:"type"` // "ephemeral"
}

// ToCacheControl converts a cache scope string to a CacheControl pointer.
// Returns nil when no caching is requested.
func ToCacheControl(scope string) *CacheControl {
	switch scope {
	case "ephemeral":
		return &CacheControl{Type: "ephemeral"}
	default:
		return nil
	}
}

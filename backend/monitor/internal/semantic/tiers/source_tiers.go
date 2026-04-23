package tiers

// tierMap holds the mapping of source names to tier levels (1-4).
// Tier 1 = most authoritative, Tier 4 = least.
var tierMap = map[string]int{
	// Tier 1: Major wire services & state media
	"reuters":     1,
	"ap":          1,
	"afp":         1,
	"bbc":         1,
	"nyt":         1,
	"washpost":    1,
	"guardian":    1,
	"ft":          1,
	"economist":   1,
	"wsj":         1,
	"bloomberg":   1,

	// Tier 2: Major national/regional outlets
	"cnn":          2,
	"aljazeera":    2,
	"dw":           2,
	"france24":     2,
	"scmp":         2,
	"nikkei":       2,
	"times":        2,
	"telegraph":    2,
	"politico":     2,
	"axios":        2,
	"npr":          2,

	// Tier 3: Specialised / digital-first
	"techcrunch":   3,
	"arstechnica":  3,
	"wired":        3,
	"theverge":     3,
	"hackernews":   3,
	"defenseone":   3,
	"thehill":      3,
	"foreignpolicy": 3,
	"brookings":    3,

	// Tier 4: Aggregators & lesser-known
	"newsapi":      4,
	"googlenews":   4,
}

// GetTier returns the tier (1-4) for a source name. Unknown sources default to 4.
func GetTier(sourceName string) int {
	if t, ok := tierMap[sourceName]; ok {
		return t
	}
	return 4
}

// AllTiers returns a copy of the full tier mapping.
func AllTiers() map[string]int {
	out := make(map[string]int, len(tierMap))
	for k, v := range tierMap {
		out[k] = v
	}
	return out
}

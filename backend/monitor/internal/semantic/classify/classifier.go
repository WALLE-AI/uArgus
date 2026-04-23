package classify

import (
	"strings"
)

// ClassificationResult holds the output of Classify().
type ClassificationResult struct {
	Categories []string
	Severity   int // 0-7
	Matched    []string
}

// ClassifyOption is a functional option for Classify.
type ClassifyOption func(*classifyConfig)

type classifyConfig struct {
	variant string
}

// WithVariant sets the digest variant for classification.
func WithVariant(v string) ClassifyOption {
	return func(c *classifyConfig) { c.variant = v }
}

// KeywordLevel defines one tier of keywords with a severity weight.
type KeywordLevel struct {
	Name     string
	Severity int
	Keywords []string
}

// Classifier performs keyword-based classification on text.
type Classifier struct {
	levels     []KeywordLevel
	exclusions map[string]bool
	shortWords map[string]bool
}

// NewNewsClassifier returns a Classifier preset with 7 geo+tech keyword layers.
func NewNewsClassifier() *Classifier {
	return &Classifier{
		levels:     defaultKeywordLevels(),
		exclusions: defaultExclusions(),
		shortWords: defaultShortWords(),
	}
}

// NewCustomClassifier creates a Classifier with caller-defined levels.
func NewCustomClassifier(levels []KeywordLevel) *Classifier {
	return &Classifier{
		levels:     levels,
		exclusions: defaultExclusions(),
		shortWords: defaultShortWords(),
	}
}

// Classify matches a title against keyword levels and returns categories + severity.
func (c *Classifier) Classify(title string, opts ...ClassifyOption) ClassificationResult {
	cfg := &classifyConfig{}
	for _, o := range opts {
		o(cfg)
	}

	lower := strings.ToLower(title)
	words := strings.Fields(lower)

	result := ClassificationResult{}
	seen := make(map[string]bool)

	for _, w := range words {
		w = strings.Trim(w, ".,;:!?\"'()-")
		if c.exclusions[w] || (len(w) <= 2 && !c.shortWords[w]) {
			continue
		}

		for _, level := range c.levels {
			for _, kw := range level.Keywords {
				if strings.Contains(lower, strings.ToLower(kw)) && !seen[kw] {
					seen[kw] = true
					result.Matched = append(result.Matched, kw)
					if !containsStr(result.Categories, level.Name) {
						result.Categories = append(result.Categories, level.Name)
					}
					if level.Severity > result.Severity {
						result.Severity = level.Severity
					}
				}
			}
		}
	}

	return result
}

func containsStr(ss []string, s string) bool {
	for _, v := range ss {
		if v == s {
			return true
		}
	}
	return false
}

package classify

// defaultExclusions returns common words that should not trigger classification.
func defaultExclusions() map[string]bool {
	words := []string{
		"the", "a", "an", "and", "or", "but", "in", "on", "at", "to",
		"for", "of", "with", "by", "from", "is", "are", "was", "were",
		"be", "been", "being", "have", "has", "had", "do", "does", "did",
		"will", "would", "could", "should", "may", "might", "can",
		"this", "that", "these", "those", "it", "its", "he", "she",
		"they", "them", "their", "we", "our", "you", "your",
		"not", "no", "so", "if", "then", "than", "as", "up",
		"new", "says", "said", "report", "reports", "news",
	}
	m := make(map[string]bool, len(words))
	for _, w := range words {
		m[w] = true
	}
	return m
}

// defaultShortWords returns short words (≤2 chars) that ARE valid keywords.
func defaultShortWords() map[string]bool {
	return map[string]bool{
		"ai": true,
		"uk": true,
		"us": true,
		"eu": true,
		"un": true,
	}
}

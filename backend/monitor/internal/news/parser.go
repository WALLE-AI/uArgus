package news

import (
	"crypto/sha256"
	"fmt"
	"strings"
	"time"

	xmlparser "github.com/WALL-AI/uArgus/backend/monitor/internal/fetcher/parser"
)

// ParsedItem is the normalised representation of one news item.
type ParsedItem struct {
	Title       string   `json:"title"`
	Link        string   `json:"link"`
	Description string   `json:"description,omitempty"`
	PubDate     string   `json:"pubDate"`
	Source      string   `json:"source"`
	Hash        string   `json:"hash"`
	Severity    int      `json:"severity,omitempty"`
	Score       float64  `json:"importanceScore,omitempty"`
	Categories  []string `json:"categories,omitempty"`

	// tracking fields (populated later)
	Stage         string `json:"stage,omitempty"`
	Corroboration int    `json:"corroboration,omitempty"`
	Tier          int    `json:"tier,omitempty"`
}

// TrackHash implements tracking.Trackable.
func (p ParsedItem) TrackHash() string { return p.Hash }

// ParseFeedResponse parses RSS body and normalises to ParsedItems.
func ParseFeedResponse(body []byte, feed FeedEntry) []ParsedItem {
	rssItems, err := xmlparser.ParseRssXML(body)
	if err == nil && len(rssItems) > 0 {
		return rssToItems(rssItems, feed)
	}
	// try atom
	atomEntries, atomErr := xmlparser.ParseAtomXML(body)
	if atomErr == nil && len(atomEntries) > 0 {
		return atomToItems(atomEntries, feed)
	}
	return nil
}

func rssToItems(items []xmlparser.RssItem, feed FeedEntry) []ParsedItem {
	result := make([]ParsedItem, 0, len(items))
	for _, it := range items {
		if it.Title == "" || it.Link == "" {
			continue
		}
		source := it.Source
		if source == "" {
			source = feed.Name
		}
		result = append(result, ParsedItem{
			Title:       it.Title,
			Link:        normalizeLink(it.Link),
			Description: truncate(it.Description, 500),
			PubDate:     normalizeDate(it.PubDate),
			Source:      source,
			Hash:        hashItem(it.Title, it.Link),
		})
	}
	return result
}

func atomToItems(entries []xmlparser.AtomEntry, feed FeedEntry) []ParsedItem {
	result := make([]ParsedItem, 0, len(entries))
	for _, e := range entries {
		if e.Title == "" {
			continue
		}
		result = append(result, ParsedItem{
			Title:       e.Title,
			Link:        e.Link,
			Description: truncate(e.Summary, 500),
			PubDate:     normalizeDate(e.Published),
			Source:      feed.Name,
			Hash:        hashItem(e.Title, e.Link),
		})
	}
	return result
}

func hashItem(title, link string) string {
	h := sha256.Sum256([]byte(strings.ToLower(title) + "|" + link))
	return fmt.Sprintf("%x", h[:16])
}

func normalizeLink(link string) string {
	return strings.TrimSpace(link)
}

func normalizeDate(s string) string {
	s = strings.TrimSpace(s)
	if s == "" {
		return time.Now().UTC().Format(time.RFC3339)
	}
	formats := []string{
		time.RFC3339,
		time.RFC1123,
		time.RFC1123Z,
		"Mon, 02 Jan 2006 15:04:05 MST",
		"Mon, 02 Jan 2006 15:04:05 -0700",
		"2006-01-02T15:04:05Z",
	}
	for _, f := range formats {
		if t, err := time.Parse(f, s); err == nil {
			return t.UTC().Format(time.RFC3339)
		}
	}
	return s
}

func truncate(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max] + "…"
}

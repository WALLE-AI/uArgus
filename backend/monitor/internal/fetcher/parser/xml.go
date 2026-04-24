package parser

import (
	"encoding/xml"
	"fmt"
	"regexp"
	"strings"
)

// ── RSS ─────────────────────────────────────────────────────

// RssItem represents one item from an RSS feed.
type RssItem struct {
	Title       string `json:"title"`
	Link        string `json:"link"`
	Description string `json:"description"`
	PubDate     string `json:"pubDate"`
	Source      string `json:"source"`
	GUID        string `json:"guid"`
}

type rssChannel struct {
	XMLName xml.Name  `xml:"channel"`
	Title   string    `xml:"title"`
	Items   []rssItem `xml:"item"`
}

type rssItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
	GUID        string `xml:"guid"`
	Source      string `xml:"source"`
}

type rssFeed struct {
	Channel rssChannel `xml:"channel"`
}

// ParseRssXML parses RSS 2.0 XML and returns structured items.
func ParseRssXML(body []byte) ([]RssItem, error) {
	var feed rssFeed
	if err := xml.Unmarshal(body, &feed); err != nil {
		return nil, fmt.Errorf("parser: rss xml: %w", err)
	}

	items := make([]RssItem, 0, len(feed.Channel.Items))
	for _, it := range feed.Channel.Items {
		items = append(items, RssItem{
			Title:       cleanText(it.Title),
			Link:        strings.TrimSpace(it.Link),
			Description: cleanText(it.Description),
			PubDate:     it.PubDate,
			Source:      it.Source,
			GUID:        it.GUID,
		})
	}
	return items, nil
}

// ── Atom ────────────────────────────────────────────────────

// AtomEntry represents one entry from an Atom feed.
type AtomEntry struct {
	Title     string `json:"title"`
	Link      string `json:"link"`
	Summary   string `json:"summary"`
	Published string `json:"published"`
	ID        string `json:"id"`
}

type atomFeed struct {
	XMLName xml.Name    `xml:"feed"`
	Entries []atomEntry `xml:"entry"`
}

type atomEntry struct {
	Title     string     `xml:"title"`
	Links     []atomLink `xml:"link"`
	Summary   string     `xml:"summary"`
	Published string     `xml:"published"`
	Updated   string     `xml:"updated"`
	ID        string     `xml:"id"`
}

type atomLink struct {
	Href string `xml:"href,attr"`
	Rel  string `xml:"rel,attr"`
	Type string `xml:"type,attr"`
}

// ParseAtomXML parses Atom XML and returns structured entries.
func ParseAtomXML(body []byte) ([]AtomEntry, error) {
	var feed atomFeed
	if err := xml.Unmarshal(body, &feed); err != nil {
		return nil, fmt.Errorf("parser: atom xml: %w", err)
	}

	entries := make([]AtomEntry, 0, len(feed.Entries))
	for _, e := range feed.Entries {
		link := ""
		for _, l := range e.Links {
			if l.Rel == "alternate" || l.Rel == "" {
				link = l.Href
				break
			}
		}
		pub := e.Published
		if pub == "" {
			pub = e.Updated
		}
		entries = append(entries, AtomEntry{
			Title:     cleanText(e.Title),
			Link:      link,
			Summary:   cleanText(e.Summary),
			Published: pub,
			ID:        e.ID,
		})
	}
	return entries, nil
}

// ── helpers ─────────────────────────────────────────────────

var (
	cdataRe    = regexp.MustCompile(`<!\[CDATA\[(.+?)\]\]>`)
	entityRe   = regexp.MustCompile(`&[a-zA-Z]+;`)
	tagStripRe = regexp.MustCompile(`<[^>]*>`)
)

func cleanXML(b []byte) []byte {
	s := cdataRe.ReplaceAllString(string(b), "$1")
	return []byte(s)
}

func cleanText(s string) string {
	s = tagStripRe.ReplaceAllString(s, "")
	s = strings.ReplaceAll(s, "&amp;", "&")
	s = strings.ReplaceAll(s, "&lt;", "<")
	s = strings.ReplaceAll(s, "&gt;", ">")
	s = strings.ReplaceAll(s, "&quot;", "\"")
	s = strings.ReplaceAll(s, "&#39;", "'")
	return strings.TrimSpace(s)
}

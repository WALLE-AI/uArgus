package news

// FeedEntry describes a single RSS/Atom feed.
type FeedEntry struct {
	URL      string `json:"url"`
	Name     string `json:"name"`
	Category string `json:"category"`
	Lang     string `json:"lang"`
	Variant  string `json:"variant"` // "full","tech","geo","finance","quick"
}

// GetFeeds returns the feeds for a given variant and language.
func GetFeeds(variant, lang string) []FeedEntry {
	if lang == "" {
		lang = "en"
	}
	var result []FeedEntry
	for _, f := range allFeeds {
		if f.Variant == variant && f.Lang == lang {
			result = append(result, f)
		}
	}
	// fallback: if no exact variant match, return "full"
	if len(result) == 0 && variant != "full" {
		return GetFeeds("full", lang)
	}
	return result
}

// AllVariants returns all available variant names.
func AllVariants() []string {
	return []string{"full", "tech", "geo", "finance", "quick"}
}

// allFeeds is the master feed list — mirrors v1 _feeds.ts (437 lines).
// Condensed here; real deployment would embed a JSON file.
var allFeeds = []FeedEntry{
	// ── full (en) ───────────────────────────────────────────
	{URL: "https://feeds.reuters.com/reuters/topNews", Name: "reuters", Category: "general", Lang: "en", Variant: "full"},
	{URL: "https://feeds.bbci.co.uk/news/world/rss.xml", Name: "bbc-world", Category: "general", Lang: "en", Variant: "full"},
	{URL: "https://rss.nytimes.com/services/xml/rss/nyt/World.xml", Name: "nyt-world", Category: "general", Lang: "en", Variant: "full"},
	{URL: "https://feeds.washingtonpost.com/rss/world", Name: "washpost-world", Category: "general", Lang: "en", Variant: "full"},
	{URL: "https://www.theguardian.com/world/rss", Name: "guardian-world", Category: "general", Lang: "en", Variant: "full"},
	{URL: "https://www.aljazeera.com/xml/rss/all.xml", Name: "aljazeera", Category: "general", Lang: "en", Variant: "full"},
	{URL: "https://rss.dw.com/xml/rss-en-world", Name: "dw-world", Category: "general", Lang: "en", Variant: "full"},
	{URL: "https://www.france24.com/en/rss", Name: "france24", Category: "general", Lang: "en", Variant: "full"},
	{URL: "https://feeds.npr.org/1004/rss.xml", Name: "npr-world", Category: "general", Lang: "en", Variant: "full"},
	{URL: "https://www.scmp.com/rss/91/feed", Name: "scmp", Category: "asia", Lang: "en", Variant: "full"},
	{URL: "https://asia.nikkei.com/rss", Name: "nikkei-asia", Category: "asia", Lang: "en", Variant: "full"},
	{URL: "https://www.politico.com/rss/politico-top-news.xml", Name: "politico", Category: "politics", Lang: "en", Variant: "full"},
	{URL: "https://feeds.axios.com/feeds/feed.rss", Name: "axios", Category: "general", Lang: "en", Variant: "full"},
	{URL: "https://www.ft.com/?format=rss", Name: "ft", Category: "finance", Lang: "en", Variant: "full"},
	{URL: "https://feeds.bloomberg.com/economics/news.rss", Name: "bloomberg-econ", Category: "economics", Lang: "en", Variant: "full"},

	// ── tech (en) ───────────────────────────────────────────
	{URL: "https://feeds.feedburner.com/TechCrunch/", Name: "techcrunch", Category: "tech", Lang: "en", Variant: "tech"},
	{URL: "https://feeds.arstechnica.com/arstechnica/index", Name: "arstechnica", Category: "tech", Lang: "en", Variant: "tech"},
	{URL: "https://www.wired.com/feed/rss", Name: "wired", Category: "tech", Lang: "en", Variant: "tech"},
	{URL: "https://www.theverge.com/rss/index.xml", Name: "theverge", Category: "tech", Lang: "en", Variant: "tech"},
	{URL: "https://www.technologyreview.com/feed/", Name: "mit-tech-review", Category: "tech", Lang: "en", Variant: "tech"},

	// ── geo (en) ────────────────────────────────────────────
	{URL: "https://reliefweb.int/updates/rss.xml", Name: "reliefweb", Category: "humanitarian", Lang: "en", Variant: "geo"},
	{URL: "https://earthquake.usgs.gov/earthquakes/feed/v1.0/summary/significant_month.atom", Name: "usgs-quakes", Category: "seismic", Lang: "en", Variant: "geo"},
	{URL: "https://www.gdacs.org/xml/rss.xml", Name: "gdacs", Category: "disaster", Lang: "en", Variant: "geo"},

	// ── finance (en) ────────────────────────────────────────
	{URL: "https://feeds.bloomberg.com/markets/news.rss", Name: "bloomberg-markets", Category: "markets", Lang: "en", Variant: "finance"},
	{URL: "https://www.cnbc.com/id/100003114/device/rss/rss.html", Name: "cnbc", Category: "markets", Lang: "en", Variant: "finance"},
	{URL: "https://feeds.marketwatch.com/marketwatch/topstories/", Name: "marketwatch", Category: "markets", Lang: "en", Variant: "finance"},

	// ── quick (en) — fewer feeds for faster refresh ─────────
	{URL: "https://feeds.reuters.com/reuters/topNews", Name: "reuters", Category: "general", Lang: "en", Variant: "quick"},
	{URL: "https://feeds.bbci.co.uk/news/world/rss.xml", Name: "bbc-world", Category: "general", Lang: "en", Variant: "quick"},
	{URL: "https://rss.nytimes.com/services/xml/rss/nyt/World.xml", Name: "nyt-world", Category: "general", Lang: "en", Variant: "quick"},
}

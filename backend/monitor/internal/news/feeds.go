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
	// ══════════════════════════════════════════════════════════
	//  full (en) — comprehensive world news (~85 feeds)
	// ══════════════════════════════════════════════════════════

	// ── Wire services ───────────────────────────────────────
	{URL: "https://feeds.reuters.com/reuters/topNews", Name: "reuters", Category: "general", Lang: "en", Variant: "full"},
	{URL: "https://feeds.reuters.com/reuters/worldNews", Name: "reuters-world", Category: "general", Lang: "en", Variant: "full"},
	{URL: "https://feeds.reuters.com/reuters/businessNews", Name: "reuters-biz", Category: "business", Lang: "en", Variant: "full"},
	{URL: "https://feeds.reuters.com/reuters/technologyNews", Name: "reuters-tech", Category: "tech", Lang: "en", Variant: "full"},

	// ── UK / Europe ─────────────────────────────────────────
	{URL: "https://feeds.bbci.co.uk/news/world/rss.xml", Name: "bbc-world", Category: "general", Lang: "en", Variant: "full"},
	{URL: "https://feeds.bbci.co.uk/news/rss.xml", Name: "bbc-news", Category: "general", Lang: "en", Variant: "full"},
	{URL: "https://feeds.bbci.co.uk/news/technology/rss.xml", Name: "bbc-tech", Category: "tech", Lang: "en", Variant: "full"},
	{URL: "https://www.theguardian.com/world/rss", Name: "guardian-world", Category: "general", Lang: "en", Variant: "full"},
	{URL: "https://www.theguardian.com/uk/rss", Name: "guardian-uk", Category: "general", Lang: "en", Variant: "full"},
	{URL: "https://www.theguardian.com/environment/rss", Name: "guardian-env", Category: "environment", Lang: "en", Variant: "full"},
	{URL: "https://www.ft.com/?format=rss", Name: "ft", Category: "finance", Lang: "en", Variant: "full"},
	{URL: "https://www.telegraph.co.uk/rss.xml", Name: "telegraph", Category: "general", Lang: "en", Variant: "full"},
	{URL: "https://www.independent.co.uk/rss", Name: "independent", Category: "general", Lang: "en", Variant: "full"},
	{URL: "https://feeds.skynews.com/feeds/rss/world.xml", Name: "sky-world", Category: "general", Lang: "en", Variant: "full"},
	{URL: "https://www.euronews.com/rss", Name: "euronews", Category: "general", Lang: "en", Variant: "full"},
	{URL: "https://rss.dw.com/xml/rss-en-world", Name: "dw-world", Category: "general", Lang: "en", Variant: "full"},
	{URL: "https://rss.dw.com/xml/rss-en-top", Name: "dw-top", Category: "general", Lang: "en", Variant: "full"},
	{URL: "https://www.france24.com/en/rss", Name: "france24", Category: "general", Lang: "en", Variant: "full"},

	// ── US ──────────────────────────────────────────────────
	{URL: "https://rss.nytimes.com/services/xml/rss/nyt/World.xml", Name: "nyt-world", Category: "general", Lang: "en", Variant: "full"},
	{URL: "https://rss.nytimes.com/services/xml/rss/nyt/HomePage.xml", Name: "nyt-home", Category: "general", Lang: "en", Variant: "full"},
	{URL: "https://rss.nytimes.com/services/xml/rss/nyt/Technology.xml", Name: "nyt-tech", Category: "tech", Lang: "en", Variant: "full"},
	{URL: "https://feeds.washingtonpost.com/rss/world", Name: "washpost-world", Category: "general", Lang: "en", Variant: "full"},
	{URL: "https://feeds.washingtonpost.com/rss/national", Name: "washpost-national", Category: "general", Lang: "en", Variant: "full"},
	{URL: "http://rss.cnn.com/rss/edition_world.rss", Name: "cnn-world", Category: "general", Lang: "en", Variant: "full"},
	{URL: "http://rss.cnn.com/rss/edition.rss", Name: "cnn-top", Category: "general", Lang: "en", Variant: "full"},
	{URL: "https://feeds.npr.org/1004/rss.xml", Name: "npr-world", Category: "general", Lang: "en", Variant: "full"},
	{URL: "https://feeds.npr.org/1001/rss.xml", Name: "npr-news", Category: "general", Lang: "en", Variant: "full"},
	{URL: "https://feeds.nbcnews.com/nbcnews/public/world", Name: "nbcnews-world", Category: "general", Lang: "en", Variant: "full"},
	{URL: "https://www.cbsnews.com/latest/rss/world", Name: "cbsnews-world", Category: "general", Lang: "en", Variant: "full"},
	{URL: "https://abcnews.go.com/abcnews/internationalheadlines", Name: "abc-us-intl", Category: "general", Lang: "en", Variant: "full"},
	{URL: "https://www.politico.com/rss/politico-top-news.xml", Name: "politico", Category: "politics", Lang: "en", Variant: "full"},
	{URL: "https://feeds.axios.com/feeds/feed.rss", Name: "axios", Category: "general", Lang: "en", Variant: "full"},
	{URL: "https://thehill.com/feed/", Name: "thehill", Category: "politics", Lang: "en", Variant: "full"},
	{URL: "https://www.vox.com/rss/index.xml", Name: "vox", Category: "explainer", Lang: "en", Variant: "full"},
	{URL: "https://www.theatlantic.com/feed/all/", Name: "theatlantic", Category: "longform", Lang: "en", Variant: "full"},
	{URL: "https://feeds.feedburner.com/foreignpolicy/topnews", Name: "foreignpolicy", Category: "geopolitics", Lang: "en", Variant: "full"},

	// ── Middle East / Africa ────────────────────────────────
	{URL: "https://www.aljazeera.com/xml/rss/all.xml", Name: "aljazeera", Category: "general", Lang: "en", Variant: "full"},
	{URL: "https://www.middleeasteye.net/rss", Name: "middleeasteye", Category: "mideast", Lang: "en", Variant: "full"},
	{URL: "https://www.arabnews.com/rss.xml", Name: "arabnews", Category: "mideast", Lang: "en", Variant: "full"},
	{URL: "https://www.thenationalnews.com/rss", Name: "thenational-ae", Category: "mideast", Lang: "en", Variant: "full"},
	{URL: "https://www.dailymaverick.co.za/rss", Name: "dailymaverick", Category: "africa", Lang: "en", Variant: "full"},

	// ── Asia-Pacific ────────────────────────────────────────
	{URL: "https://www.scmp.com/rss/91/feed", Name: "scmp", Category: "asia", Lang: "en", Variant: "full"},
	{URL: "https://asia.nikkei.com/rss", Name: "nikkei-asia", Category: "asia", Lang: "en", Variant: "full"},
	{URL: "https://www.japantimes.co.jp/feed/", Name: "japantimes", Category: "asia", Lang: "en", Variant: "full"},
	{URL: "https://www.straitstimes.com/news/world/rss.xml", Name: "straitstimes", Category: "asia", Lang: "en", Variant: "full"},
	{URL: "https://www.bangkokpost.com/rss/data/topstories.xml", Name: "bangkokpost", Category: "asia", Lang: "en", Variant: "full"},
	{URL: "https://timesofindia.indiatimes.com/rssfeeds/296589292.cms", Name: "toi-world", Category: "asia", Lang: "en", Variant: "full"},
	{URL: "https://www.channelnewsasia.com/api/v1/rss-outbound-feed?_format=xml", Name: "cna-asia", Category: "asia", Lang: "en", Variant: "full"},
	{URL: "https://en.yna.co.kr/RSS/news.xml", Name: "yonhap-en", Category: "asia", Lang: "en", Variant: "full"},
	{URL: "https://focustaiwan.tw/rss", Name: "focustaiwan", Category: "asia", Lang: "en", Variant: "full"},

	// ── Americas (non-US) ───────────────────────────────────
	{URL: "https://rss.cbc.ca/lineup/world.xml", Name: "cbc-world", Category: "general", Lang: "en", Variant: "full"},
	{URL: "https://www.abc.net.au/news/feed/51120/rss.xml", Name: "abc-au", Category: "general", Lang: "en", Variant: "full"},

	// ── Economics / business ────────────────────────────────
	{URL: "https://feeds.bloomberg.com/economics/news.rss", Name: "bloomberg-econ", Category: "economics", Lang: "en", Variant: "full"},
	{URL: "https://www.economist.com/the-world-this-week/rss.xml", Name: "economist", Category: "economics", Lang: "en", Variant: "full"},
	{URL: "https://feeds.wsj.com/xml/rss/3_7085.xml", Name: "wsj-world", Category: "general", Lang: "en", Variant: "full"},

	// ── Defense / security ──────────────────────────────────
	{URL: "https://www.defenseone.com/rss/", Name: "defenseone", Category: "defense", Lang: "en", Variant: "full"},
	{URL: "https://www.janes.com/feeds/news", Name: "janes", Category: "defense", Lang: "en", Variant: "full"},

	// ── Think tanks ─────────────────────────────────────────
	{URL: "https://www.brookings.edu/feed/", Name: "brookings", Category: "thinktank", Lang: "en", Variant: "full"},
	{URL: "https://www.cfr.org/rss.xml", Name: "cfr", Category: "thinktank", Lang: "en", Variant: "full"},
	{URL: "https://www.chathamhouse.org/rss.xml", Name: "chathamhouse", Category: "thinktank", Lang: "en", Variant: "full"},
	{URL: "https://carnegieendowment.org/rss/solr/?q=*:*", Name: "carnegie", Category: "thinktank", Lang: "en", Variant: "full"},

	// ══════════════════════════════════════════════════════════
	//  tech (en) — technology-focused (~10 feeds)
	// ══════════════════════════════════════════════════════════
	{URL: "https://feeds.feedburner.com/TechCrunch/", Name: "techcrunch", Category: "tech", Lang: "en", Variant: "tech"},
	{URL: "https://feeds.arstechnica.com/arstechnica/index", Name: "arstechnica", Category: "tech", Lang: "en", Variant: "tech"},
	{URL: "https://www.wired.com/feed/rss", Name: "wired", Category: "tech", Lang: "en", Variant: "tech"},
	{URL: "https://www.theverge.com/rss/index.xml", Name: "theverge", Category: "tech", Lang: "en", Variant: "tech"},
	{URL: "https://www.technologyreview.com/feed/", Name: "mit-tech-review", Category: "tech", Lang: "en", Variant: "tech"},
	{URL: "https://www.zdnet.com/news/rss.xml", Name: "zdnet", Category: "tech", Lang: "en", Variant: "tech"},
	{URL: "https://www.engadget.com/rss.xml", Name: "engadget", Category: "tech", Lang: "en", Variant: "tech"},
	{URL: "https://www.tomshardware.com/feeds/all", Name: "tomshardware", Category: "tech", Lang: "en", Variant: "tech"},
	{URL: "https://restofworld.org/feed/", Name: "restofworld", Category: "tech", Lang: "en", Variant: "tech"},
	{URL: "https://spectrum.ieee.org/feeds/feed.rss", Name: "ieee-spectrum", Category: "tech", Lang: "en", Variant: "tech"},

	// ══════════════════════════════════════════════════════════
	//  geo (en) — disaster / seismic / humanitarian (~10 feeds)
	// ══════════════════════════════════════════════════════════
	{URL: "https://reliefweb.int/updates/rss.xml", Name: "reliefweb", Category: "humanitarian", Lang: "en", Variant: "geo"},
	{URL: "https://earthquake.usgs.gov/earthquakes/feed/v1.0/summary/significant_month.atom", Name: "usgs-quakes", Category: "seismic", Lang: "en", Variant: "geo"},
	{URL: "https://earthquake.usgs.gov/earthquakes/feed/v1.0/summary/4.5_week.atom", Name: "usgs-4.5", Category: "seismic", Lang: "en", Variant: "geo"},
	{URL: "https://www.gdacs.org/xml/rss.xml", Name: "gdacs", Category: "disaster", Lang: "en", Variant: "geo"},
	{URL: "https://volcanoes.usgs.gov/rss/vhpcap.xml", Name: "usgs-volcanoes", Category: "volcanic", Lang: "en", Variant: "geo"},
	{URL: "https://www.nhc.noaa.gov/index-at.xml", Name: "nhc-atlantic", Category: "hurricane", Lang: "en", Variant: "geo"},
	{URL: "https://www.nhc.noaa.gov/index-ep.xml", Name: "nhc-pacific", Category: "hurricane", Lang: "en", Variant: "geo"},
	{URL: "https://www.spc.noaa.gov/products/spcwwrss.xml", Name: "spc-severe", Category: "severe-weather", Lang: "en", Variant: "geo"},
	{URL: "https://tsunami.gov/events/xml/PHEBkml/AtomFeed.xml", Name: "tsunami-gov", Category: "tsunami", Lang: "en", Variant: "geo"},
	{URL: "https://www.who.int/feeds/entity/csr/don/en/rss.xml", Name: "who-outbreaks", Category: "health", Lang: "en", Variant: "geo"},

	// ══════════════════════════════════════════════════════════
	//  finance (en) — markets & economics (~10 feeds)
	// ══════════════════════════════════════════════════════════
	{URL: "https://feeds.bloomberg.com/markets/news.rss", Name: "bloomberg-markets", Category: "markets", Lang: "en", Variant: "finance"},
	{URL: "https://www.cnbc.com/id/100003114/device/rss/rss.html", Name: "cnbc", Category: "markets", Lang: "en", Variant: "finance"},
	{URL: "https://feeds.marketwatch.com/marketwatch/topstories/", Name: "marketwatch", Category: "markets", Lang: "en", Variant: "finance"},
	{URL: "https://feeds.marketwatch.com/marketwatch/marketpulse/", Name: "marketwatch-pulse", Category: "markets", Lang: "en", Variant: "finance"},
	{URL: "https://feeds.finance.yahoo.com/rss/2.0/headline?s=^GSPC&region=US&lang=en-US", Name: "yahoo-sp500", Category: "markets", Lang: "en", Variant: "finance"},
	{URL: "https://www.investing.com/rss/news.rss", Name: "investing-com", Category: "markets", Lang: "en", Variant: "finance"},
	{URL: "https://seekingalpha.com/feed.xml", Name: "seekingalpha", Category: "markets", Lang: "en", Variant: "finance"},
	{URL: "https://www.ft.com/markets?format=rss", Name: "ft-markets", Category: "markets", Lang: "en", Variant: "finance"},
	{URL: "https://feeds.reuters.com/news/wealth", Name: "reuters-wealth", Category: "finance", Lang: "en", Variant: "finance"},
	{URL: "https://www.imf.org/en/News/rss", Name: "imf-news", Category: "economics", Lang: "en", Variant: "finance"},

	// ══════════════════════════════════════════════════════════
	//  quick (en) — fast-refresh subset (~10 feeds)
	// ══════════════════════════════════════════════════════════
	{URL: "https://feeds.reuters.com/reuters/topNews", Name: "reuters", Category: "general", Lang: "en", Variant: "quick"},
	{URL: "https://feeds.reuters.com/reuters/worldNews", Name: "reuters-world", Category: "general", Lang: "en", Variant: "quick"},
	{URL: "https://feeds.bbci.co.uk/news/world/rss.xml", Name: "bbc-world", Category: "general", Lang: "en", Variant: "quick"},
	{URL: "https://rss.nytimes.com/services/xml/rss/nyt/World.xml", Name: "nyt-world", Category: "general", Lang: "en", Variant: "quick"},
	{URL: "http://rss.cnn.com/rss/edition_world.rss", Name: "cnn-world", Category: "general", Lang: "en", Variant: "quick"},
	{URL: "https://www.aljazeera.com/xml/rss/all.xml", Name: "aljazeera", Category: "general", Lang: "en", Variant: "quick"},
	{URL: "https://www.theguardian.com/world/rss", Name: "guardian-world", Category: "general", Lang: "en", Variant: "quick"},
	{URL: "https://feeds.washingtonpost.com/rss/world", Name: "washpost-world", Category: "general", Lang: "en", Variant: "quick"},
	{URL: "https://feeds.bloomberg.com/economics/news.rss", Name: "bloomberg-econ", Category: "economics", Lang: "en", Variant: "quick"},
	{URL: "https://asia.nikkei.com/rss", Name: "nikkei-asia", Category: "asia", Lang: "en", Variant: "quick"},
}

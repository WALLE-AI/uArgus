package news

import "testing"

func TestGetFeeds_FullVariant(t *testing.T) {
	feeds := GetFeeds("full", "en")
	if len(feeds) < 50 {
		t.Fatalf("expected ≥50 full feeds, got %d", len(feeds))
	}
	// verify at least one Reuters entry
	found := false
	for _, f := range feeds {
		if f.Name == "reuters" {
			found = true
			break
		}
	}
	if !found {
		t.Fatal("reuters not found in full feeds")
	}
}

func TestGetFeeds_AllVariants(t *testing.T) {
	for _, v := range AllVariants() {
		feeds := GetFeeds(v, "en")
		if len(feeds) == 0 {
			t.Fatalf("variant %q returned 0 feeds", v)
		}
	}
}

func TestGetFeeds_FallbackToFull(t *testing.T) {
	feeds := GetFeeds("nonexistent", "en")
	if len(feeds) == 0 {
		t.Fatal("expected fallback to full variant")
	}
}

func TestParseFeedResponse_RSS(t *testing.T) {
	rss := []byte(`<?xml version="1.0"?>
<rss version="2.0"><channel><title>T</title>
<item><title>Test Article</title><link>https://example.com/1</link>
<pubDate>Mon, 02 Jan 2006 15:04:05 GMT</pubDate></item>
</channel></rss>`)

	feed := FeedEntry{Name: "test-src", Variant: "full", Lang: "en"}
	items := ParseFeedResponse(rss, feed)
	if len(items) != 1 {
		t.Fatalf("expected 1 item, got %d", len(items))
	}
	if items[0].Title != "Test Article" {
		t.Fatalf("unexpected title: %s", items[0].Title)
	}
	if items[0].Source != "test-src" {
		t.Fatalf("expected source test-src, got %s", items[0].Source)
	}
	if items[0].Hash == "" {
		t.Fatal("expected non-empty hash")
	}
}

func TestParseFeedResponse_Atom(t *testing.T) {
	atom := []byte(`<?xml version="1.0"?>
<feed xmlns="http://www.w3.org/2005/Atom">
<entry><title>Atom Entry</title>
<link href="https://example.com/a" rel="alternate"/>
<published>2024-01-01T00:00:00Z</published>
<id>1</id></entry>
</feed>`)

	feed := FeedEntry{Name: "atom-src", Variant: "full", Lang: "en"}
	items := ParseFeedResponse(atom, feed)
	if len(items) != 1 {
		t.Fatalf("expected 1 item, got %d", len(items))
	}
	if items[0].Title != "Atom Entry" {
		t.Fatalf("unexpected title: %s", items[0].Title)
	}
}

func TestParseFeedResponse_Invalid(t *testing.T) {
	items := ParseFeedResponse([]byte("not xml"), FeedEntry{Name: "bad"})
	if items != nil {
		t.Fatalf("expected nil items for invalid input, got %d", len(items))
	}
}

func TestParsedItem_TrackHash(t *testing.T) {
	p := ParsedItem{Hash: "abc123"}
	if p.TrackHash() != "abc123" {
		t.Fatalf("expected abc123, got %s", p.TrackHash())
	}
}

func TestTruncate(t *testing.T) {
	long := "abcdefghijklmnopqrstuvwxyz"
	result := truncate(long, 10)
	if len(result) > 14 { // 10 + "…" (3 bytes in UTF-8)
		t.Fatalf("expected truncated, got len=%d", len(result))
	}
	short := "abc"
	if truncate(short, 10) != "abc" {
		t.Fatal("short string should not be truncated")
	}
}

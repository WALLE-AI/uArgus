package parser

import "testing"

const sampleRSS = `<?xml version="1.0" encoding="UTF-8"?>
<rss version="2.0">
<channel>
  <title>Test Feed</title>
  <item>
    <title>Breaking: Major Event</title>
    <link>https://example.com/1</link>
    <description>Something happened.</description>
    <pubDate>Mon, 02 Jan 2006 15:04:05 GMT</pubDate>
  </item>
  <item>
    <title>Second Story</title>
    <link>https://example.com/2</link>
    <description><![CDATA[<p>HTML content</p>]]></description>
    <pubDate>Tue, 03 Jan 2006 10:00:00 GMT</pubDate>
  </item>
</channel>
</rss>`

const sampleAtom = `<?xml version="1.0" encoding="UTF-8"?>
<feed xmlns="http://www.w3.org/2005/Atom">
  <title>Atom Feed</title>
  <entry>
    <title>Paper Title</title>
    <link href="https://arxiv.org/abs/1234" rel="alternate"/>
    <summary>Abstract text</summary>
    <published>2024-01-01T00:00:00Z</published>
    <id>urn:arxiv:1234</id>
  </entry>
  <entry>
    <title>Second Paper</title>
    <link href="https://arxiv.org/abs/5678" rel="alternate"/>
    <summary>Another abstract</summary>
    <updated>2024-02-01T00:00:00Z</updated>
    <id>urn:arxiv:5678</id>
  </entry>
</feed>`

const sampleICS = `BEGIN:VCALENDAR
BEGIN:VEVENT
SUMMARY:Tech Conference 2024
DESCRIPTION:Annual tech event
LOCATION:San Francisco
DTSTART;VALUE=DATE:20240601
DTEND;VALUE=DATE:20240603
URL:https://example.com/conf
UID:conf-2024@example.com
END:VEVENT
BEGIN:VEVENT
SUMMARY:AI Summit
DESCRIPTION:AI focused event
LOCATION:London
DTSTART:20240715T090000Z
DTEND:20240716T180000Z
UID:ai-summit@example.com
END:VEVENT
END:VCALENDAR`

func TestParseRssXML(t *testing.T) {
	items, err := ParseRssXML([]byte(sampleRSS))
	if err != nil {
		t.Fatal(err)
	}
	if len(items) != 2 {
		t.Fatalf("expected 2 items, got %d", len(items))
	}
	if items[0].Title != "Breaking: Major Event" {
		t.Fatalf("unexpected title: %s", items[0].Title)
	}
	if items[0].Link != "https://example.com/1" {
		t.Fatalf("unexpected link: %s", items[0].Link)
	}
	// CDATA is unwrapped, then HTML tags stripped
	if items[1].Description == "" {
		t.Fatal("expected non-empty description after CDATA unwrap")
	}
}

func TestParseRssXML_Empty(t *testing.T) {
	_, err := ParseRssXML([]byte("not xml"))
	if err == nil {
		t.Fatal("expected error on invalid XML")
	}
}

func TestParseAtomXML(t *testing.T) {
	entries, err := ParseAtomXML([]byte(sampleAtom))
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(entries))
	}
	if entries[0].Link != "https://arxiv.org/abs/1234" {
		t.Fatalf("unexpected link: %s", entries[0].Link)
	}
	// second entry uses Updated as Published fallback
	if entries[1].Published == "" {
		t.Fatal("expected Published to be populated from Updated")
	}
}

func TestParseICS(t *testing.T) {
	events := ParseICS([]byte(sampleICS))
	if len(events) != 2 {
		t.Fatalf("expected 2 events, got %d", len(events))
	}
	if events[0].Summary != "Tech Conference 2024" {
		t.Fatalf("unexpected summary: %s", events[0].Summary)
	}
	if events[0].Location != "San Francisco" {
		t.Fatalf("unexpected location: %s", events[0].Location)
	}
	if events[1].DTStart != "20240715T090000Z" {
		t.Fatalf("unexpected dtstart: %s", events[1].DTStart)
	}
}

func TestCleanText(t *testing.T) {
	cases := []struct{ in, want string }{
		{"<b>bold</b>", "bold"},
		{"&amp; &lt; &gt;", "& < >"},
		{"  spaces  ", "spaces"},
	}
	for _, c := range cases {
		got := cleanText(c.in)
		if got != c.want {
			t.Errorf("cleanText(%q) = %q, want %q", c.in, got, c.want)
		}
	}
}

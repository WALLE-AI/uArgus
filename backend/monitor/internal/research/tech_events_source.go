package research

import (
	"context"
	"log/slog"
	"strings"
	"time"

	"github.com/WALL-AI/uArgus/backend/monitor/internal/fetcher"
	"github.com/WALL-AI/uArgus/backend/monitor/internal/fetcher/parser"
	"github.com/WALL-AI/uArgus/backend/monitor/internal/registry"
)

const (
	techmemeICSURL  = "https://www.techmeme.com/events.ics"
	devEventsRSSURL = "https://dev.events/rss.xml"
)

// TechEvent is the unified representation of a tech event from any source.
type TechEvent struct {
	Name        string  `json:"name"`
	Description string  `json:"description,omitempty"`
	Location    string  `json:"location,omitempty"`
	StartDate   string  `json:"startDate"`
	EndDate     string  `json:"endDate,omitempty"`
	URL         string  `json:"url,omitempty"`
	Source      string  `json:"source"` // "techmeme" | "dev.events" | "curated"
	Lat         float64 `json:"lat,omitempty"`
	Lng         float64 `json:"lng,omitempty"`
}

// TechEventsSource fetches tech events from Techmeme ICS, dev.events RSS, and a curated list.
type TechEventsSource struct {
	spec       registry.SourceSpec
	httpClient *fetcher.Client
}

// NewTechEventsSource creates a TechEventsSource.
func NewTechEventsSource(spec registry.SourceSpec, httpClient *fetcher.Client) *TechEventsSource {
	return &TechEventsSource{spec: spec, httpClient: httpClient}
}

func (s *TechEventsSource) Name() string              { return s.spec.CanonicalKey }
func (s *TechEventsSource) Spec() registry.SourceSpec { return s.spec }
func (s *TechEventsSource) Dependencies() []string    { return nil }

func (s *TechEventsSource) Run(ctx context.Context) (*registry.FetchResult, error) {
	start := time.Now()
	logger := slog.With("source", s.Name())

	var allEvents []TechEvent

	// 1. Techmeme ICS
	if body, _, err := s.httpClient.Get(ctx, techmemeICSURL); err == nil {
		icsEvents := parser.ParseICS(body)
		for _, e := range icsEvents {
			te := TechEvent{
				Name:        e.Summary,
				Description: e.Description,
				Location:    e.Location,
				StartDate:   e.DTStart,
				EndDate:     e.DTEnd,
				URL:         e.URL,
				Source:      "techmeme",
			}
			enrichCoords(&te)
			allEvents = append(allEvents, te)
		}
		logger.Info("techmeme fetched", "events", len(icsEvents))
	} else {
		logger.Warn("techmeme fetch failed", "err", err)
	}

	// 2. dev.events RSS
	if body, _, err := s.httpClient.Get(ctx, devEventsRSSURL); err == nil {
		rssItems, parseErr := parser.ParseRssXML(body)
		if parseErr == nil {
			for _, item := range rssItems {
				te := TechEvent{
					Name:      item.Title,
					URL:       item.Link,
					StartDate: item.PubDate,
					Source:    "dev.events",
				}
				allEvents = append(allEvents, te)
			}
			logger.Info("dev.events fetched", "events", len(rssItems))
		}
	} else {
		logger.Warn("dev.events fetch failed", "err", err)
	}

	// 3. Curated major conferences
	allEvents = append(allEvents, curatedTechEvents()...)

	// 4. Filter to future events only + dedup
	allEvents = filterFutureEvents(allEvents)
	allEvents = dedupEvents(allEvents)

	return &registry.FetchResult{
		Data: allEvents,
		Metrics: registry.FetchMetrics{
			Duration:    time.Since(start),
			RecordCount: len(allEvents),
		},
	}, nil
}

// enrichCoords attempts to match the event location to known city coordinates.
func enrichCoords(te *TechEvent) {
	if te.Location == "" {
		return
	}
	loc := strings.ToLower(te.Location)
	for _, c := range CityCoords {
		if strings.Contains(loc, strings.ToLower(c.Name)) {
			te.Lat = c.Lat
			te.Lng = c.Lng
			return
		}
	}
}

// filterFutureEvents keeps only events whose start date is in the future (or unparseable).
func filterFutureEvents(events []TechEvent) []TechEvent {
	now := time.Now()
	var result []TechEvent
	for _, e := range events {
		t, err := parseFlexDate(e.StartDate)
		if err != nil || t.After(now.AddDate(0, 0, -1)) {
			result = append(result, e)
		}
	}
	return result
}

// dedupEvents removes duplicate events by normalised name.
func dedupEvents(events []TechEvent) []TechEvent {
	seen := make(map[string]bool)
	var result []TechEvent
	for _, e := range events {
		key := strings.ToLower(strings.TrimSpace(e.Name))
		if !seen[key] {
			seen[key] = true
			result = append(result, e)
		}
	}
	return result
}

// parseFlexDate tries several date formats.
func parseFlexDate(s string) (time.Time, error) {
	s = strings.TrimSpace(s)
	formats := []string{
		"20060102",
		"20060102T150405Z",
		"2006-01-02",
		"2006-01-02T15:04:05Z",
		time.RFC3339,
		time.RFC1123,
		time.RFC1123Z,
		"Mon, 02 Jan 2006 15:04:05 MST",
		"Mon, 02 Jan 2006 15:04:05 -0700",
	}
	for _, f := range formats {
		if t, err := time.Parse(f, s); err == nil {
			return t, nil
		}
	}
	return time.Time{}, nil
}

// curatedTechEvents returns a static list of major annual tech conferences.
func curatedTechEvents() []TechEvent {
	year := time.Now().Year()
	return []TechEvent{
		{Name: "CES", Location: "Las Vegas", StartDate: formatDate(year, 1, 7), EndDate: formatDate(year, 1, 10), URL: "https://www.ces.tech", Source: "curated", Lat: 36.1699, Lng: -115.1398},
		{Name: "MWC Barcelona", Location: "Barcelona", StartDate: formatDate(year, 2, 26), EndDate: formatDate(year, 3, 1), URL: "https://www.mwcbarcelona.com", Source: "curated", Lat: 41.3874, Lng: 2.1686},
		{Name: "GDC", Location: "San Francisco", StartDate: formatDate(year, 3, 17), EndDate: formatDate(year, 3, 21), URL: "https://gdconf.com", Source: "curated", Lat: 37.7749, Lng: -122.4194},
		{Name: "Google I/O", Location: "Mountain View", StartDate: formatDate(year, 5, 14), EndDate: formatDate(year, 5, 15), URL: "https://io.google", Source: "curated", Lat: 37.3861, Lng: -122.0839},
		{Name: "WWDC", Location: "Cupertino", StartDate: formatDate(year, 6, 9), EndDate: formatDate(year, 6, 13), URL: "https://developer.apple.com/wwdc", Source: "curated", Lat: 37.3230, Lng: -122.0322},
		{Name: "AWS re:Invent", Location: "Las Vegas", StartDate: formatDate(year, 12, 1), EndDate: formatDate(year, 12, 5), URL: "https://reinvent.awsevents.com", Source: "curated", Lat: 36.1699, Lng: -115.1398},
		{Name: "Microsoft Build", Location: "Seattle", StartDate: formatDate(year, 5, 19), EndDate: formatDate(year, 5, 22), URL: "https://build.microsoft.com", Source: "curated", Lat: 47.6062, Lng: -122.3321},
		{Name: "KubeCon NA", Location: "Salt Lake City", StartDate: formatDate(year, 11, 10), EndDate: formatDate(year, 11, 13), URL: "https://events.linuxfoundation.org/kubecon-cloudnativecon-north-america/", Source: "curated", Lat: 40.7608, Lng: -111.8910},
		{Name: "DEF CON", Location: "Las Vegas", StartDate: formatDate(year, 8, 7), EndDate: formatDate(year, 8, 10), URL: "https://defcon.org", Source: "curated", Lat: 36.1699, Lng: -115.1398},
		{Name: "Black Hat USA", Location: "Las Vegas", StartDate: formatDate(year, 8, 2), EndDate: formatDate(year, 8, 7), URL: "https://www.blackhat.com", Source: "curated", Lat: 36.1699, Lng: -115.1398},
	}
}

func formatDate(year, month, day int) string {
	return time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC).Format("2006-01-02")
}

package parser

import (
	"strings"
)

// ICSEvent represents one event parsed from an iCalendar file.
type ICSEvent struct {
	Summary     string `json:"summary"`
	Description string `json:"description"`
	Location    string `json:"location"`
	DTStart     string `json:"dtStart"`
	DTEnd       string `json:"dtEnd"`
	URL         string `json:"url"`
	UID         string `json:"uid"`
}

// ParseICS parses an iCalendar (.ics) body into events.
func ParseICS(body []byte) []ICSEvent {
	lines := strings.Split(strings.ReplaceAll(string(body), "\r\n", "\n"), "\n")
	var events []ICSEvent
	var current *ICSEvent

	for _, line := range lines {
		line = strings.TrimSpace(line)
		switch {
		case line == "BEGIN:VEVENT":
			current = &ICSEvent{}
		case line == "END:VEVENT" && current != nil:
			events = append(events, *current)
			current = nil
		case current != nil:
			key, val := splitICSLine(line)
			switch key {
			case "SUMMARY":
				current.Summary = val
			case "DESCRIPTION":
				current.Description = val
			case "LOCATION":
				current.Location = val
			case "DTSTART":
				current.DTStart = val
			case "DTEND":
				current.DTEnd = val
			case "URL":
				current.URL = val
			case "UID":
				current.UID = val
			}
		}
	}
	return events
}

func splitICSLine(line string) (string, string) {
	// handle DTSTART;VALUE=DATE:20240101 format
	idx := strings.Index(line, ":")
	if idx < 0 {
		return line, ""
	}
	key := line[:idx]
	val := line[idx+1:]
	// strip parameters from key (e.g. DTSTART;VALUE=DATE → DTSTART)
	if semi := strings.Index(key, ";"); semi >= 0 {
		key = key[:semi]
	}
	return strings.ToUpper(key), val
}

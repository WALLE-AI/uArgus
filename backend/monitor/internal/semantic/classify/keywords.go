package classify

// defaultKeywordLevels returns 7 layers of geo+tech keywords.
func defaultKeywordLevels() []KeywordLevel {
	return []KeywordLevel{
		{
			Name:     "conflict",
			Severity: 7,
			Keywords: []string{
				"war", "invasion", "missile", "airstrike", "bombing",
				"nuclear", "military", "troops", "ceasefire", "casualties",
				"conflict", "drone strike", "artillery", "combat",
			},
		},
		{
			Name:     "terrorism",
			Severity: 7,
			Keywords: []string{
				"terrorist", "terrorism", "attack", "hostage", "extremist",
				"insurgent", "suicide bomb", "militant", "jihad",
			},
		},
		{
			Name:     "disaster",
			Severity: 6,
			Keywords: []string{
				"earthquake", "tsunami", "hurricane", "typhoon", "flood",
				"wildfire", "tornado", "volcanic", "landslide", "cyclone",
				"drought", "famine",
			},
		},
		{
			Name:     "political",
			Severity: 5,
			Keywords: []string{
				"election", "coup", "protest", "sanctions", "regime",
				"parliament", "impeachment", "referendum", "summit",
				"diplomatic", "treaty", "embargo",
			},
		},
		{
			Name:     "economic",
			Severity: 4,
			Keywords: []string{
				"recession", "inflation", "gdp", "trade war", "tariff",
				"default", "bankruptcy", "market crash", "supply chain",
				"interest rate", "central bank", "federal reserve",
			},
		},
		{
			Name:     "health",
			Severity: 4,
			Keywords: []string{
				"pandemic", "epidemic", "outbreak", "vaccine", "who",
				"virus", "quarantine", "covid", "h5n1", "bird flu",
			},
		},
		{
			Name:     "technology",
			Severity: 3,
			Keywords: []string{
				"cyber attack", "data breach", "ransomware", "hacking",
				"AI", "artificial intelligence", "quantum", "semiconductor",
				"chip shortage", "export control",
			},
		},
	}
}

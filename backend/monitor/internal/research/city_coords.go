package research

import "math"

// LatLng is a geographic coordinate.
type LatLng struct {
	Lat float64
	Lng float64
}

// City holds a city name and its coordinates.
type City struct {
	Name string
	LatLng
}

// CityCoords is a curated list of major cities for geospatial lookups.
var CityCoords = []City{
	{"Tokyo", LatLng{35.6762, 139.6503}},
	{"New York", LatLng{40.7128, -74.0060}},
	{"London", LatLng{51.5074, -0.1278}},
	{"Beijing", LatLng{39.9042, 116.4074}},
	{"Moscow", LatLng{55.7558, 37.6173}},
	{"Washington DC", LatLng{38.9072, -77.0369}},
	{"Berlin", LatLng{52.5200, 13.4050}},
	{"Paris", LatLng{48.8566, 2.3522}},
	{"Seoul", LatLng{37.5665, 126.9780}},
	{"New Delhi", LatLng{28.6139, 77.2090}},
	{"Singapore", LatLng{1.3521, 103.8198}},
	{"Sydney", LatLng{-33.8688, 151.2093}},
	{"Dubai", LatLng{25.2048, 55.2708}},
	{"São Paulo", LatLng{-23.5505, -46.6333}},
	{"Istanbul", LatLng{41.0082, 28.9784}},
	{"Taipei", LatLng{25.0330, 121.5654}},
	{"Bangkok", LatLng{13.7563, 100.5018}},
	{"Cairo", LatLng{30.0444, 31.2357}},
	{"Lagos", LatLng{6.5244, 3.3792}},
	{"Nairobi", LatLng{-1.2921, 36.8219}},
}

// Nearest returns the closest city to the given coordinates.
func Nearest(lat, lng float64) City {
	var best City
	bestDist := math.MaxFloat64

	for _, c := range CityCoords {
		d := haversine(lat, lng, c.Lat, c.Lng)
		if d < bestDist {
			bestDist = d
			best = c
		}
	}
	return best
}

func haversine(lat1, lon1, lat2, lon2 float64) float64 {
	const R = 6371.0
	dLat := (lat2 - lat1) * math.Pi / 180
	dLon := (lon2 - lon1) * math.Pi / 180
	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(lat1*math.Pi/180)*math.Cos(lat2*math.Pi/180)*
			math.Sin(dLon/2)*math.Sin(dLon/2)
	return R * 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
}

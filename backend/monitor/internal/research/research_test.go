package research

import (
	"math"
	"testing"
)

func TestCityCoords_HasEntries(t *testing.T) {
	if len(CityCoords) < 20 {
		t.Fatalf("expected ≥20 cities, got %d", len(CityCoords))
	}
}

func TestNearest_Tokyo(t *testing.T) {
	city := Nearest(35.68, 139.69)
	if city.Name != "Tokyo" {
		t.Fatalf("expected Tokyo, got %s", city.Name)
	}
}

func TestNearest_London(t *testing.T) {
	city := Nearest(51.50, -0.13)
	if city.Name != "London" {
		t.Fatalf("expected London, got %s", city.Name)
	}
}

func TestNearest_Sydney(t *testing.T) {
	city := Nearest(-33.87, 151.21)
	if city.Name != "Sydney" {
		t.Fatalf("expected Sydney, got %s", city.Name)
	}
}

func TestHaversine_KnownDistance(t *testing.T) {
	// London to Paris ≈ 343 km
	d := haversine(51.5074, -0.1278, 48.8566, 2.3522)
	if math.Abs(d-343) > 10 {
		t.Fatalf("London-Paris distance expected ~343km, got %.1f", d)
	}
}

func TestHaversine_SamePoint(t *testing.T) {
	d := haversine(0, 0, 0, 0)
	if d != 0 {
		t.Fatalf("expected 0 distance for same point, got %f", d)
	}
}

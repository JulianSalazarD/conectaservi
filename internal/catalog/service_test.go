package catalog

import (
	"errors"
	"testing"
)

func TestNewLocation(t *testing.T) {
	tests := []struct {
		name    string
		lat     float64
		lng     float64
		wantErr error
	}{
		{"bogotá", 4.65, -74.08, nil},
		{"origen", 0, 0, nil},
		{"límite sur", -90, 0, nil},
		{"límite norte", 90, 0, nil},
		{"límite oeste", 0, -180, nil},
		{"límite este", 0, 180, nil},

		{"lat > 90", 91, 0, ErrInvalidLocation},
		{"lat < -90", -91, 0, ErrInvalidLocation},
		{"lng > 180", 0, 181, ErrInvalidLocation},
		{"lng < -180", 0, -181, ErrInvalidLocation},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			loc, err := NewLocation(tc.lat, tc.lng)
			if !errors.Is(err, tc.wantErr) {
				t.Fatalf("err = %v, want %v", err, tc.wantErr)
			}
			if tc.wantErr == nil {
				if loc.Lat != tc.lat || loc.Lng != tc.lng {
					t.Errorf("got (%v,%v), want (%v,%v)", loc.Lat, loc.Lng, tc.lat, tc.lng)
				}
			}
		})
	}
}

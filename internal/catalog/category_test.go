package catalog

import (
	"errors"
	"testing"
)

func TestNewCategory(t *testing.T) {
	tests := []struct {
		name    string
		nombre  string
		slug    string
		wantErr error
	}{
		{"válido", "Plomería", "plomeria", nil},
		{"slug con guión", "Electricidad Industrial", "electricidad-industrial", nil},
		{"slug alfanumérico", "Rango 24/7", "rango-24-7", nil},

		{"nombre vacío", "", "plomeria", ErrInvalidName},
		{"nombre solo espacios", "   ", "plomeria", ErrInvalidName},

		{"slug vacío", "Plomería", "", ErrInvalidSlug},
		{"slug con espacios", "Plomería", "plomeria industrial", ErrInvalidSlug},
		{"slug con mayúsculas", "Plomería", "Plomeria", ErrInvalidSlug},
		{"slug con guión inicial", "Plomería", "-plomeria", ErrInvalidSlug},
		{"slug con guión final", "Plomería", "plomeria-", ErrInvalidSlug},
		{"slug con doble guión", "Plomería", "plomeria--industrial", ErrInvalidSlug},
		{"slug con carácter especial", "Plomería", "plomería", ErrInvalidSlug},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			c, err := NewCategory(tc.nombre, tc.slug, nil)
			if !errors.Is(err, tc.wantErr) {
				t.Fatalf("err = %v, want %v", err, tc.wantErr)
			}
			if tc.wantErr == nil {
				if c == nil {
					t.Fatal("expected non-nil category")
				}
				if c.ID.String() == "" {
					t.Error("expected auto-generated ID")
				}
			}
		})
	}
}

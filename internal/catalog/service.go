package catalog

import (
	"strings"
	"time"

	"github.com/google/uuid"
)

type Location struct {
	Lat float64
	Lng float64
}

func NewLocation(lat, lng float64) (Location, error) {
	if lat < -90 || lat > 90 {
		return Location{}, ErrInvalidLocation
	}
	if lng < -180 || lng > 180 {
		return Location{}, ErrInvalidLocation
	}
	return Location{Lat: lat, Lng: lng}, nil
}

type Service struct {
	ID          uuid.UUID
	ProviderID  uuid.UUID
	CategoryID  uuid.UUID
	Titulo      string
	Descripcion string
	PrecioBase  float64
	Lat         *float64
	Lng         *float64
	RadioKm     *float64
	IsActive    bool
	CreatedAt   time.Time
}

func NewService(
	providerID, categoryID uuid.UUID,
	titulo, descripcion string,
	precioBase float64,
	lat, lng, radioKm *float64,
) (*Service, error) {
	titulo = strings.TrimSpace(titulo)
	if titulo == "" {
		return nil, ErrInvalidTitle
	}
	if precioBase < 0 {
		return nil, ErrInvalidPrice
	}
	if (lat == nil) != (lng == nil) {
		return nil, ErrInvalidLocation
	}
	if lat != nil {
		if _, err := NewLocation(*lat, *lng); err != nil {
			return nil, err
		}
	}
	return &Service{
		ID:          uuid.New(),
		ProviderID:  providerID,
		CategoryID:  categoryID,
		Titulo:      titulo,
		Descripcion: strings.TrimSpace(descripcion),
		PrecioBase:  precioBase,
		Lat:         lat,
		Lng:         lng,
		RadioKm:     radioKm,
		IsActive:    true,
		CreatedAt:   time.Now().UTC(),
	}, nil
}

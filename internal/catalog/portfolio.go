package catalog

import (
	"strings"

	"github.com/google/uuid"
)

type PortfolioItem struct {
	ID         uuid.UUID
	ServiceID  uuid.UUID
	StorageURL string
	Titulo     string
	Orden      int
}

func NewPortfolioItem(serviceID uuid.UUID, storageURL, titulo string, orden int) (*PortfolioItem, error) {
	storageURL = strings.TrimSpace(storageURL)
	if storageURL == "" {
		return nil, ErrInvalidStorageURL
	}
	return &PortfolioItem{
		ID:         uuid.New(),
		ServiceID:  serviceID,
		StorageURL: storageURL,
		Titulo:     strings.TrimSpace(titulo),
		Orden:      orden,
	}, nil
}

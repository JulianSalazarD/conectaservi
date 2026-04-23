package catalog

import (
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"
)

var slugRegex = regexp.MustCompile(`^[a-z0-9]+(?:-[a-z0-9]+)*$`)

type Category struct {
	ID        uuid.UUID
	Nombre    string
	Slug      string
	ParentID  *uuid.UUID
	CreatedAt time.Time
}

func NewCategory(nombre, slug string, parentID *uuid.UUID) (*Category, error) {
	nombre = strings.TrimSpace(nombre)
	if nombre == "" {
		return nil, ErrInvalidName
	}
	slug = strings.TrimSpace(slug)
	if !slugRegex.MatchString(slug) {
		return nil, ErrInvalidSlug
	}
	return &Category{
		ID:        uuid.New(),
		Nombre:    nombre,
		Slug:      slug,
		ParentID:  parentID,
		CreatedAt: time.Now().UTC(),
	}, nil
}

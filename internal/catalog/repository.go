package catalog

import (
	"context"

	"github.com/google/uuid"
)

type CategoryRepository interface {
	Insert(ctx context.Context, c *Category) error
	FindAll(ctx context.Context) ([]*Category, error)
	FindByID(ctx context.Context, id uuid.UUID) (*Category, error)
	Update(ctx context.Context, c *Category) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type ServiceFilter struct {
	CategoryID *uuid.UUID
	IsActive   *bool
}

type ServiceRepository interface {
	Insert(ctx context.Context, s *Service) error
	FindAll(ctx context.Context, filter ServiceFilter) ([]*Service, error)
	FindByID(ctx context.Context, id uuid.UUID) (*Service, error)
	Update(ctx context.Context, s *Service) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type PortfolioRepository interface {
	Insert(ctx context.Context, p *PortfolioItem) error
	FindByServiceID(ctx context.Context, serviceID uuid.UUID) ([]*PortfolioItem, error)
	FindByID(ctx context.Context, id uuid.UUID) (*PortfolioItem, error)
	Update(ctx context.Context, p *PortfolioItem) error
	Delete(ctx context.Context, id uuid.UUID) error
}

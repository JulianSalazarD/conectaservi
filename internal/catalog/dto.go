package catalog

import "github.com/google/uuid"

type CreateCategoryRequest struct {
	Nombre   string     `json:"nombre" binding:"required"`
	Slug     string     `json:"slug" binding:"required"`
	ParentID *uuid.UUID `json:"parent_id"`
}

type UpdateCategoryRequest struct {
	Nombre   string     `json:"nombre" binding:"required"`
	Slug     string     `json:"slug" binding:"required"`
	ParentID *uuid.UUID `json:"parent_id"`
}

type CreateServiceRequest struct {
	ProviderID  uuid.UUID `json:"provider_id" binding:"required"`
	CategoryID  uuid.UUID `json:"category_id" binding:"required"`
	Titulo      string    `json:"titulo" binding:"required"`
	Descripcion string    `json:"descripcion"`
	PrecioBase  float64   `json:"precio_base" binding:"gte=0"`
	Lat         *float64  `json:"lat"`
	Lng         *float64  `json:"lng"`
	RadioKm     *float64  `json:"radio_km"`
}

type UpdateServiceRequest struct {
	ProviderID  uuid.UUID `json:"provider_id" binding:"required"`
	CategoryID  uuid.UUID `json:"category_id" binding:"required"`
	Titulo      string    `json:"titulo" binding:"required"`
	Descripcion string    `json:"descripcion"`
	PrecioBase  float64   `json:"precio_base" binding:"gte=0"`
	Lat         *float64  `json:"lat"`
	Lng         *float64  `json:"lng"`
	RadioKm     *float64  `json:"radio_km"`
	IsActive    bool      `json:"is_active"`
}

type CreatePortfolioItemRequest struct {
	StorageURL string `json:"storage_url" binding:"required"`
	Titulo     string `json:"titulo"`
	Orden      int    `json:"orden"`
}

type UpdatePortfolioItemRequest struct {
	StorageURL string `json:"storage_url" binding:"required"`
	Titulo     string `json:"titulo"`
	Orden      int    `json:"orden"`
}

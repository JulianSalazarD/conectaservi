package catalog

import "errors"

var (
	ErrCategoryNotFound      = errors.New("category not found")
	ErrServiceNotFound       = errors.New("service not found")
	ErrPortfolioItemNotFound = errors.New("portfolio item not found")
	ErrProviderNotFound      = errors.New("provider not found")

	ErrDuplicateSlug       = errors.New("category slug already exists")
	ErrCategoryHasServices = errors.New("category has associated services")

	ErrInvalidName       = errors.New("invalid name: must not be empty")
	ErrInvalidSlug       = errors.New("invalid slug: must be lowercase kebab-case")
	ErrInvalidTitle      = errors.New("invalid title: must not be empty")
	ErrInvalidPrice      = errors.New("invalid price: must be >= 0")
	ErrInvalidLocation   = errors.New("invalid location: lat in [-90,90], lng in [-180,180]")
	ErrInvalidStorageURL = errors.New("invalid storage url: must not be empty")
)

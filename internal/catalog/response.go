package catalog

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

func writeError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, ErrCategoryNotFound),
		errors.Is(err, ErrServiceNotFound),
		errors.Is(err, ErrPortfolioItemNotFound),
		errors.Is(err, ErrProviderNotFound):
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	case errors.Is(err, ErrDuplicateSlug),
		errors.Is(err, ErrCategoryHasServices):
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
	case errors.Is(err, ErrInvalidName),
		errors.Is(err, ErrInvalidSlug),
		errors.Is(err, ErrInvalidTitle),
		errors.Is(err, ErrInvalidPrice),
		errors.Is(err, ErrInvalidLocation),
		errors.Is(err, ErrInvalidStorageURL):
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
}

func writeBindingError(c *gin.Context, err error) {
	c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
}

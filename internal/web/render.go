package web

import (
	"bytes"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/JulianSalazarD/conectaservi/internal/catalog"
)

func (m *Module) render(c *gin.Context, name string, data gin.H) {
	if data == nil {
		data = gin.H{}
	}
	if _, ok := data["Title"]; !ok {
		data["Title"] = "ConectaServi"
	}
	var buf bytes.Buffer
	if err := m.tpl.ExecuteTemplate(&buf, name, data); err != nil {
		c.String(http.StatusInternalServerError, "template error: %v", err)
		return
	}
	c.Data(http.StatusOK, "text/html; charset=utf-8", buf.Bytes())
}

// userMessage maps domain errors to friendly Spanish messages for the UI.
func userMessage(err error) string {
	switch {
	case errors.Is(err, catalog.ErrInvalidName):
		return "El nombre no puede estar vacío."
	case errors.Is(err, catalog.ErrInvalidSlug):
		return "Slug inválido: usa solo minúsculas, números y guiones (kebab-case)."
	case errors.Is(err, catalog.ErrDuplicateSlug):
		return "Ya existe una categoría con ese slug."
	case errors.Is(err, catalog.ErrCategoryNotFound):
		return "La categoría no existe."
	case errors.Is(err, catalog.ErrCategoryHasServices):
		return "No se puede eliminar: la categoría tiene servicios asociados."
	case errors.Is(err, catalog.ErrInvalidTitle):
		return "El título no puede estar vacío."
	case errors.Is(err, catalog.ErrInvalidPrice):
		return "El precio base debe ser mayor o igual a 0."
	case errors.Is(err, catalog.ErrInvalidLocation):
		return "Ubicación inválida: latitud en [-90, 90], longitud en [-180, 180], y ambas requeridas si se especifica una."
	case errors.Is(err, catalog.ErrServiceNotFound):
		return "El servicio no existe."
	case errors.Is(err, catalog.ErrProviderNotFound):
		return "El prestador no existe."
	default:
		return "Ocurrió un error: " + err.Error()
	}
}

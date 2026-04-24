package catalog

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type CategoryHandler struct {
	repo CategoryRepository
}

func NewCategoryHandler(repo CategoryRepository) *CategoryHandler {
	return &CategoryHandler{repo: repo}
}

func (h *CategoryHandler) mountOn(r *gin.RouterGroup) {
	g := r.Group("/categories")
	g.POST("", h.Create)
	g.GET("", h.List)
	g.GET("/:id", h.GetByID)
	g.PUT("/:id", h.Update)
	g.DELETE("/:id", h.Delete)
}

func (h *CategoryHandler) Create(c *gin.Context) {
	var req CreateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		writeBindingError(c, err)
		return
	}
	cat, err := NewCategory(req.Nombre, req.Slug, req.ParentID)
	if err != nil {
		writeError(c, err)
		return
	}
	if err := h.repo.Insert(c.Request.Context(), cat); err != nil {
		writeError(c, err)
		return
	}
	c.JSON(http.StatusCreated, cat)
}

func (h *CategoryHandler) List(c *gin.Context) {
	items, err := h.repo.FindAll(c.Request.Context())
	if err != nil {
		writeError(c, err)
		return
	}
	c.JSON(http.StatusOK, items)
}

func (h *CategoryHandler) GetByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		writeBindingError(c, err)
		return
	}
	cat, err := h.repo.FindByID(c.Request.Context(), id)
	if err != nil {
		writeError(c, err)
		return
	}
	c.JSON(http.StatusOK, cat)
}

func (h *CategoryHandler) Update(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		writeBindingError(c, err)
		return
	}
	var req UpdateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		writeBindingError(c, err)
		return
	}
	existing, err := h.repo.FindByID(c.Request.Context(), id)
	if err != nil {
		writeError(c, err)
		return
	}
	if !slugRegex.MatchString(req.Slug) {
		writeError(c, ErrInvalidSlug)
		return
	}
	existing.Nombre = req.Nombre
	existing.Slug = req.Slug
	existing.ParentID = req.ParentID
	if err := h.repo.Update(c.Request.Context(), existing); err != nil {
		writeError(c, err)
		return
	}
	c.JSON(http.StatusOK, existing)
}

func (h *CategoryHandler) Delete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		writeBindingError(c, err)
		return
	}
	if err := h.repo.Delete(c.Request.Context(), id); err != nil {
		writeError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}

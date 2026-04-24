package catalog

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type PortfolioHandler struct {
	repo PortfolioRepository
}

func NewPortfolioHandler(repo PortfolioRepository) *PortfolioHandler {
	return &PortfolioHandler{repo: repo}
}

func (h *PortfolioHandler) mountOn(r *gin.RouterGroup) {
	g := r.Group("/services/:sid/portfolio")
	g.POST("", h.Create)
	g.GET("", h.List)
	g.PUT("/:id", h.Update)
	g.DELETE("/:id", h.Delete)
}

func (h *PortfolioHandler) Create(c *gin.Context) {
	sid, err := uuid.Parse(c.Param("sid"))
	if err != nil {
		writeBindingError(c, err)
		return
	}
	var req CreatePortfolioItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		writeBindingError(c, err)
		return
	}
	item, err := NewPortfolioItem(sid, req.StorageURL, req.Titulo, req.Orden)
	if err != nil {
		writeError(c, err)
		return
	}
	if err := h.repo.Insert(c.Request.Context(), item); err != nil {
		writeError(c, err)
		return
	}
	c.JSON(http.StatusCreated, item)
}

func (h *PortfolioHandler) List(c *gin.Context) {
	sid, err := uuid.Parse(c.Param("sid"))
	if err != nil {
		writeBindingError(c, err)
		return
	}
	items, err := h.repo.FindByServiceID(c.Request.Context(), sid)
	if err != nil {
		writeError(c, err)
		return
	}
	c.JSON(http.StatusOK, items)
}

func (h *PortfolioHandler) Update(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		writeBindingError(c, err)
		return
	}
	var req UpdatePortfolioItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		writeBindingError(c, err)
		return
	}
	existing, err := h.repo.FindByID(c.Request.Context(), id)
	if err != nil {
		writeError(c, err)
		return
	}
	if req.StorageURL == "" {
		writeError(c, ErrInvalidStorageURL)
		return
	}
	existing.StorageURL = req.StorageURL
	existing.Titulo = req.Titulo
	existing.Orden = req.Orden
	if err := h.repo.Update(c.Request.Context(), existing); err != nil {
		writeError(c, err)
		return
	}
	c.JSON(http.StatusOK, existing)
}

func (h *PortfolioHandler) Delete(c *gin.Context) {
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

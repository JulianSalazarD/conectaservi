package catalog

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ServiceHandler struct {
	repo ServiceRepository
}

func NewServiceHandler(repo ServiceRepository) *ServiceHandler {
	return &ServiceHandler{repo: repo}
}

func (h *ServiceHandler) mountOn(r *gin.RouterGroup) {
	g := r.Group("/services")
	g.POST("", h.Create)
	g.GET("", h.List)
	g.GET("/:id", h.GetByID)
	g.PUT("/:id", h.Update)
	g.DELETE("/:id", h.Delete)
}

func (h *ServiceHandler) Create(c *gin.Context) {
	var req CreateServiceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		writeBindingError(c, err)
		return
	}
	svc, err := NewService(
		req.ProviderID, req.CategoryID,
		req.Titulo, req.Descripcion, req.PrecioBase,
		req.Lat, req.Lng, req.RadioKm,
	)
	if err != nil {
		writeError(c, err)
		return
	}
	if err := h.repo.Insert(c.Request.Context(), svc); err != nil {
		writeError(c, err)
		return
	}
	c.JSON(http.StatusCreated, svc)
}

func (h *ServiceHandler) List(c *gin.Context) {
	filter := ServiceFilter{}
	if v := c.Query("category_id"); v != "" {
		id, err := uuid.Parse(v)
		if err != nil {
			writeBindingError(c, err)
			return
		}
		filter.CategoryID = &id
	}
	if v := c.Query("is_active"); v != "" {
		b, err := strconv.ParseBool(v)
		if err != nil {
			writeBindingError(c, err)
			return
		}
		filter.IsActive = &b
	}
	items, err := h.repo.FindAll(c.Request.Context(), filter)
	if err != nil {
		writeError(c, err)
		return
	}
	c.JSON(http.StatusOK, items)
}

func (h *ServiceHandler) GetByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		writeBindingError(c, err)
		return
	}
	svc, err := h.repo.FindByID(c.Request.Context(), id)
	if err != nil {
		writeError(c, err)
		return
	}
	c.JSON(http.StatusOK, svc)
}

func (h *ServiceHandler) Update(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		writeBindingError(c, err)
		return
	}
	var req UpdateServiceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		writeBindingError(c, err)
		return
	}
	existing, err := h.repo.FindByID(c.Request.Context(), id)
	if err != nil {
		writeError(c, err)
		return
	}
	if req.Titulo == "" {
		writeError(c, ErrInvalidTitle)
		return
	}
	if req.PrecioBase < 0 {
		writeError(c, ErrInvalidPrice)
		return
	}
	if (req.Lat == nil) != (req.Lng == nil) {
		writeError(c, ErrInvalidLocation)
		return
	}
	if req.Lat != nil {
		if _, err := NewLocation(*req.Lat, *req.Lng); err != nil {
			writeError(c, err)
			return
		}
	}
	existing.ProviderID = req.ProviderID
	existing.CategoryID = req.CategoryID
	existing.Titulo = req.Titulo
	existing.Descripcion = req.Descripcion
	existing.PrecioBase = req.PrecioBase
	existing.Lat = req.Lat
	existing.Lng = req.Lng
	existing.RadioKm = req.RadioKm
	existing.IsActive = req.IsActive
	if err := h.repo.Update(c.Request.Context(), existing); err != nil {
		writeError(c, err)
		return
	}
	c.JSON(http.StatusOK, existing)
}

func (h *ServiceHandler) Delete(c *gin.Context) {
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

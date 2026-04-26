package web

import (
	"net/http"
	"net/url"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/JulianSalazarD/conectaservi/internal/catalog"
)

type serviceForm struct {
	ProviderID  string
	CategoryID  string
	Titulo      string
	Descripcion string
	PrecioBase  string
	Lat         string
	Lng         string
	RadioKm     string
}

func (m *Module) servicesList(c *gin.Context) {
	items, err := m.svcRepo.FindAll(c.Request.Context(), catalog.ServiceFilter{})
	if err != nil {
		c.String(http.StatusInternalServerError, "load services: %v", err)
		return
	}
	cats, err := m.catRepo.FindAll(c.Request.Context())
	if err != nil {
		c.String(http.StatusInternalServerError, "load categories: %v", err)
		return
	}
	names := make(map[uuid.UUID]string, len(cats))
	for _, cat := range cats {
		names[cat.ID] = cat.Nombre
	}
	m.render(c, "services_list.html", gin.H{
		"Title":          "Servicios",
		"Items":          items,
		"CategoryNames":  names,
		"Flash":          c.Query("flash"),
	})
}

func (m *Module) serviceNewForm(c *gin.Context) {
	cats, err := m.catRepo.FindAll(c.Request.Context())
	if err != nil {
		c.String(http.StatusInternalServerError, "load categories: %v", err)
		return
	}
	m.render(c, "services_form.html", gin.H{
		"Title":      "Nuevo servicio",
		"Categories": cats,
		"Form":       serviceForm{ProviderID: "00000000-0000-0000-0000-000000000001"},
	})
}

func (m *Module) serviceCreate(c *gin.Context) {
	form := serviceForm{
		ProviderID:  c.PostForm("provider_id"),
		CategoryID:  c.PostForm("category_id"),
		Titulo:      c.PostForm("titulo"),
		Descripcion: c.PostForm("descripcion"),
		PrecioBase:  c.PostForm("precio_base"),
		Lat:         c.PostForm("lat"),
		Lng:         c.PostForm("lng"),
		RadioKm:     c.PostForm("radio_km"),
	}

	providerID, err := uuid.Parse(form.ProviderID)
	if err != nil {
		m.renderServiceForm(c, form, "El UUID del prestador es inválido.")
		return
	}
	categoryID, err := uuid.Parse(form.CategoryID)
	if err != nil {
		m.renderServiceForm(c, form, "Selecciona una categoría.")
		return
	}
	precio, err := strconv.ParseFloat(form.PrecioBase, 64)
	if err != nil {
		m.renderServiceForm(c, form, "El precio base debe ser un número.")
		return
	}
	lat, err := parseOptionalFloat(form.Lat)
	if err != nil {
		m.renderServiceForm(c, form, "Latitud inválida.")
		return
	}
	lng, err := parseOptionalFloat(form.Lng)
	if err != nil {
		m.renderServiceForm(c, form, "Longitud inválida.")
		return
	}
	radio, err := parseOptionalFloat(form.RadioKm)
	if err != nil {
		m.renderServiceForm(c, form, "Radio inválido.")
		return
	}

	svc, err := catalog.NewService(providerID, categoryID, form.Titulo, form.Descripcion, precio, lat, lng, radio)
	if err != nil {
		m.renderServiceForm(c, form, userMessage(err))
		return
	}
	if err := m.svcRepo.Insert(c.Request.Context(), svc); err != nil {
		m.renderServiceForm(c, form, userMessage(err))
		return
	}
	c.Redirect(http.StatusSeeOther, "/web/services?flash=Servicio+creado")
}

func (m *Module) serviceDelete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.String(http.StatusBadRequest, "id inválido")
		return
	}
	if err := m.svcRepo.Delete(c.Request.Context(), id); err != nil {
		c.Redirect(http.StatusSeeOther, "/web/services?flash="+urlEscape(userMessage(err)))
		return
	}
	c.Redirect(http.StatusSeeOther, "/web/services?flash=Servicio+eliminado")
}

func (m *Module) renderServiceForm(c *gin.Context, form serviceForm, errMsg string) {
	cats, err := m.catRepo.FindAll(c.Request.Context())
	if err != nil {
		c.String(http.StatusInternalServerError, "load categories: %v", err)
		return
	}
	m.render(c, "services_form.html", gin.H{
		"Title":      "Nuevo servicio",
		"Categories": cats,
		"Form":       form,
		"Error":      errMsg,
	})
}

func parseOptionalFloat(s string) (*float64, error) {
	if s == "" {
		return nil, nil
	}
	v, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return nil, err
	}
	return &v, nil
}

func urlEscape(s string) string { return url.QueryEscape(s) }

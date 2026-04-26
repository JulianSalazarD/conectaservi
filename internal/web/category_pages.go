package web

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/JulianSalazarD/conectaservi/internal/catalog"
)

type categoryForm struct {
	Nombre string
	Slug   string
}

func (m *Module) categoriesList(c *gin.Context) {
	items, err := m.catRepo.FindAll(c.Request.Context())
	if err != nil {
		c.String(http.StatusInternalServerError, "load categories: %v", err)
		return
	}
	m.render(c, "categories_list.html", gin.H{
		"Title": "Categorías",
		"Items": items,
		"Flash": c.Query("flash"),
	})
}

func (m *Module) categoryNewForm(c *gin.Context) {
	m.render(c, "categories_form.html", gin.H{
		"Title":  "Nueva categoría",
		"Action": "/web/categories",
		"IsEdit": false,
		"Form":   categoryForm{},
	})
}

func (m *Module) categoryCreate(c *gin.Context) {
	form := categoryForm{
		Nombre: c.PostForm("nombre"),
		Slug:   c.PostForm("slug"),
	}
	cat, err := catalog.NewCategory(form.Nombre, form.Slug, nil)
	if err != nil {
		m.renderCategoryForm(c, "/web/categories", false, form, err)
		return
	}
	if err := m.catRepo.Insert(c.Request.Context(), cat); err != nil {
		m.renderCategoryForm(c, "/web/categories", false, form, err)
		return
	}
	c.Redirect(http.StatusSeeOther, "/web/categories?flash=Categor%C3%ADa+creada")
}

func (m *Module) categoryEditForm(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.String(http.StatusBadRequest, "id inválido")
		return
	}
	cat, err := m.catRepo.FindByID(c.Request.Context(), id)
	if err != nil {
		c.String(http.StatusNotFound, userMessage(err))
		return
	}
	m.render(c, "categories_form.html", gin.H{
		"Title":  "Editar categoría",
		"Action": "/web/categories/" + id.String(),
		"IsEdit": true,
		"Form":   categoryForm{Nombre: cat.Nombre, Slug: cat.Slug},
	})
}

func (m *Module) categoryUpdate(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.String(http.StatusBadRequest, "id inválido")
		return
	}
	form := categoryForm{
		Nombre: c.PostForm("nombre"),
		Slug:   c.PostForm("slug"),
	}
	existing, err := m.catRepo.FindByID(c.Request.Context(), id)
	if err != nil {
		c.String(http.StatusNotFound, userMessage(err))
		return
	}
	updated, err := catalog.NewCategory(form.Nombre, form.Slug, existing.ParentID)
	if err != nil {
		m.renderCategoryForm(c, "/web/categories/"+id.String(), true, form, err)
		return
	}
	existing.Nombre = updated.Nombre
	existing.Slug = updated.Slug
	if err := m.catRepo.Update(c.Request.Context(), existing); err != nil {
		m.renderCategoryForm(c, "/web/categories/"+id.String(), true, form, err)
		return
	}
	c.Redirect(http.StatusSeeOther, "/web/categories?flash=Categor%C3%ADa+actualizada")
}

func (m *Module) categoryDelete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.String(http.StatusBadRequest, "id inválido")
		return
	}
	if err := m.catRepo.Delete(c.Request.Context(), id); err != nil {
		c.Redirect(http.StatusSeeOther, "/web/categories?flash="+urlEscape(userMessage(err)))
		return
	}
	c.Redirect(http.StatusSeeOther, "/web/categories?flash=Categor%C3%ADa+eliminada")
}

func (m *Module) renderCategoryForm(c *gin.Context, action string, isEdit bool, form categoryForm, err error) {
	m.render(c, "categories_form.html", gin.H{
		"Title":  "Categoría",
		"Action": action,
		"IsEdit": isEdit,
		"Form":   form,
		"Error":  userMessage(err),
	})
}

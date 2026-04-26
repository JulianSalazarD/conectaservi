package web

import (
	"database/sql"
	"fmt"
	"html/template"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/JulianSalazarD/conectaservi/internal/catalog"
)

// Module wires the HTML/forms web feature on top of the catalog repositories.
type Module struct {
	catRepo catalog.CategoryRepository
	svcRepo catalog.ServiceRepository
	tpl     *template.Template
}

// New builds a Module using Postgres-backed repositories from the catalog package.
func New(db *sql.DB) (*Module, error) {
	tpl, err := parseTemplates()
	if err != nil {
		return nil, fmt.Errorf("parse templates: %w", err)
	}
	return &Module{
		catRepo: catalog.NewPgCategoryRepo(db),
		svcRepo: catalog.NewPgServiceRepo(db),
		tpl:     tpl,
	}, nil
}

// NewWithRepos lets tests inject in-memory repositories.
func NewWithRepos(catRepo catalog.CategoryRepository, svcRepo catalog.ServiceRepository) (*Module, error) {
	tpl, err := parseTemplates()
	if err != nil {
		return nil, fmt.Errorf("parse templates: %w", err)
	}
	return &Module{catRepo: catRepo, svcRepo: svcRepo, tpl: tpl}, nil
}

// Mount registers the web routes (GET/POST) and serves static assets.
func (m *Module) Mount(r *gin.Engine) {
	sub, err := staticSubFS()
	if err == nil {
		r.StaticFS("/static", http.FS(sub))
	}

	r.GET("/", m.home)

	cats := r.Group("/web/categories")
	cats.GET("", m.categoriesList)
	cats.GET("/new", m.categoryNewForm)
	cats.POST("", m.categoryCreate)
	cats.GET("/:id/edit", m.categoryEditForm)
	cats.POST("/:id", m.categoryUpdate)
	cats.POST("/:id/delete", m.categoryDelete)

	svcs := r.Group("/web/services")
	svcs.GET("", m.servicesList)
	svcs.GET("/new", m.serviceNewForm)
	svcs.POST("", m.serviceCreate)
	svcs.POST("/:id/delete", m.serviceDelete)
}

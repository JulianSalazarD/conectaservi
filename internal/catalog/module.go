package catalog

import (
	"database/sql"

	"github.com/gin-gonic/gin"
)

type Module struct {
	categoryHandler  *CategoryHandler
	serviceHandler   *ServiceHandler
	portfolioHandler *PortfolioHandler
}

func New(db *sql.DB) *Module {
	catRepo := NewPgCategoryRepo(db)
	svcRepo := NewPgServiceRepo(db)
	pfRepo := NewPgPortfolioRepo(db)
	return &Module{
		categoryHandler:  NewCategoryHandler(catRepo),
		serviceHandler:   NewServiceHandler(svcRepo),
		portfolioHandler: NewPortfolioHandler(pfRepo),
	}
}

func (m *Module) Mount(r *gin.RouterGroup) {
	m.categoryHandler.mountOn(r)
	m.serviceHandler.mountOn(r)
	m.portfolioHandler.mountOn(r)
}

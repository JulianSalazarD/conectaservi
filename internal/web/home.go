package web

import "github.com/gin-gonic/gin"

func (m *Module) home(c *gin.Context) {
	m.render(c, "home.html", gin.H{"Title": "Inicio"})
}

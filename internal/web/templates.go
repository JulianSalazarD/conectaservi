package web

import (
	"embed"
	"html/template"
	"io/fs"
)

//go:embed templates/*.html
var templatesFS embed.FS

//go:embed static
var staticFS embed.FS

func parseTemplates() (*template.Template, error) {
	return template.ParseFS(templatesFS, "templates/*.html")
}

func staticSubFS() (fs.FS, error) {
	return fs.Sub(staticFS, "static")
}

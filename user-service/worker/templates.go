package worker

import (
	"embed"
	"fmt"
	"html/template"
)

//go:embed templates/*.html
var templateFS embed.FS
var templates *template.Template

func init() {
	t, err := templates.ParseFS(templateFS, "templates/*.html")
	if err != nil {
		panic(fmt.Sprintf("worker: fail to parse email templates: %v", err))
	}
	templates = t
}

package server

import (
	"cyberrange/db"
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"
)

type TemplateConfig struct {
	layoutPath  []string
	includePath []string
	templates   map[string]*template.Template
}

type TemplateDataContext struct {
	User     *User
	Machines *[]db.Machine
}

var (
	loginTemplate    = CreateTemplateConfig("login")
	homeTemplate     = CreateTemplateConfig("home")
	machinesTemplate = CreateTemplateConfig("machines")
	adminTemplate    = CreateTemplateConfig("admin")
)

func checkPath(path []string, err error) []string {
	if err != nil {
		panic(err)
	}
	return path
}

func CreateTemplateConfig(name string) TemplateConfig {
	layoutPath := fmt.Sprintf("web/templates/layouts/%s/*.html", name)
	includePath := fmt.Sprintf("web/templates/includes/%s/*.html", name)

	return TemplateConfig{
		layoutPath:  checkPath(filepath.Glob(layoutPath)),
		includePath: checkPath(filepath.Glob(includePath)),
		templates:   make(map[string]*template.Template),
	}
}

func serveTemplate(w http.ResponseWriter, name string, data interface{}, tmplCfg TemplateConfig) {
	tmpl, ok := tmplCfg.templates[name]
	if !ok {
		http.Error(w, http.StatusText(404), 404)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := tmpl.ExecuteTemplate(w, "base", data); err != nil {
		http.Error(w, http.StatusText(500), 500)
	}
}

func InstantiateTemplates() {
	globalTemplates := []TemplateConfig{
		loginTemplate,
		homeTemplate,
		adminTemplate,
		machinesTemplate,
	}

	for _, tmplCfg := range globalTemplates {
		for _, tmpl := range tmplCfg.includePath {
			files := append(tmplCfg.layoutPath, tmpl)
			tmplCfg.templates[filepath.Base(tmpl)] = template.Must(template.ParseFiles(files...))
		}
	}
}

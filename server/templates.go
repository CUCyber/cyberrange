package main

import (
	"html/template"
	"net/http"
	"path/filepath"
)

type TemplateConfig struct {
	layoutPath  []string
	includePath []string
	templates   map[string]*template.Template
}

var loginTemplate = TemplateConfig{
	layoutPath:  checkPath(filepath.Glob("web/templates/layouts/login/*.html")),
	includePath: checkPath(filepath.Glob("web/templates/includes/login/*.html")),
	templates:   make(map[string]*template.Template),
}

var homeTemplate = TemplateConfig{
	layoutPath:  checkPath(filepath.Glob("web/templates/layouts/home/*.html")),
	includePath: checkPath(filepath.Glob("web/templates/includes/home/*.html")),
	templates:   make(map[string]*template.Template),
}

var adminTemplate = TemplateConfig{
	layoutPath:  checkPath(filepath.Glob("web/templates/layouts/admin/*.html")),
	includePath: checkPath(filepath.Glob("web/templates/includes/admin/*.html")),
	templates:   make(map[string]*template.Template),
}

func checkPath(path []string, err error) []string {
	if err != nil {
		panic(err)
	}
	return path
}

func serveTemplateLogin(w http.ResponseWriter, name string, data interface{}) {
	tmpl, ok := loginTemplate.templates[name]
	if !ok {
		http.Error(w, http.StatusText(404), 404)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := tmpl.ExecuteTemplate(w, "base", data); err != nil {
		http.Error(w, http.StatusText(500), 500)
	}
}

func serveTemplateHome(w http.ResponseWriter, name string, data interface{}) {
	tmpl, ok := homeTemplate.templates[name]
	if !ok {
		http.Error(w, http.StatusText(404), 404)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := tmpl.ExecuteTemplate(w, "base", data); err != nil {
		http.Error(w, http.StatusText(500), 500)
	}
}

func serveTemplateAdmin(w http.ResponseWriter, name string, data interface{}) {
	tmpl, ok := adminTemplate.templates[name]
	if !ok {
		http.Error(w, http.StatusText(404), 404)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := tmpl.ExecuteTemplate(w, "base", data); err != nil {
		http.Error(w, http.StatusText(500), 500)
	}
}

func instantiateTemplates() {
	for _, tmplCfg := range []TemplateConfig{loginTemplate, homeTemplate, adminTemplate} {
		for _, tmpl := range tmplCfg.includePath {
			files := append(tmplCfg.layoutPath, tmpl)
			tmplCfg.templates[filepath.Base(tmpl)] = template.Must(template.ParseFiles(files...))
		}
	}
}

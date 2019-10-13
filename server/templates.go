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

var portalTemplate = TemplateConfig{
	layoutPath:  checkPath(filepath.Glob("web/templates/layouts/portal/*.html")),
	includePath: checkPath(filepath.Glob("web/templates/includes/portal/*.html")),
	templates:   make(map[string]*template.Template),
}

var trainingTemplate = TemplateConfig{
	layoutPath:  checkPath(filepath.Glob("web/templates/layouts/training/*.html")),
	includePath: checkPath(filepath.Glob("web/templates/includes/training/*.html")),
	templates:   make(map[string]*template.Template),
}

var scenarioTemplate = TemplateConfig{
	layoutPath:  checkPath(filepath.Glob("web/templates/layouts/scenario/*.html")),
	includePath: checkPath(filepath.Glob("web/templates/includes/scenario/*.html")),
	templates:   make(map[string]*template.Template),
}

func checkPath(path []string, err error) []string {
	if err != nil {
		panic(err)
	}
	return path
}

func serveTemplatePortal(w http.ResponseWriter, name string) {
	tmpl, ok := portalTemplate.templates[name]
	if !ok {
		http.Error(w, http.StatusText(404), 404)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := tmpl.ExecuteTemplate(w, "base", nil); err != nil {
		http.Error(w, http.StatusText(500), 500)
	}
}

func serveTemplateScenario(w http.ResponseWriter, name string) {
	tmpl, ok := scenarioTemplate.templates[name]
	if !ok {
		http.Error(w, http.StatusText(404), 404)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := tmpl.ExecuteTemplate(w, "base", nil); err != nil {
		http.Error(w, http.StatusText(500), 500)
	}
}

func serveTemplateTraining(w http.ResponseWriter, name string) {
	tmpl, ok := trainingTemplate.templates[name]
	if !ok {
		http.Error(w, http.StatusText(404), 404)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := tmpl.ExecuteTemplate(w, "base", nil); err != nil {
		http.Error(w, http.StatusText(500), 500)
	}
}

func instantiateTemplates() {
	for _, tmplCfg := range []TemplateConfig{portalTemplate, trainingTemplate, scenarioTemplate} {
		for _, tmpl := range tmplCfg.includePath {
			files := append(tmplCfg.layoutPath, tmpl)
			tmplCfg.templates[filepath.Base(tmpl)] = template.Must(template.ParseFiles(files...))
		}
	}
}

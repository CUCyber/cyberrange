package server

import (
	"fmt"
	"github.com/cucyber/cyberrange/services/webserver/db"
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
	User         *User
	Profile      *db.User
	Users        *[]db.User
	Machines     *[]db.Machine
	OwnsTimeline *[]db.OwnMetadata
}

var (
	loginTemplate      = CreateTemplateConfig("login", "login")
	homeTemplate       = CreateTemplateConfig("home", "common")
	profileTemplate    = CreateTemplateConfig("profile", "common")
	scoreboardTemplate = CreateTemplateConfig("scoreboard", "common")
	machinesTemplate   = CreateTemplateConfig("machines", "common")
	adminTemplate      = CreateTemplateConfig("admin", "common")
)

func checkPath(path []string, err error) []string {
	if err != nil {
		panic(err)
	}
	return path
}

func CreateTemplateConfig(include, layout string) TemplateConfig {
	layoutPath := fmt.Sprintf("web/templates/layouts/%s/*.html", layout)
	includePath := fmt.Sprintf("web/templates/includes/%s/*.html", include)

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
		profileTemplate,
		scoreboardTemplate,
		adminTemplate,
		machinesTemplate,
	}

	for _, tmplCfg := range globalTemplates {
		for _, tmpl := range tmplCfg.includePath {
			files := append(tmplCfg.layoutPath, tmpl)

			t, err := template.New("").Funcs(template.FuncMap{
				"inc": func(x int) int {
					return x + 1
				},
				"icon": func(s string) string {
					switch s {
					case "user":
						return "dollar-sign"
					case "root":
						return "hashtag"
					default:
						return "flag"
					}
				},
				"rank": func(user *db.User) uint64 {
					rank, _ := db.GetRank(user)
					return rank
				},
			}).ParseFiles(files...)

			tmplCfg.templates[filepath.Base(tmpl)] = template.Must(t, err)
		}
	}
}

package main

import (
	"fmt"
	"net/http"
	"strings"
	"sync/atomic"
	"time"
)

var routes = []struct {
	path   string
	router func(w http.ResponseWriter, req *http.Request)
}{
	{"/login", c.login},
	{"/logout", c.requiresLogin(c.logout)},
	{"/home", c.requiresLogin(c.home)},
	{"/admin", c.requiresAdmin((c.admin))},
}

func (c *controller) index(w http.ResponseWriter, req *http.Request) {
	for _, route := range routes {
		if strings.HasPrefix(req.URL.Path, route.path) {
			route.router(w, req)
			return
		}
	}

	session, err := store.Get(req, "auth-cookie")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	user := getUser(session)

	if auth := user.Authenticated; !auth {
		err = session.Save(req, w)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		http.Redirect(w, req, "/login", http.StatusFound)
		return
	}

	http.Redirect(w, req, "/home", http.StatusFound)
}

func (c *controller) home(w http.ResponseWriter, req *http.Request) {
	session, err := store.Get(req, "auth-cookie")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	path := strings.TrimPrefix(req.URL.Path, "/home")
	if path == "" {
		path = "index"
	} else {
		path = path[1:]
	}

	user := getUser(session)

	switch req.Method {
	case "GET":
		serveTemplateHome(w, path+".html", user)
	default:
		fmt.Fprintf(w, "Unsupported HTTP option.")
	}
}

func (c *controller) login(w http.ResponseWriter, req *http.Request) {
	session, err := store.Get(req, "auth-cookie")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	switch req.Method {
	case "POST":
		err := req.ParseForm()
		if err != nil {
			serveTemplateLogin(w, "index.html",
				struct{ Error string }{
					"Invalid form data.",
				},
			)
			return
		}

		username := req.PostFormValue("username")
		password := req.PostFormValue("password")

		err = LDAPAuthenticate(username, password)
		if err != nil {
			serveTemplateLogin(w, "index.html",
				struct{ Error string }{
					LDAPError(err),
				},
			)
			return
		}

		user := &User{
			Username:      username,
			Authenticated: true,
			IsAdmin:       LDAPIsAdmin(username),
		}

		session.Values["user"] = user

		err = session.Save(req, w)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		http.Redirect(w, req, "/home", http.StatusSeeOther)
	default:
		serveTemplateLogin(w, "index.html", nil)
	}
}

func (c *controller) admin(w http.ResponseWriter, req *http.Request) {
	session, err := store.Get(req, "auth-cookie")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	path := strings.TrimPrefix(req.URL.Path, "/admin")
	if path == "" {
		path = "index"
	} else {
		path = path[1:]
	}

	user := getUser(session)

	switch req.Method {
	case "GET":
		serveTemplateAdmin(w, path+".html", user)
	default:
		fmt.Fprintf(w, "Unsupported HTTP option.")
	}
}

func (c *controller) logout(w http.ResponseWriter, req *http.Request) {
	session, err := store.Get(req, "auth-cookie")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	session.Values["user"] = User{}
	session.Options.MaxAge = -1

	err = session.Save(req, w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, req, "/", http.StatusFound)
}

func (c *controller) health(w http.ResponseWriter, req *http.Request) {
	if h := atomic.LoadInt64(&c.healthy); h == 0 {
		w.WriteHeader(http.StatusServiceUnavailable)
	} else {
		fmt.Fprintf(w, "uptime: %s\n", time.Since(time.Unix(0, h)))
	}
}

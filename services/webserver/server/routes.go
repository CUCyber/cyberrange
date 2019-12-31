package server

import (
	"fmt"
	"github.com/CUCyber/cyberrange/services/webserver/db"
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
	{"/me", c.requiresLogin(c.me)},
	{"/users", c.requiresLogin(c.users)},
	{"/machines", c.requiresLogin(c.machines)},
	{"/scoreboard", c.requiresLogin(c.scoreboard)},
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

	if auth := user.Authenticated; !auth || user.User == nil {
		err = destroyUserSession(session, w, req)
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

	user := getUser(session)

	data := &TemplateDataContext{
		User: &user,
	}

	switch req.Method {
	case "GET":
		serveTemplate(w, "index.html", data, homeTemplate)
	default:
		fmt.Fprintf(w, "Unsupported HTTP option.")
	}
}

func (c *controller) machines(w http.ResponseWriter, req *http.Request) {
	session, err := store.Get(req, "auth-cookie")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	user := getUser(session)

	switch req.Method {
	case "GET":
		machines, err := db.GetMachines()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		data := &TemplateDataContext{
			Machines: machines,
			User:     &user,
		}

		serveTemplate(w, "index.html", data, machinesTemplate)
	case "POST":
		flag := req.FormValue("flag")
		name := req.FormValue("machine-name")

		err := db.OwnMachine(flag, name, user.User)
		if err != nil {
			w.Write([]byte(err.Error()))
			return
		}

		w.Write([]byte("success"))
	default:
		fmt.Fprintf(w, "Unsupported HTTP option.")
	}
}

func (c *controller) me(w http.ResponseWriter, req *http.Request) {
	session, err := store.Get(req, "auth-cookie")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	user := getUser(session)
	profile := fmt.Sprintf("/users/%d", user.User.Id)

	http.Redirect(w, req, profile, http.StatusFound)
}

func (c *controller) users(w http.ResponseWriter, req *http.Request) {
	session, err := store.Get(req, "auth-cookie")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	user := getUser(session)

	uid, err := ParseUID(req.URL.Path)
	if err != nil {
		http.Redirect(w, req, "/home", http.StatusFound)
	}

	profile, err := db.FindUserById(&db.User{Id: uid})
	if err != nil {
		http.Redirect(w, req, "/home", http.StatusFound)
	}

	timeline, err := db.GetOwnsTimeline(uid)
	if err != nil {
		http.Redirect(w, req, "/home", http.StatusFound)
	}

	switch req.Method {
	case "GET":
		data := &TemplateDataContext{
			User:         &user,
			Profile:      profile,
			OwnsTimeline: timeline,
		}

		serveTemplate(w, "index.html", data, profileTemplate)
	default:
		fmt.Fprintf(w, "Unsupported HTTP option.")
	}
}

func (c *controller) scoreboard(w http.ResponseWriter, req *http.Request) {
	session, err := store.Get(req, "auth-cookie")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	user := getUser(session)

	users, err := db.Scoreboard()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := &TemplateDataContext{
		User:  &user,
		Users: users,
	}

	switch req.Method {
	case "GET":
		serveTemplate(w, "index.html", data, scoreboardTemplate)
	default:
		fmt.Fprintf(w, "Unsupported HTTP option.")
	}
}

func (c *controller) admin(w http.ResponseWriter, req *http.Request) {
	session, err := store.Get(req, "auth-cookie")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	user := getUser(session)

	data := &TemplateDataContext{
		User: &user,
	}

	switch req.Method {
	case "POST":
		err := req.ParseForm()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		data, err := MachineForm(req)
		if err != nil {
			w.Write([]byte(err.Error()))
			return
		}

		_, err = db.FindOrCreateMachine(data)
		if err != nil {
			panic(err.Error())
		}

		w.Write([]byte("success"))
	default:
		serveTemplate(w, "index.html", data, adminTemplate)
	}
}

func (c *controller) login(w http.ResponseWriter, req *http.Request) {
	session, err := store.Get(req, "auth-cookie")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	user := getUser(session)
	if user.Authenticated && user.User != nil {
		http.Redirect(w, req, "/home", http.StatusSeeOther)
		return
	}

	switch req.Method {
	case "POST":
		err := req.ParseForm()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		data, err := LoginForm(req)
		if err != nil {
			serveTemplate(w, "index.html",
				struct{ Error string }{
					"Invalid form data.",
				},
				loginTemplate,
			)
			return
		}

		err = createUserSession(data, w, req)
		if err != nil {
			ldaperr := LDAPError(err)
			if ldaperr != "" {
				serveTemplate(w, "index.html",
					struct{ Error string }{
						ldaperr,
					},
					loginTemplate,
				)
			}
		}

		http.Redirect(w, req, "/home", http.StatusSeeOther)
	default:
		serveTemplate(w, "index.html", nil, loginTemplate)
	}
}

func (c *controller) logout(w http.ResponseWriter, req *http.Request) {
	session, err := store.Get(req, "auth-cookie")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = destroyUserSession(session, w, req)
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

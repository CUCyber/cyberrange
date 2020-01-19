package server

import (
	"fmt"
	"net/http"
	"net/http/pprof"
	"runtime"
	"strings"
	"sync/atomic"
	"time"

	"github.com/cucyber/cyberrange/services/webserver/db"
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
	{"/start", c.requiresLogin(c.start)},
	{"/stop", c.requiresLogin(c.stop)},
	{"/revert", c.requiresLogin(c.revert)},
	{"/list", c.requiresLogin(c.list)},
	{"/chat", c.requiresLogin(c.chat)},
	{"/scoreboard", c.requiresLogin(c.scoreboard)},
	{"/admin/owns", c.requiresAdmin(c.admin_owns)},
	{"/admin/create", c.requiresAdmin(c.admin_create)},
	{"/admin/delete", c.requiresAdmin(c.admin_delete)},
	{"/debug/pprof", c.requiresAdmin(c.profile)},
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

		http.Redirect(w, req, "/login", http.StatusSeeOther)
		return
	}

	http.Redirect(w, req, "/home", http.StatusSeeOther)
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

	path := strings.TrimPrefix(req.URL.Path, "/machines")
	switch path {
	case "/flag":
		conn, err := Websockify(w, req)
		if err != nil {
			return
		}

		go func() {
			defer conn.Close()

			data := struct {
				Flag db.Flag
				Name string
			}{}

			err = ReadJSON(conn, &data)
			if err != nil {
				WriteJSONError(conn, err)
				return
			}

			machine, err := db.FindMachineByName(
				&db.Machine{Name: data.Name},
			)
			if err != nil {
				WriteJSONError(conn, err)
				return
			}

			user, err := db.FindUserById(
				&db.User{Id: user.User.Id},
			)
			if err != nil {
				WriteJSONError(conn, err)
				return
			}

			err = db.OwnMachine(data.Flag, user, machine)
			if err != nil {
				WriteJSONError(conn, err)
				return
			}

			WriteJSONSuccess(conn, "Correct Flag.")
		}()
	default:
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
	}
}

func (c *controller) list(w http.ResponseWriter, req *http.Request) {
	conn, err := Websockify(w, req)
	if err != nil {
		return
	}

	go func() {
		defer conn.Close()

		machines, err := db.GetMachines()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		err = WriteJSON(conn, &machines)
		if err != nil {
			return
		}
	}()
}

func (c *controller) chat(w http.ResponseWriter, req *http.Request) {
	session, err := store.Get(req, "auth-cookie")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	user := getUser(session)

	path := strings.TrimPrefix(req.URL.Path, "/chat")

	switch path {
	case "/ws":
		conn, err := Websockify(w, req)
		if err != nil {
			return
		}

		go func() {
			defer conn.Close()

			ChatHandler(&conn, &user)
		}()
	default:
		data := &TemplateDataContext{
			User: &user,
		}
		serveTemplate(w, "chat.html", data, homeTemplate)
	}
}

func (c *controller) start(w http.ResponseWriter, req *http.Request) {
	conn, err := Websockify(w, req)
	if err != nil {
		return
	}

	go func() {
		defer conn.Close()

		var machine db.Machine
		err = ReadJSON(conn, &machine)
		if err != nil {
			WriteJSONError(conn, err)
			return
		}

		exists, err := db.MachineExists(&machine)
		if err != nil {
			WriteJSONError(conn, err)
			return
		} else if exists == false {
			WriteJSONError(conn, db.ErrMachineNotFound)
			return
		}

		dbMachine, err := db.FindOrCreateMachine(&machine)
		if err != nil {
			WriteJSONError(conn, err)
			return
		}

		if dbMachine.Status == "up" {
			WriteJSONError(conn, db.ErrMachineStarted)
			return
		}

		err = StartMachine(&machine)
		if err != nil {
			WriteJSONError(conn, err)
			return
		}

		WriteJSONSuccess(conn, "Start Machine Request Initiated.")
	}()
}

func (c *controller) stop(w http.ResponseWriter, req *http.Request) {
	conn, err := Websockify(w, req)
	if err != nil {
		return
	}

	go func() {
		defer conn.Close()

		var machine db.Machine
		err = ReadJSON(conn, &machine)
		if err != nil {
			WriteJSONError(conn, err)
			return
		}

		exists, err := db.MachineExists(&machine)
		if err != nil {
			WriteJSONError(conn, err)
			return
		} else if exists == false {
			WriteJSONError(conn, db.ErrMachineNotFound)
			return
		}

		dbMachine, err := db.FindOrCreateMachine(&machine)
		if err != nil {
			WriteJSONError(conn, err)
			return
		}

		if dbMachine.Status == "down" {
			WriteJSONError(conn, db.ErrMachineStopped)
			return
		}

		err = StopMachine(&machine)
		if err != nil {
			WriteJSONError(conn, err)
			return
		}

		WriteJSONSuccess(conn, "Stop Machine Request Initiated.")
	}()
}

func (c *controller) revert(w http.ResponseWriter, req *http.Request) {
	conn, err := Websockify(w, req)
	if err != nil {
		return
	}

	go func() {
		defer conn.Close()

		var machine db.Machine
		err = ReadJSON(conn, &machine)
		if err != nil {
			WriteJSONError(conn, err)
			return
		}

		exists, err := db.MachineExists(&machine)
		if err != nil {
			WriteJSONError(conn, err)
			return
		} else if exists == false {
			WriteJSONError(conn, db.ErrMachineNotFound)
			return
		}

		WriteJSONInfo(conn, "Revert Machine Request Initiated.", 20)

		err = RevertMachine(&machine)
		if err != nil {
			WriteJSONError(conn, err)
			return
		}

		WriteJSONSuccess(conn, "Revert Machine Request Completed.")
	}()
}

func (c *controller) me(w http.ResponseWriter, req *http.Request) {
	session, err := store.Get(req, "auth-cookie")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	user := getUser(session)
	profile := fmt.Sprintf("/users/%d", user.User.Id)

	http.Redirect(w, req, profile, http.StatusSeeOther)
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
		http.Redirect(w, req, "/home", http.StatusSeeOther)
	}

	profile, err := db.FindUserById(&db.User{Id: uid})
	if err != nil {
		http.Redirect(w, req, "/home", http.StatusSeeOther)
	}

	timeline, err := db.GetOwnsTimeline(uid)
	if err != nil {
		http.Redirect(w, req, "/home", http.StatusSeeOther)
	}

	data := &TemplateDataContext{
		User:         &user,
		Profile:      profile,
		OwnsTimeline: timeline,
	}

	serveTemplate(w, "index.html", data, profileTemplate)
}

func (c *controller) scoreboard(w http.ResponseWriter, req *http.Request) {
	session, err := store.Get(req, "auth-cookie")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	user := getUser(session)

	users, err := db.GetUsers()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := &TemplateDataContext{
		User:  &user,
		Users: users,
	}

	serveTemplate(w, "index.html", data, scoreboardTemplate)
}

func (c *controller) admin_owns(w http.ResponseWriter, req *http.Request) {
	path := strings.TrimPrefix(req.URL.Path, "/admin/owns")

	switch path {
	case "/create_user":
		conn, err := Websockify(w, req)
		if err != nil {
			return
		}

		go func() {
			defer conn.Close()

			data := struct {
				Username    string
				MachineName string
			}{}

			err = ReadJSON(conn, &data)
			if err != nil {
				WriteJSONError(conn, err)
				return
			}

			machine, err := db.FindMachineByName(
				&db.Machine{Name: data.MachineName},
			)
			if err != nil {
				WriteJSONError(conn, err)
				return
			}

			user, err := db.FindUserByUsername(
				&db.User{Username: data.Username},
			)
			if err != nil {
				WriteJSONError(conn, err)
				return
			}

			err = db.UserOwnMachine(user, machine)
			if err != nil {
				WriteJSONError(conn, err)
				return
			}

			WriteJSONSuccess(conn, "Created User Own")
		}()
	case "/create_root":
		conn, err := Websockify(w, req)
		if err != nil {
			return
		}

		go func() {
			defer conn.Close()

			data := struct {
				Username    string
				MachineName string
			}{}

			err = ReadJSON(conn, &data)
			if err != nil {
				WriteJSONError(conn, err)
				return
			}

			machine, err := db.FindMachineByName(
				&db.Machine{Name: data.MachineName},
			)
			if err != nil {
				WriteJSONError(conn, err)
				return
			}

			user, err := db.FindUserByUsername(
				&db.User{Username: data.Username},
			)
			if err != nil {
				WriteJSONError(conn, err)
				return
			}

			err = db.RootOwnMachine(user, machine)
			if err != nil {
				WriteJSONError(conn, err)
				return
			}

			WriteJSONSuccess(conn, "Created Root Own")
		}()
	case "/delete_user":
		conn, err := Websockify(w, req)
		if err != nil {
			return
		}

		go func() {
			defer conn.Close()

			data := struct {
				Username    string
				MachineName string
			}{}

			err = ReadJSON(conn, &data)
			if err != nil {
				WriteJSONError(conn, err)
				return
			}

			machine, err := db.FindMachineByName(
				&db.Machine{Name: data.MachineName},
			)
			if err != nil {
				WriteJSONError(conn, err)
				return
			}

			user, err := db.FindUserByUsername(
				&db.User{Username: data.Username},
			)
			if err != nil {
				WriteJSONError(conn, err)
				return
			}

			err = db.DeleteUserOwn(user, machine)
			if err != nil {
				WriteJSONError(conn, err)
				return
			}

			WriteJSONSuccess(conn, "Deleted User Own")
		}()
	case "/delete_root":
		conn, err := Websockify(w, req)
		if err != nil {
			return
		}

		go func() {
			defer conn.Close()

			data := struct {
				Username    string
				MachineName string
			}{}

			err = ReadJSON(conn, &data)
			if err != nil {
				WriteJSONError(conn, err)
				return
			}

			machine, err := db.FindMachineByName(
				&db.Machine{Name: data.MachineName},
			)
			if err != nil {
				WriteJSONError(conn, err)
				return
			}

			user, err := db.FindUserByUsername(
				&db.User{Username: data.Username},
			)
			if err != nil {
				WriteJSONError(conn, err)
				return
			}

			err = db.DeleteRootOwn(user, machine)
			if err != nil {
				WriteJSONError(conn, err)
				return
			}

			WriteJSONSuccess(conn, "Deleted Root Own")
		}()
	default:
		session, err := store.Get(req, "auth-cookie")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		user := getUser(session)

		data := &TemplateDataContext{
			User: &user,
		}

		serveTemplate(w, "owns.html", data, adminTemplate)
	}
}

func (c *controller) admin_create(w http.ResponseWriter, req *http.Request) {
	path := strings.TrimPrefix(req.URL.Path, "/admin/create")

	switch path {
	case "/machine":
		conn, err := Websockify(w, req)
		if err != nil {
			return
		}

		go func() {
			defer conn.Close()

			var machine db.Machine
			err = ReadJSON(conn, &machine)
			if err != nil {
				WriteJSONError(conn, err)
				return
			}

			WriteJSONInfo(conn, "Checking Database for Duplicate Machines", 0)

			exists, err := db.MachineExists(&machine)
			if err != nil {
				WriteJSONError(conn, err)
				return
			} else if exists != false {
				WriteJSONError(conn, db.ErrMachineExists)
				return
			}

			WriteJSONInfo(conn, "Checking Machine Creation", 10)

			err = CheckCreateMachine(&machine)
			if err != nil {
				WriteJSONError(conn, err)
				return
			}

			WriteJSONInfo(conn, "Creating Machine", 20)

			err = CreateMachine(&machine)
			if err != nil {
				WriteJSONError(conn, err)
				return
			}

			WriteJSONInfo(conn, "Creating Initial Machine Snapshot", 70)

			err = SnapshotMachine(&machine)
			if err != nil {
				WriteJSONError(conn, err)
				return
			}

			WriteJSONInfo(conn, "Creating Database Entry", 90)

			_, err = db.FindOrCreateMachine(&machine)
			if err != nil {
				WriteJSONError(conn, err)
				return
			}

			WriteJSONSuccess(conn, "Finished Creating Machine")
		}()
	default:
		session, err := store.Get(req, "auth-cookie")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		user := getUser(session)

		data := &TemplateDataContext{
			User: &user,
		}

		serveTemplate(w, "create.html", data, adminTemplate)
	}
}

func (c *controller) admin_delete(w http.ResponseWriter, req *http.Request) {
	path := strings.TrimPrefix(req.URL.Path, "/admin/delete")
	switch path {
	case "/machine":
		conn, err := Websockify(w, req)
		if err != nil {
			return
		}

		go func() {
			defer conn.Close()

			var machine db.Machine
			err = ReadJSON(conn, &machine)
			if err != nil {
				WriteJSONError(conn, err)
				return
			}

			WriteJSONInfo(conn, "Checking Database for Machine", 0)

			exists, err := db.MachineExists(&machine)
			if err != nil {
				WriteJSONError(conn, err)
				return
			} else if exists == false {
				WriteJSONError(conn, db.ErrMachineNotFound)
				return
			}

			WriteJSONInfo(conn, "Deleting Machine", 50)

			err = db.DeleteMachine(&machine)
			if err != nil {
				WriteJSONError(conn, err)
				return
			}

			WriteJSONInfo(conn, "Reverting User Owns", 75)

			users, err := db.GetUsers()
			if err != nil {
				WriteJSONError(conn, err)
				return
			}

			for i := range *users {
				err = db.ResetOwns(&(*users)[i])
				if err != nil {
					WriteJSONError(conn, err)
					return
				}
			}

			WriteJSONSuccess(conn, "Deleted Machine")
		}()
	default:
		session, err := store.Get(req, "auth-cookie")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		user := getUser(session)

		data := &TemplateDataContext{
			User: &user,
		}

		serveTemplate(w, "delete.html", data, adminTemplate)
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

	http.Redirect(w, req, "/", http.StatusSeeOther)
}

func (c *controller) profile(w http.ResponseWriter, req *http.Request) {
	path := strings.TrimPrefix(req.URL.Path, "/debug/pprof")
	switch path {
	case "/profile":
		pprof.Profile(w, req)
	case "/symbol":
		pprof.Symbol(w, req)
	case "/trace":
		pprof.Trace(w, req)
	default:
		pprof.Index(w, req)
	}
}

func (c *controller) health(w http.ResponseWriter, req *http.Request) {
	if h := atomic.LoadInt64(&c.healthy); h == 0 {
		w.WriteHeader(http.StatusServiceUnavailable)
	} else {
		var m runtime.MemStats
		runtime.ReadMemStats(&m)

		fmt.Fprintf(w, "Uptime: %s\n", time.Since(time.Unix(0, h)))
		fmt.Fprintf(w, "TotalAlloc: %+v\n", m.TotalAlloc)
		fmt.Fprintf(w, "HeapAlloc: %+v\n", m.HeapAlloc)
		fmt.Fprintf(w, "StackInuse: %+v\n", m.StackInuse)
	}
}

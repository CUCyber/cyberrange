package main

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"sync/atomic"
	"time"
)

var routes = []struct {
	path   string
	router func(w http.ResponseWriter, req *http.Request)
}{
	{"/training", c.training},
	{"/scenario", c.scenario},
}

func (c *controller) index(w http.ResponseWriter, req *http.Request) {
	for _, route := range routes {
		if strings.HasPrefix(req.URL.Path, route.path) {
			route.router(w, req)
			return
		}
	}

	if req.URL.Path != "/" {
		http.NotFound(w, req)
		return
	}

	serveTemplatePortal(w, "index.html")
}

func (c *controller) training(w http.ResponseWriter, req *http.Request) {
	path := strings.TrimPrefix(req.URL.Path, "/training")[1:]

	switch req.Method {
	case "GET":
		serveTemplateTraining(w, path+".html")
	case "POST":
		flag := req.FormValue("flag")
		isCorrect := checkFlag(flag, path)
		w.Write([]byte(strconv.FormatBool(isCorrect)))
	default:
		fmt.Fprintf(w, "Unsupported HTTP option.")
	}
}

func (c *controller) scenario(w http.ResponseWriter, req *http.Request) {
	path := strings.TrimPrefix(req.URL.Path, "/scenario")[1:]

	switch req.Method {
	case "GET":
		serveTemplateScenario(w, path+".html")
	default:
		fmt.Fprintf(w, "Unsupported HTTP option.")
	}
}

func (c *controller) health(w http.ResponseWriter, req *http.Request) {
	if h := atomic.LoadInt64(&c.healthy); h == 0 {
		w.WriteHeader(http.StatusServiceUnavailable)
	} else {
		fmt.Fprintf(w, "uptime: %s\n", time.Since(time.Unix(0, h)))
	}
}

var (
	_ http.Handler = http.HandlerFunc((&controller{}).index)
	_ http.Handler = http.HandlerFunc((&controller{}).training)
	_ http.Handler = http.HandlerFunc((&controller{}).health)
)

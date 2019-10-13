package main

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"sync/atomic"
	"time"
)

func (c *controller) index(w http.ResponseWriter, req *http.Request) {
	if strings.HasPrefix(req.URL.Path, "/training") {
		c.training(w, req)
		return
	}

	if strings.HasPrefix(req.URL.Path, "/scenario") {
		c.scenario(w, req)
		return
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
	switch req.Method {
	case "GET":
		path := strings.TrimPrefix(req.URL.Path, "/scenario")
		if path == "" {
			path = "index"
		} else {
			path = path[1:]
		}
		serveTemplateScenario(w, path+".html")
	case "POST":
		path := strings.TrimPrefix(req.URL.Path, "/scenario")
		if path == "" {
			path = "index"
		} else {
			path = path[1:]
		}
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

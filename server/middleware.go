package main

import (
	"net/http"
	"time"
)

type middleware func(http.Handler) http.Handler
type middlewares []middleware
type handlerwrapper func(http.ResponseWriter, *http.Request)

func (mws middlewares) apply(hdlr http.Handler) http.Handler {
	if len(mws) == 0 {
		return hdlr
	}
	return mws[1:].apply(mws[0](hdlr))
}

func (c *controller) requiresLogin(hdlr handlerwrapper) handlerwrapper {
	return func(w http.ResponseWriter, req *http.Request) {
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

		hdlr(w, req)
	}
}

func (c *controller) requiresAdmin(hdlr handlerwrapper) handlerwrapper {
	return func(w http.ResponseWriter, req *http.Request) {
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

		if admin := user.IsAdmin; !admin {
			err = session.Save(req, w)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			http.Redirect(w, req, "/home", http.StatusFound)
		}

		hdlr(w, req)
	}
}

func (c *controller) logging(hdlr http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		defer func(start time.Time) {
			requestID := w.Header().Get("X-Request-Id")
			if requestID == "" {
				requestID = "unknown"
			}
			c.logger.Println(requestID, req.Method, req.URL.Path, req.RemoteAddr, req.UserAgent(), time.Since(start))
		}(time.Now())
		hdlr.ServeHTTP(w, req)
	})
}

func (c *controller) tracing(hdlr http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		requestID := req.Header.Get("X-Request-Id")
		if requestID == "" {
			requestID = c.nextRequestID()
		}
		w.Header().Set("X-Request-Id", requestID)
		hdlr.ServeHTTP(w, req)
	})
}

func (c *controller) restore(hdlr http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				c.logger.Printf("panic: %+v", err)
				http.Error(w, http.StatusText(500), 500)
			}
		}()
		hdlr.ServeHTTP(w, req)
	})
}

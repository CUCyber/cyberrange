package server

import (
	"encoding/gob"
	"net/http"

	"github.com/cucyber/cyberrange/services/webserver/db"
	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
)

type User struct {
	User          *db.User
	Authenticated bool
	IsAdmin       bool
}

var store *sessions.CookieStore

func getUser(s *sessions.Session) User {
	val := s.Values["user"]
	var user = User{}
	user, ok := val.(User)
	if !ok {
		return User{Authenticated: false}
	}
	return user
}

func destroyUserSession(s *sessions.Session, w http.ResponseWriter, req *http.Request) error {
	s.Values["user"] = User{}
	s.Options.MaxAge = -1

	err := s.Save(req, w)
	if err != nil {
		return err
	}

	return nil
}

func createUserSession(username, password string, w http.ResponseWriter, req *http.Request) error {
	err := LDAPAuthenticate(username, password)
	if err != nil {
		return err
	}

	user, err := db.FindOrCreateUser(
		&db.User{
			Username: username,
			Points:   0,
		},
	)
	if err != nil {
		return err
	}

	session, err := store.Get(req, "auth-cookie")
	if err != nil {
		return err
	}

	session.Values["user"] = &User{
		User:          user,
		Authenticated: true,
		IsAdmin:       LDAPIsAdmin(username),
	}

	err = session.Save(req, w)
	if err != nil {
		return err
	}

	return nil
}

func InitializeSessions() {
	authKey := securecookie.GenerateRandomKey(64)
	encryptionKey := securecookie.GenerateRandomKey(32)

	store = sessions.NewCookieStore(
		authKey,
		encryptionKey,
	)

	store.Options = &sessions.Options{
		MaxAge:   60 * 60,
		HttpOnly: true,
	}

	gob.Register(User{})
}

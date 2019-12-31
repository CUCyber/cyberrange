package server

import (
	"github.com/CUCyber/cyberrange/services/webserver/db"
	"encoding/gob"
	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
	"net/http"
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

func createUserSession(data *LoginFormData, w http.ResponseWriter, req *http.Request) error {
	err := LDAPAuthenticate(data.Username, data.Password)
	if err != nil {
		return err
	}

	user, err := db.FindOrCreateUser(
		&db.User{
			Username: data.Username,
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
		IsAdmin:       LDAPIsAdmin(data.Username),
	}

	err = session.Save(req, w)
	if err != nil {
		return err
	}

	return nil
}

func InitializeSessions() {
	authKey, _ := []byte("ABCDABCDABCDABCD"), securecookie.GenerateRandomKey(64)
	encryptionKey, _ := []byte("ABCDABCDABCDABCD"), securecookie.GenerateRandomKey(32)

	store = sessions.NewCookieStore(
		authKey,
		encryptionKey,
	)

	store.Options = &sessions.Options{
		MaxAge:   60 * 15,
		HttpOnly: true,
	}

	gob.Register(User{})
}

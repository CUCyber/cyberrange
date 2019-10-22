package server

import (
	"encoding/gob"
	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
)

type User struct {
	Username      string
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

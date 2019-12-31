package server

import (
	"github.com/CUCyber/cyberrange/services/webserver/db"
	"errors"
	"net/http"
)

type LoginFormData struct {
	Username string
	Password string
}

var (
	ErrInvalidFormData = errors.New("cyberrange: form data is invalid")
)

func MachineForm(req *http.Request) (*db.Machine, error) {
	name := req.PostFormValue("name")
	points := req.PostFormValue("points")
	difficulty := req.PostFormValue("difficulty")
	userflag := req.PostFormValue("userflag")
	rootflag := req.PostFormValue("rootflag")

	if name == "" {
		return nil, ErrInvalidFormData
	}

	pointval, err := ConvUint(points)
	if err != nil {
		return nil, ErrInvalidFormData
	}

	validDifficulty := false
	for _, val := range db.MachineDifficulty {
		if difficulty == val {
			validDifficulty = true
			break
		}
	}
	if !validDifficulty {
		return nil, ErrInvalidFormData
	}

	if userflag == "" {
		return nil, ErrInvalidFormData
	}

	if rootflag == "" {
		return nil, ErrInvalidFormData
	}

	return &db.Machine{
		Name:       name,
		Points:     pointval,
		Difficulty: difficulty,
		UserFlag:   userflag,
		RootFlag:   rootflag,
	}, nil
}

func LoginForm(req *http.Request) (*LoginFormData, error) {
	username := req.PostFormValue("username")
	password := req.PostFormValue("password")

	if username == "" {
		return nil, ErrInvalidFormData
	}

	if password == "" {
		return nil, ErrInvalidFormData
	}

	return &LoginFormData{
		Username: username,
		Password: password,
	}, nil
}

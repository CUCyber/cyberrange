package db

import "github.com/gocraft/dbr/v2"

type User struct {
	Id       uint64
	Username string
	Points   uint64
}

func Scoreboard() (*[]User, error) {
	var users []User

	_, err := db.Select("*").From("users").
		OrderDesc("points").
		Load(&users)
	if err != nil {
		return nil, err
	}

	return &users, nil
}

func GetUsers() (*[]User, error) {
	var users []User

	_, err := db.Select("*").From("users").
		Load(&users)
	if err != nil {
		return nil, err
	}

	return &users, nil
}

func FindUserById(user *User) (*User, error) {
	var queryUser User

	_, err := db.Select("*").From("users").
		Where("id = ?", user.Id).
		Load(&queryUser)
	if err != nil {
		return nil, err
	}

	if queryUser.Username == "" {
		return nil, dbr.ErrNotFound
	}

	return &queryUser, nil
}

func FindUserByUsername(user *User) (*User, error) {
	var queryUser User

	_, err := db.Select("*").From("users").
		Where("username = ?", user.Username).
		Load(&queryUser)
	if err != nil {
		return nil, err
	}

	if queryUser.Username == "" {
		return nil, dbr.ErrNotFound
	}

	return &queryUser, nil
}

func FindOrCreateUser(user *User) (*User, error) {
	/*
	   While insert...returning isn't available, we need
	   to get the updated user id through a second query

	   https://mariadb.com/kb/en/library/insertreturning/
	*/

	queryUser, err := FindUserByUsername(user)
	if err != nil {
		/* Unexpected error */
		if err != dbr.ErrNotFound {
			return nil, err
		}

		/* User not found, create user */
		_, err := db.InsertInto("users").
			Columns("username", "points").
			Record(user).
			Exec()
		if err != nil {
			return nil, err
		}

		/* Get new user data */
		queryUser, err = FindUserByUsername(user)
		if err != nil {
			if err != dbr.ErrNotFound {
				return nil, err
			}
		}

		/* Pre-populate machine owns */
		err = UserCreateMachineOwns(queryUser.Id)
		if err != nil {
			panic(err.Error())
		}
	}

	return queryUser, nil
}

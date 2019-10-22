package db

import (
	"github.com/jinzhu/gorm"
)

type User struct {
	gorm.Model
	Username string
}

func FindUser(user *User) (*User, error) {
	result := db.Find(&user, user)
	if result.Error != nil {
		return nil, result.Error
	}
	return user, nil
}

func FirstOrCreateUser(user *User) *User {
	db.FirstOrCreate(&user, user)
	return user
}

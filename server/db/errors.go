package db

import "errors"

var (
	ErrUserNotFound    = errors.New("cyberrange: user not found")
	ErrMachineNotFound = errors.New("cyberrange: machine not found")
)

var (
	ErrWrongFlag = errors.New("cyberrange: flag submission is incorrect")
	ErrUserOwned = errors.New("cyberrange: user has already been owned by this user")
	ErrRootOwned = errors.New("cyberrange: root has already been owned by this user")
	ErrBothOwned = errors.New("cyberrange: user and root have already been owned by this user")
)

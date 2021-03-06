package db

import "errors"

var (
	ErrMachineStarted  = errors.New("cyberrange: machine is already active")
	ErrMachineStopped  = errors.New("cyberrange: machine is already inactive")
	ErrUserNotFound    = errors.New("cyberrange: user not found")
	ErrMachineNotFound = errors.New("cyberrange: machine not found")
	ErrMachineExists   = errors.New("cyberrange: machine already exists")
)

var (
	ErrWrongFlag    = errors.New("cyberrange: flag submission is incorrect")
	ErrUserOwned    = errors.New("cyberrange: user has already been owned by this user")
	ErrRootOwned    = errors.New("cyberrange: root has already been owned by this user")
	ErrBothOwned    = errors.New("cyberrange: user and root have already been owned by this user")
	ErrUserNotOwned = errors.New("cyberrange: user has not been owned by this user")
	ErrRootNotOwned = errors.New("cyberrange: root has not been owned by this user")
)

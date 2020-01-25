package server

import (
	"errors"
	"sync"
	"time"

	"github.com/cucyber/cyberrange/services/webserver/db"
)

var (
	ErrActionExists   = errors.New("cyberrange: action request already exists")
	ErrActionNotFound = errors.New("cyberrange: action request not found")
)

type ActionType uint8

const (
	None ActionType = iota
	Stop
	Revert
	Cancel
)

type Action struct {
	action  ActionType
	channel chan ActionType
}

type ActionList struct {
	*sync.Map
}

var al *ActionList

func (al *ActionList) Execute(machine *db.Machine, action *Action) {
	switch action.action {
	case Revert:
		RevertMachine(machine)
	case Stop:
		StopMachine(machine)
	}
}

func (al *ActionList) ActionWait(machine *db.Machine, action *Action) {
	for {
		select {
		case status := <-action.channel:
			action.action = status
			al.Map.Store(machine.Name, action)
		}

		select {
		case <-action.channel:
		case <-time.After(30 * time.Second):
			al.Execute(machine, action)
		}

		action.action = None
		al.Map.Store(machine.Name, action)
	}
}

func (al *ActionList) CancelAction(machine *db.Machine) (string, error) {
	action, exists := al.Map.Load(machine.Name)
	if exists == false || action.(*Action).action == None {
		return "", ErrActionNotFound
	}

	var cancelled string
	switch action.(*Action).action {
	case Stop:
		cancelled = "stop"
	case Revert:
		cancelled = "revert"
	default:
		cancelled = "action"
	}

	action.(*Action).channel <- Cancel

	return cancelled, nil
}

func (al *ActionList) CreateAction(machine *db.Machine, newAction ActionType) error {
	action, _ := al.Map.Load(machine.Name)
	if action != nil && action.(*Action).action != None {
		return ErrActionExists
	}

	if action == nil {
		action = &Action{
			action:  newAction,
			channel: make(chan ActionType),
		}
		al.Map.Store(machine.Name, action)
		go al.ActionWait(machine, action.(*Action))
	}

	action.(*Action).channel <- newAction

	return nil
}

func ClearActionList() {
	al.Map.Range(func(key, value interface{}) bool {
		al.Map.Delete(key)
		return true
	})
}

func InitializeActionList() {
	al = &ActionList{&sync.Map{}}
}

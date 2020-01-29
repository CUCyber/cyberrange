package server

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/cucyber/cyberrange/services/webserver/db"
)

type Command struct {
	command   string
	arguments []string
}

var (
	ErrNotCommand = errors.New("cyberrange: message is not a command")
)

func (cmd *Command) Cancel(mc *MessageContext) {
	mc.Message.Auto = true
	mc.Message.Time = "!"

	if len(cmd.arguments) < 1 {
		mc.Message.Data = "/cancel <MachineName>"
		hub.PrivateMessage(mc)
		return
	}

	machine := &db.Machine{
		Name: cmd.arguments[0],
	}

	exists, err := db.MachineExists(machine)
	if err != nil {
		mc.Message.Data = err.Error()
		hub.PrivateMessage(mc)
		return
	} else if exists == false {
		mc.Message.Data = db.ErrMachineNotFound.Error()
		hub.PrivateMessage(mc)
		return
	}

	event, err := al.CancelAction(machine)
	if err != nil {
		mc.Message.Data = err.Error()
		hub.PrivateMessage(mc)
		return
	}

	mc.Message.Time = time.Now().UTC().Format("Jan 02 15:04")
	mc.Message.Data = fmt.Sprintf(
		"%s cancelled %s for %s.",
		mc.ChatConn.User.User.Username,
		event, machine.Name,
	)
	hub.BroadcastMessage(mc.Message)
}

func (cmd *Command) Help(mc *MessageContext) {
	mc.Message.Auto = true
	mc.Message.Time = "?"

	mc.Message.Data = "Available Commands:"
	hub.PrivateMessage(mc)

	mc.Message.Data = "/help" +
		" | " +
		"Displays this help screen"
	hub.PrivateMessage(mc)

	mc.Message.Data = "/cancel <MachineName>" +
		" | " +
		"Cancels a stop or revert on <MachineName>"
	hub.PrivateMessage(mc)
}

func ParseCommand(input string) (*Command, error) {
	command := &Command{}

	if !strings.HasPrefix(input, "/") {
		return nil, ErrNotCommand
	}

	input = strings.TrimSpace(strings.TrimPrefix(input, "/"))
	if input == "" {
		return nil, ErrNotCommand
	}

	parts := strings.SplitN(input, " ", 2)

	command.command = strings.ToLower(parts[0])
	if len(parts) > 1 {
		command.arguments = strings.Fields(parts[1])
	}

	return command, nil
}

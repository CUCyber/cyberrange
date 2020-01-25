package server

import (
	"container/ring"
	"encoding/json"
	"fmt"
	"net"
	"strings"
	"sync"
	"time"

	"github.com/cucyber/cyberrange/services/webserver/db"
)

var hub *ChatHub
var mq *MessageQueue

type ChatConn struct {
	User *User
	Conn *net.Conn
}

type ChatHub struct {
	Clients          map[*User]*net.Conn
	CreateClientChan chan *ChatConn
	DeleteClientChan chan *ChatConn
	BroadcastMsgChan chan *MessageContext
	CloseChatHubChan chan int
}

type MessageQueue struct {
	*ring.Ring
	*sync.Mutex
}

type Message struct {
	Global bool   `json:"Global"`
	Auto   bool   `json:"Auto"`
	Time   string `json:"Time"`
	Name   string `json:"Name"`
	Data   string `json:"Data"`
}

type MessageContext struct {
	*Message
	*ChatConn
}

type Chat struct {
	Messages []*Message `json:"Messages"`
}

func (mq *MessageQueue) MarshalJSON() ([]byte, error) {
	mq.Mutex.Lock()
	defer mq.Mutex.Unlock()

	var chat Chat
	for j := 0; j < mq.Len(); j++ {
		chat.Messages = append(chat.Messages, mq.Value.(*Message))
		mq.Ring = mq.Next()
	}
	return json.Marshal(&chat)
}

func (mq *MessageQueue) Add(v *Message) {
	mq.Mutex.Lock()
	defer mq.Mutex.Unlock()

	mq.Ring.Value = v
	mq.Ring = mq.Next()
}

func (mq *MessageQueue) init(size int) *MessageQueue {
	mq.Ring = ring.New(size)
	mq.Mutex = new(sync.Mutex)
	for i := 0; i < size; i++ {
		mq.Ring.Value = &Message{}
		mq.Ring = mq.Ring.Next()
	}
	return mq
}

func NewMessageQueue(size int) *MessageQueue {
	return new(MessageQueue).init(size)
}

func CreateChatBroadcast(message string) {
	m := &Message{
		Global: true,
		Auto:   true,
		Time:   time.Now().UTC().Format("Jan 02 15:04"),
		Name:   "",
		Data:   message,
	}
	hub.BroadcastMessage(m)
}

func (h *ChatHub) HandleMessage(mc *MessageContext) {
	var cmd string
	var args []string

	if !strings.HasPrefix(mc.Message.Data, "/") {
		hub.BroadcastMessage(mc.Message)
		return
	}

	message := strings.TrimSpace(strings.TrimPrefix(mc.Message.Data, "/"))
	if message == "" {
		hub.BroadcastMessage(mc.Message)
		return
	}

	parts := strings.SplitN(message, " ", 2)

	cmd = strings.ToLower(parts[0])
	if len(parts) > 1 {
		args = strings.Fields(parts[1])
	}

	switch cmd {
	case "help":
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
	case "cancel":
		mc.Message.Auto = true
		mc.Message.Time = "!"

		if len(args) < 1 {
			mc.Message.Data = "/cancel <MachineName>"
			hub.PrivateMessage(mc)
			return
		}

		machine := &db.Machine{Name: args[0]}

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
			event, args[0],
		)
		hub.BroadcastMessage(mc.Message)
	}
}

func (h *ChatHub) PrivateMessage(mc *MessageContext) {
	err := WriteJSON(*mc.ChatConn.Conn, mc.Message)
	if err != nil {
		return
	}
}

func (h *ChatHub) BroadcastMessage(m *Message) {
	for _, conn := range h.Clients {
		err := WriteJSON(*conn, *m)
		if err != nil {
			return
		}
	}

	go mq.Add(m)
}

func (h *ChatHub) DeleteClient(conn *ChatConn) {
	delete(h.Clients, conn.User)
}

func (h *ChatHub) CreateClient(conn *ChatConn) {
	h.Clients[conn.User] = conn.Conn

	err := WriteJSON(*conn.Conn, mq)
	if err != nil {
		return
	}
}

func (h *ChatHub) run() {
	for {
		select {
		case conn := <-h.CreateClientChan:
			h.CreateClient(conn)
		case conn := <-h.DeleteClientChan:
			h.DeleteClient(conn)
		case mesg := <-h.BroadcastMsgChan:
			h.HandleMessage(mesg)
		case <-h.CloseChatHubChan:
			return
		}
	}
}

func ChatHandler(ws *net.Conn, user *User) {
	go hub.run()

	ClientConnection := &ChatConn{user, ws}
	hub.CreateClientChan <- ClientConnection

	for {
		var m Message

		err := ReadJSON(*ws, &m)
		if err != nil {
			hub.DeleteClient(&ChatConn{user, ws})
			return
		}

		hub.BroadcastMsgChan <- &MessageContext{
			&Message{
				Auto: false,
				Time: time.Now().UTC().Format("Jan 02 15:04"),
				Name: user.User.Username,
				Data: m.Data,
			},
			ClientConnection,
		}
	}
}

func CloseChatHub() {
	hub.CloseChatHubChan <- 1

	for client := range hub.Clients {
		delete(hub.Clients, client)
	}
}

func InitializeChatHub() {
	hub = &ChatHub{
		Clients:          make(map[*User]*net.Conn),
		CreateClientChan: make(chan *ChatConn),
		DeleteClientChan: make(chan *ChatConn),
		BroadcastMsgChan: make(chan *MessageContext),
		CloseChatHubChan: make(chan int),
	}
	mq = NewMessageQueue(16)
}

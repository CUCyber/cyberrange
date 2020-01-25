package server

import (
	"net"
	"time"
)

var hub *ChatHub

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

func (h *ChatHub) PrivateMessage(mc *MessageContext) {
	err := WriteJSON(*mc.ChatConn.Conn, mc.Message)
	if err != nil {
		return
	}
}

func (h *ChatHub) HandleMessage(mc *MessageContext) {
	command, err := ParseCommand(mc.Message.Data)
	if err != nil {
		hub.BroadcastMessage(mc.Message)
		return
	}

	switch command.command {
	case "help":
		command.Help(mc)
	case "cancel":
		command.Cancel(mc)
	}
}

func (h *ChatHub) BroadcastMessage(m *Message) {
	for _, conn := range h.Clients {
		err := WriteJSON(*conn, *m)
		if err != nil {
			return
		}
	}

	go ch.Add(m)
}

func (h *ChatHub) DeleteClient(conn *ChatConn) {
	delete(h.Clients, conn.User)
}

func (h *ChatHub) CreateClient(conn *ChatConn) {
	h.Clients[conn.User] = conn.Conn

	err := WriteJSON(*conn.Conn, ch)
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
	ch = NewChatHistory(16)
}

package server

import (
	"container/ring"
	"encoding/json"
	"net"
	"sync"
	"time"
)

var hub *ChatHub
var mq *MessageQueue

type ChatConn struct {
	user *User
	conn *net.Conn
}

type ChatHub struct {
	Clients          map[*User]*net.Conn
	CreateClientChan chan *ChatConn
	DeleteClientChan chan *ChatConn
	BroadcastMsgChan chan Message
	CloseChatHubChan chan int
}

type MessageQueue struct {
	*ring.Ring
	*sync.Mutex
}

type Message struct {
	Time string `json:"Time"`
	Name string `json:"Name"`
	Data string `json:"Data"`
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

func (mq *MessageQueue) Add(v interface{}) {
	mq.Mutex.Lock()
	defer mq.Mutex.Unlock()

	mq.Ring.Value = v
	mq.Ring = mq.Next()
}

func NewMessageQueue() *MessageQueue {
	q := &MessageQueue{ring.New(16), new(sync.Mutex)}
	for i := 0; i < q.Len(); i++ {
		q.Ring.Value = &Message{}
		q.Ring = q.Next()
	}
	return q
}

func (h *ChatHub) BroadcastMessage(m Message) {
	for _, conn := range h.Clients {
		err := WriteJSON(*conn, m)
		if err != nil {
			return
		}
	}

	mq.Add(&m)
}

func (h *ChatHub) DeleteClient(conn *ChatConn) {
	delete(h.Clients, conn.user)
}

func (h *ChatHub) CreateClient(conn *ChatConn) {
	h.Clients[conn.user] = conn.conn

	err := WriteJSON(*conn.conn, mq)
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
			h.BroadcastMessage(mesg)
		case <-h.CloseChatHubChan:
			return
		}
	}
}

func ChatHandler(ws *net.Conn, user *User) {
	go hub.run()

	hub.CreateClientChan <- &ChatConn{user, ws}

	for {
		var m Message
		err := ReadJSON(*ws, &m)
		if err != nil {
			hub.DeleteClient(&ChatConn{user, ws})
			return
		}
		m.Time = time.Now().UTC().Format("Jan 02 15:04")
		m.Name = user.User.Username
		hub.BroadcastMsgChan <- m
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
		BroadcastMsgChan: make(chan Message),
		CloseChatHubChan: make(chan int),
	}
	mq = NewMessageQueue()
}

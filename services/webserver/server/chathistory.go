package server

import (
	"container/ring"
	"encoding/json"
	"sync"
)

var ch *ChatHistory

type ChatHistory struct {
	*ring.Ring
	*sync.Mutex
}

func (ch *ChatHistory) MarshalJSON() ([]byte, error) {
	ch.Mutex.Lock()
	defer ch.Mutex.Unlock()

	var chat Chat
	for j := 0; j < ch.Len(); j++ {
		chat.Messages = append(chat.Messages, ch.Value.(*Message))
		ch.Ring = ch.Next()
	}
	return json.Marshal(&chat)
}

func (ch *ChatHistory) Add(v *Message) {
	ch.Mutex.Lock()
	defer ch.Mutex.Unlock()

	ch.Ring.Value = v
	ch.Ring = ch.Next()
}

func (ch *ChatHistory) init(size int) *ChatHistory {
	ch.Ring = ring.New(size)
	ch.Mutex = new(sync.Mutex)
	for i := 0; i < size; i++ {
		ch.Ring.Value = &Message{}
		ch.Ring = ch.Ring.Next()
	}
	return ch
}

func NewChatHistory(size int) *ChatHistory {
	return new(ChatHistory).init(size)
}

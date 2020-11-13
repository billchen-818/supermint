package message

import (
	"fmt"
	"github.com/vbhp/supermint/libs/log"
	"time"
)

type Message struct {
	log log.Logger
}

func NewMessage(logger log.Logger) Message {
	message := Message{
		log: logger,
	}

	return message
}

func (m *Message) Start() {
	for {
		fmt.Print("hello message! ")
		time.Sleep(time.Second * 5)
	}
}

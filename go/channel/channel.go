package channel

import (
	"time"
)

// Channel 管理多个聊天室
type Channel interface {
	CreateRoom(name string) bool
	ExitRoom(roomname string)
	GetInfo(roomname string) (history []string, removeTime time.Time, exist bool)
	Init()
	SendMessage(roomname string, message string) bool
}

var C = New()

func New() Channel {
	return &pubsub_channel{}
}

func init() {
	C.Init()
}

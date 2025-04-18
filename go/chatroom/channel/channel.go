// Package channal 基于redis实现聊天室
package channel

import (
	"time"

	"gitee.com/qiulaidongfeng/chatroom/go/chatroom/internal/config"
)

// Channel 管理多个聊天室
type Channel interface {
	// CreateRoom 创建一个聊天室
	CreateRoom(name string) bool
	// EnterRoom 进入聊天室
	EnterRoom(roonname string, expire time.Duration) (id string)
	// SetIdExpire 设置id的过期时间为当前时间经过一定时期
	SetIdExpire(roomname, id string, expire time.Duration)
	// ExitRoom 退出聊天室
	ExitRoom(roomname, id string)
	// GetInfo 获取聊天室的信息
	GetInfo(roomname string, id string) (history []string, ttl time.Duration, exist bool, online int64)
	// Init 进行连接数据库之类的初始化
	Init()
	// SendMessage 发送一条消息到聊天室
	SendMessage(roomname string, message string) bool
	// waitMessage 等待任意聊天室收到一条消息
	// 用于测试 仅在test==true时可用
	waitMessage()
}

var C = New()

func New() Channel {
	if config.GetMode() == 2 {
		return &list_channel{}
	}
	return &pubsub_channel{}
}

func init() {
	C.Init()
}

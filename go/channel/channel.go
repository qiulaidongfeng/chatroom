package channel

import (
	"fmt"
	"os"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

// Channel 管理多个聊天室
type Channel interface {
	// CreateRoom 创建一个聊天室
	CreateRoom(name string) bool
	// ExitRoom 退出聊天室
	ExitRoom(roomname string)
	// GetInfo 获取聊天室的信息
	GetInfo(roomname string) (history []string, ttl time.Duration, exist bool)
	// Init 进行连接数据库之类的初始化
	Init()
	// SendMessage 发送一条消息到聊天室
	SendMessage(roomname string, message string) bool
	// waitMessage 等待任意聊天室收到一条消息
	// 用于测试 仅在test==true时可用
	waitMessage()
}

var C = New()

var v *viper.Viper = func() *viper.Viper {
	v := viper.New()
	prefix := ""
	if Test {
		prefix = "../"
	}
	v.SetConfigFile(prefix + "config.ini")
	v.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("Config file changed:", e.Name)
	})
	v.WatchConfig()
	err := v.ReadInConfig()
	if err != nil {
		panic(err)
	}
	return v
}()

var Test bool = os.Getenv("TEST") != ""

func New() Channel {
	if v.GetInt("chatroom.mode") == 2 {
		return &list_channel{}
	}
	return &pubsub_channel{}
}

func init() {
	C.Init()
}

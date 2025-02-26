package channel

import (
	"context"
	"sync"

	"github.com/redis/go-redis/v9"
)

// Channel 管理多个聊天室
type Channel struct {
	rdb  *redis.Client
	all  sync.Map //map[string]*room
	lock sync.Mutex
}

// room 聊天室
// TODO:支持多个用户进入并退出聊天室
type room struct {
	pubsub  *redis.PubSub
	history []string
	lock    sync.Mutex
}

var seam = make(chan struct{})

var c Channel

func init() {
	c.rdb = redis.NewClient(&redis.Options{Addr: "127.0.0.1:6379", DB: 15})
	if err := c.rdb.Ping(context.Background()).Err(); err != nil {
		panic(err)
	}
}

// CreateRoom 创建一个聊天室
func CreateRoom(name string) {
	if _, ok := c.all.Load(name); ok {
		return
	}
	c.lock.Lock()
	defer c.lock.Unlock()
	r := &room{}
	r.pubsub = c.rdb.Subscribe(context.Background(), name)
	_, err := r.pubsub.Receive(context.Background())
	if err != nil {
		panic(err)
	}
	go func() {
		c := r.pubsub.Channel()
		for {
			m := <-c
			if test {
				seam <- struct{}{}
			}
			r.lock.Lock()
			r.history = append(r.history, m.Payload)
			r.lock.Unlock()
		}
	}()
	c.all.Store(name, r)
}

var test bool

// GetHistory 获取聊天室的历史消息
func GetHistory(roomname string) []string {
	v, ok := c.all.Load(roomname)
	if !ok {
		return nil
	}
	r := v.(*room)
	r.lock.Lock()
	defer r.lock.Unlock()
	return r.history
}

// SendMessage 发送一条消息到聊天室
func SendMessage(roomname string, message string) bool {
	_, ok := c.all.Load(roomname)
	if !ok {
		return false
	}
	for i := range 10 {
		err := c.rdb.Publish(context.Background(), roomname, message).Err()
		if err != nil {
			continue
		}
		if i == 9 {
			panic(err)
		}
		break
	}
	return true
}

// waitMessage 等待任意聊天室收到一条消息
// 用于测试 仅在test==true时可用
func waitMessage() {
	<-seam
}

// ExitRoom 退出聊天室
func ExitRoom(roomname string) {
	v, ok := c.all.Load(roomname)
	if !ok {
		return
	}
	r := v.(*room)
	r.pubsub.Close()
	r.pubsub.Unsubscribe(context.Background(), roomname)
}

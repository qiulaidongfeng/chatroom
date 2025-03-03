package channel

import (
	"context"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
)

var _ Channel = (*pubsub_channel)(nil)

// pubsub_channel 管理基于redis发布订阅实现的聊天室
type pubsub_channel struct {
	rdb  *redis.Client
	all  sync.Map //map[string]*room
	lock sync.Mutex
}

// pubsub_room 聊天室
// TODO:支持多个用户进入并退出聊天室
type pubsub_room struct {
	pubsub     *redis.PubSub
	history    []string
	removeTime time.Time
	t          *time.Timer
	//Note:这里故意不用读写锁，因为一个聊天室的并发量不会很大
	lock sync.Mutex
}

var seam = make(chan struct{})

func (c *pubsub_channel) Init() {
	c.rdb = redis.NewClient(&redis.Options{Addr: "127.0.0.1:6379", DB: 15})
	if err := c.rdb.Ping(context.Background()).Err(); err != nil {
		panic(err)
	}
}

func (c *pubsub_channel) CreateRoom(name string) bool {
	if _, ok := c.all.Load(name); ok {
		return false
	}
	c.lock.Lock()
	defer c.lock.Unlock()
	r := &pubsub_room{}
	r.pubsub = c.rdb.Subscribe(context.Background(), name)
	r.removeTime = time.Now().Add(2 * time.Hour)
	r.t = time.NewTimer(r.removeTime.Sub(time.Now()))
	_, err := r.pubsub.Receive(context.Background())
	if err != nil {
		panic(err)
	}
	go func() {
		ch := r.pubsub.Channel()
		for {
			var m *redis.Message
			select {
			case m = <-ch:
			case <-r.t.C:
				c.ExitRoom(name)
				return
			}
			if Test {
				seam <- struct{}{}
			}
			r.lock.Lock()
			r.history = append(r.history, m.Payload)
			r.removeTime = time.Now().Add(2 * time.Hour)
			r.t.Reset(r.removeTime.Sub(time.Now()))
			r.lock.Unlock()
		}
	}()
	c.all.Store(name, r)
	return true
}

func (c *pubsub_channel) GetInfo(roomname string) (history []string, ttl time.Duration, exist bool) {
	v, ok := c.all.Load(roomname)
	if !ok {
		return nil, 0, false
	}
	r := v.(*pubsub_room)
	r.lock.Lock()
	defer r.lock.Unlock()
	return r.history, r.removeTime.Sub(time.Now()), true
}

func (c *pubsub_channel) SendMessage(roomname string, message string) bool {
	_, ok := c.all.Load(roomname)
	if !ok {
		return false
	}
	for i := range 10 {
		err := c.rdb.Publish(context.Background(), roomname, message).Err()
		if i == 9 {
			panic(err)
		}
		if err != nil {
			continue
		}
		break
	}
	return true
}

func (c *pubsub_channel) waitMessage() {
	<-seam
}

func (c *pubsub_channel) ExitRoom(roomname string) {
	v, ok := c.all.LoadAndDelete(roomname)
	if !ok {
		return
	}
	r := v.(*pubsub_room)
	r.pubsub.Close()
	r.pubsub.Unsubscribe(context.Background(), roomname)
}

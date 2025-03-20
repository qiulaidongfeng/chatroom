package channel

import (
	"context"
	"crypto/tls"
	"errors"
	"time"

	"gitee.com/qiulaidongfeng/chatroom/go/chatroom/internal/config"
	"github.com/redis/go-redis/v9"
)

var _ Channel = (*list_channel)(nil)

// list_channel 管理基于redis列表实现的聊天室
// TODO:支持多个用户进入并退出聊天室
// TODO:支持返回在线人数
type list_channel struct {
	rdb *redis.Client
}

func (c *list_channel) Init() {
	c.rdb = redis.NewClient(&redis.Options{Addr: "127.0.0.1:6379", DB: 15, Password: config.GetRedisPassword(), TLSConfig: &tls.Config{MinVersion: tls.VersionTLS13}})
	if err := c.rdb.Ping(context.Background()).Err(); err != nil {
		panic(err)
	}
	if config.Test {
		if len(c.rdb.Keys(context.Background(), "*").Val()) != 0 {
			panic("测试应该使用空数据库")
		}
	}
}

var exists = errors.New("exists")

func (c *list_channel) CreateRoom(name string) (ret bool) {
	for i := 0; i < 10; i++ {
		err := c.rdb.Watch(context.Background(), func(tx *redis.Tx) error {
			//如果聊天室已经存在
			if tx.Exists(context.Background(), name).Val() == 1 {
				return exists
			}
			_, err := tx.TxPipelined(context.Background(), func(pipe redis.Pipeliner) error {
				err := pipe.LPush(context.Background(), name, "").Err()
				if err != nil {
					return err
				}
				return pipe.Expire(context.Background(), name, 2*time.Hour).Err()
			})
			return err
		}, name)
		if err == exists {
			return false
		}
		if err == redis.TxFailedErr {
			i--
			continue
		}
		if i == 9 {
			panic(err)
		}
		if err == nil {
			break
		}
	}
	return true
}

func (c *list_channel) SendMessage(roomname string, message string) bool {
	val := c.rdb.RPushX(context.Background(), roomname, message).Val()
	if val == 0 { //可能在这台服务器执行到了这里，在聊天室刚好自动删除，或聊天室本身就不存在
		return false
	}
	//TODO:处理如果在这里出现错误
	if !c.rdb.Expire(context.Background(), roomname, 2*time.Hour).Val() {
		return false
	}
	return true
}

func (c *list_channel) GetInfo(roomname string) (history []string, ttl time.Duration, exist bool) {
	//TODO:处理下面三条语句出现错误
	l := c.rdb.LLen(context.Background(), roomname).Val()
	history = c.rdb.LRange(context.Background(), roomname, 0, l).Val()
	if len(history) == 0 {
		return nil, 0, false
	}
	ttl = c.rdb.TTL(context.Background(), roomname).Val()
	return history[1:], ttl, true
}

func (c *list_channel) ExitRoom(roomname string) {
	//TODO:处理如果在这里出现错误
	c.rdb.Del(context.Background(), roomname)
	if config.Test {
		c.rdb.FlushAll(context.Background())
	}
}

// waitMessage 实现接口存在
func (c *list_channel) waitMessage() {
}

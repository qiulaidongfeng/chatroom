package channel

import (
	"slices"
	"testing"
	"time"

	"gitee.com/qiulaidongfeng/chatroom/go/chatroom/internal/config"
)

func TestAll(t *testing.T) {
	C.CreateRoom("test")
	id := C.EnterRoom("test", 10*time.Second)
	_ = C.EnterRoom("test", 10*time.Second)
	C.SendMessage("test", "k")
	C.waitMessage()
	C.SendMessage("test", "k")
	C.waitMessage()
	got, _, _, online := C.GetInfo("test", id)
	if online != 2 {
		t.Fatalf("got %d, want 2", online)
	}
	C.ExitRoom("test")
	want := []string{"k", "k"}
	if !slices.Equal(got, want) {
		t.Fatalf("got %v, want %v", got, want)
	}
}

func init() {
	if !config.Test {
		panic("测试应该设置TEST环境变量为非空")
	}
}

package channel

import (
	"slices"
	"testing"
)

func TestAll(t *testing.T) {
	C.CreateRoom("test")
	C.SendMessage("test", "k")
	C.waitMessage()
	C.SendMessage("test", "k")
	C.waitMessage()
	got, _, _ := C.GetInfo("test")
	C.ExitRoom("test")
	want := []string{"k", "k"}
	if !slices.Equal(got, want) {
		t.Fatalf("got %v, want %v", got, want)
	}
}

func init() {
	if !Test {
		panic("测试应该设置TEST环境变量为非空")
	}
}

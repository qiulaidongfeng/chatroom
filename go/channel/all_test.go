package channel

import (
	"slices"
	"testing"
)

func TestAll(t *testing.T) {
	CreateRoom("test")
	SendMessage("test", "k")
	waitMessage()
	SendMessage("test", "k")
	waitMessage()
	ExitRoom("test")
	got := GetHistory("test")
	want := []string{"k", "k"}
	if !slices.Equal(got, want) {
		t.Fatalf("got %v, want %v", got, want)
	}
}

func init() {
	test = true
}

package main

import (
	"gitee.com/qiulaidongfeng/chatroom/go/chatroom"
)

func main() {
	chatroom.S.RunTLS(":4431", "./cert.pem", "./key.pem")
}

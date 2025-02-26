package main

import (
	"bytes"
	"chatroom/channel"
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"strings"
	"unsafe"

	"github.com/gin-gonic/gin"
)

var static = filepath.Join("."+string(filepath.Separator), "static")
var html = filepath.Join(static, "html")
var index = filepath.Join(html, "index.html")
var createroom = filepath.Join(html, "createroom.html")
var enterroom = filepath.Join(html, "enterroom.html")
var tmpl = filepath.Join(static, "template")

func main() {
	s := gin.Default()
	s.GET("/", func(ctx *gin.Context) {
		ctx.File(index)
	})
	s.GET("/createroom", func(ctx *gin.Context) {
		ctx.File(createroom)
	})
	s.POST("/createroom", func(ctx *gin.Context) {
		name := ctx.PostForm("roomName")
		if name == "" {
			ctx.String(400, "聊天室名不能为空")
			return
		}
		//TODO:在聊天室已创建时报错
		channel.CreateRoom(name)
		enterRoom(ctx, name)
	})
	s.GET("/enterroom", func(ctx *gin.Context) {
		if name := ctx.Query("roomname"); name != "" {
			enterRoom(ctx, name)
			return
		}
		ctx.File(enterroom)
	})
	s.POST("/enterroom", func(ctx *gin.Context) {
		name := ctx.PostForm("roomName")
		if name == "" {
			ctx.String(400, "聊天室名不能为空")
			return
		}
		//TODO:在聊天室不存在时报错
		enterRoom(ctx, name)
	})
	s.POST("/sendMessage", func(ctx *gin.Context) {
		name := ctx.Query("roomname")
		msg := ctx.PostForm("message")
		if msg == "" {
			ctx.String(400, "消息不能为空")
			return
		}
		//Note:发送消息到不存在的聊天室时报错
		if !channel.SendMessage(name, msg) {
			ctx.String(400, "不能发送消息到不存在的聊天室")
			return
		}

		ret := `
		<!DOCTYPE html>
			<head>
				<meta charset="UTF-8">
			</head>
			<body>
			</body>
			<script>
				function f() {
					window.location.href = "%s";
				}
				f();
			</script>
		</html>
		`
		ret = fmt.Sprintf(ret, strings.Join([]string{"http://", ctx.Request.Host, "/enterroom?roomname=", name}, ""))
		ctx.Data(200, "text/html", unsafe.Slice(unsafe.StringData(ret), len(ret)))
	})
	//TODO:使用https
	s.Run(":801")
}

func enterRoom(ctx *gin.Context, name string) {
	var buf bytes.Buffer
	err := roomtmpl.Execute(&buf, map[string]string{"roomname": name})
	if err != nil {
		panic(err)
	}
	ctx.Data(200, "text/html", buf.Bytes())
}

var roomtmpl = func() *template.Template {
	t := template.New("room")
	t.Funcs(template.FuncMap{"getAllMsg": func(roomname string) []string {
		return channel.GetHistory(roomname)
	}})
	file, err := os.ReadFile(filepath.Join(tmpl, "room.temp"))
	if err != nil {
		panic(err)
	}
	t, err = t.Parse(string(file))
	if err != nil {
		panic(err)
	}
	return t
}()

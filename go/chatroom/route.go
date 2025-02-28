package chatroom

import (
	"bytes"
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"strings"
	"time"
	"unsafe"

	"gitee.com/qiulaidongfeng/chatroom/go/chatroom/channel"
	"github.com/gin-gonic/gin"
)

func Handle(s *gin.Engine) {
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
		if !channel.CreateRoom(name) {
			ctx.String(400, "聊天室 %s 已创建", name)
			return
		}
		redirect(ctx, name)
	})
	s.GET("/enterroom", func(ctx *gin.Context) {
		if name := ctx.Query("roomname"); name != "" {
			if !enterRoom(ctx, name) {
				ctx.String(400, "聊天室 %s 不存在", name)
				return
			}
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
		if !enterRoom(ctx, name) {
			ctx.String(400, "聊天室 %s 不存在", name)
			return
		}
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
		redirect(ctx, name)
	})
}

func enterRoom(ctx *gin.Context, name string) bool {
	var buf bytes.Buffer
	h, r, exist := channel.GetInfo(name)
	if !exist {
		return false
	}
	err := roomtmpl.Execute(&buf, map[string]any{"roomname": name, "history": h, "removetime": r.Sub(time.Now()).String()})
	if err != nil {
		panic(err)
	}
	ctx.Data(200, "text/html", buf.Bytes())
	return true
}

var roomtmpl = func() *template.Template {
	t := template.New("room")
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

func redirect(ctx *gin.Context, name string) {
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
	ret = fmt.Sprintf(ret, strings.Join([]string{"https://", ctx.Request.Host, "/enterroom?roomname=", name}, ""))
	ctx.Data(200, "text/html", unsafe.Slice(unsafe.StringData(ret), len(ret)))
}

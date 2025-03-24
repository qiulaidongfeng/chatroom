package chatroom

import (
	"bytes"
	"fmt"
	"log/slog"
	"strconv"
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
		if !channel.C.CreateRoom(name) {
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
	s.POST("/sendMessage", func(ctx *gin.Context) {
		name := ctx.Query("roomname")
		msg := ctx.PostForm("message")
		if msg == "" {
			ctx.String(400, "消息不能为空")
			return
		}
		//Note:发送消息到不存在的聊天室时报错
		if !channel.C.SendMessage(name, msg) {
			ctx.String(400, "不能发送消息到不存在的聊天室")
			return
		}
		redirect(ctx, name)
	})
	s.GET("/exitroom", func(ctx *gin.Context) {
		name := ctx.Query("roomname")
		id, _ := ctx.Cookie(name + "_id")
		channel.C.ExitRoom(name, id)
		redirect(ctx, "/")
	})
}

func enterRoom(ctx *gin.Context, name string) bool {
	id, _ := ctx.Cookie(name + "_id")
	step, _ := ctx.Cookie(name + "_step")
	var id_expire = 10
	if step != "" && step != "10" {
		var err error
		id_expire, err = strconv.Atoi(step)
		if err != nil {
			slog.Error("", "err", err)
			id_expire = 10
		}
	}
	id_expire *= 2 //TODO:设置更好的值
	if id == "" {
		id = channel.C.EnterRoom(name, time.Duration(id_expire)*time.Second)
	} else {
		channel.C.SetIdExpire(name, id, time.Duration(id_expire)*time.Second)
	}
	ctx.SetCookie(name+"_id", id, id_expire, "", "", true, true)
	var buf bytes.Buffer
	h, r, exist, online := channel.C.GetInfo(name, id)
	if !exist {
		return false
	}
	err := roomtmpl.Execute(&buf, map[string]any{"roomname": name, "history": h, "removetime": r.String(), "expire": 2 * 60 * 60, "online": online})
	if err != nil {
		panic(err)
	}
	ctx.Data(200, "text/html", buf.Bytes())
	return true
}

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
	if name == "/" {
		ret = fmt.Sprintf(ret, "https://"+ctx.Request.Host)
	} else {
		ret = fmt.Sprintf(ret, strings.Join([]string{"https://", ctx.Request.Host, "/enterroom?roomname=", name}, ""))
	}
	ctx.Data(200, "text/html", unsafe.Slice(unsafe.StringData(ret), len(ret)))
}

package main

import (
	"path/filepath"

	"github.com/gin-gonic/gin"
)

var static = filepath.Join("."+string(filepath.Separator), "static")
var html = filepath.Join(static, "html")
var index = filepath.Join(html, "index.html")
var createroom = filepath.Join(html, "createroom.html")
var enterroom = filepath.Join(html, "enterroom.html")

func main() {
	s := gin.Default()
	s.GET("/", func(ctx *gin.Context) {
		ctx.File(index)
	})
	s.GET("/createroom", func(ctx *gin.Context) {
		ctx.File(createroom)
	})
	s.GET("/enterroom", func(ctx *gin.Context) {
		ctx.File(enterroom)
	})
	//TODO:使用https
	s.Run(":801")
}

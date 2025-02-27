package chatroom

import (
	"path/filepath"

	"github.com/gin-gonic/gin"
)

var static = filepath.Join("."+string(filepath.Separator), "static")
var html = filepath.Join(static, "html")
var index = filepath.Join(html, "index.html")
var createroom = filepath.Join(html, "createroom.html")
var enterroom = filepath.Join(html, "enterroom.html")
var tmpl = filepath.Join(static, "template")

var S = gin.Default()

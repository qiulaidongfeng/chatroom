package config

import (
	"fmt"
	"os"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

var v *viper.Viper = func() *viper.Viper {
	v := viper.New()
	prefix := ""
	if Test {
		prefix = "../"
	}
	v.SetConfigFile(prefix + "config.ini")
	v.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("Config file changed:", e.Name)
	})
	v.WatchConfig()
	err := v.ReadInConfig()
	if err != nil {
		panic(err)
	}
	return v
}()

var Test bool = os.Getenv("TEST") != ""

func GetMode() int {
	return v.GetInt("chatroom.mode")
}

func GetRedisPassword() string {
	return v.GetString("redis.password")
}

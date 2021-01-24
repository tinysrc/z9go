package conf

import (
	"github.com/spf13/viper"
)

// Global instance
var Global *viper.Viper

func init() {
	Global = viper.New()
	if Global == nil {
		panic("create global config failed")
	}
	Global.SetConfigName("config")
	Global.AddConfigPath("/etc/z9/")
	Global.AddConfigPath("./")
	Global.ReadInConfig()
	Global.SetDefault("service.name", "z9")
	Global.SetDefault("service.addr", ":8080")
	Global.SetDefault("service.gateway", "localhost:80")
}

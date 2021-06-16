package conf

import (
	"flag"

	"github.com/spf13/viper"
)

// Global instance
var Global *viper.Viper

func init() {
	Global = viper.New()
	if Global == nil {
		panic("create global config failed")
	}
	cfg := flag.String("c", "config", "specify config filename")
	flag.Parse()
	Global.SetConfigName(*cfg)
	Global.AddConfigPath("/etc/z9/")
	Global.AddConfigPath("./")
	err := Global.ReadInConfig()
	if err != nil {
		panic(err)
	}
}

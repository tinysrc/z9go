package svc

import "github.com/tinysrc/z9go/pkg/conf"

func initConfig() {
	conf.Global.SetDefault("service.name", "z9")
	conf.Global.SetDefault("service.addr", ":8080")
	conf.Global.SetDefault("service.gateway", "127.0.0.1:80")
}

func init() {
	initConfig()
}

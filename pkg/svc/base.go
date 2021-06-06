package svc

import "github.com/tinysrc/z9go/pkg/conf"

func init() {
	conf.Global.SetDefault("service.name", "z9")
	conf.Global.SetDefault("service.addr", ":8080")
	conf.Global.SetDefault("service.gateway", "127.0.0.1:80")
	conf.Global.SetDefault("service.tls.caFile", "./certs/ca.pem")
	conf.Global.SetDefault("service.tls.client.serverName", "z9os.com")
	conf.Global.SetDefault("service.tls.client.certFile", "./certs/client.pem")
	conf.Global.SetDefault("service.tls.client.keyFile", "./certs/client.key")
	conf.Global.SetDefault("service.tls.server.certFile", "./certs/server.pem")
	conf.Global.SetDefault("service.tls.server.keyFile", "./certs/server.key")
}

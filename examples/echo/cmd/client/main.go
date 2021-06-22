package main

import (
	"context"
	"strconv"

	"github.com/tinysrc/z9go/examples/echo/pb"
	"github.com/tinysrc/z9go/pkg/conf"
	"github.com/tinysrc/z9go/pkg/log"
	"github.com/tinysrc/z9go/pkg/rpc"
	"go.uber.org/zap"
)

func main() {
	defer log.Close()
	sign := conf.Global.GetString("jwt.sign")
	cli := rpc.NewClient(sign)
	conn, err := cli.Dial("echo")
	if err != nil {
		return
	}
	defer conn.Close()
	echoService := pb.NewEchoServiceClient(conn)
	for i := 0; i < 10; i++ {
		out, err := echoService.Echo(context.Background(), &pb.StringMessage{Value: strconv.Itoa(i)})
		if err != nil {
			log.Error("call echo failed", zap.Error(err))
		} else {
			log.Debug("call echo success", zap.Any("result", out))
		}
	}
}

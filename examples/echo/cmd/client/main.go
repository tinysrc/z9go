package main

import (
	"context"
	"strconv"
	"time"

	"github.com/tinysrc/z9go/examples/echo/pb"
	"github.com/tinysrc/z9go/pkg/log"
	"github.com/tinysrc/z9go/pkg/svr"
	"go.uber.org/zap"
)

func main() {
	defer log.Close()
	cli, err := svr.NewClient("echo")
	if err != nil {
		return
	}
	defer cli.Close()
	echoClient := pb.NewEchoClient(cli.Conn)
	for i := 0; i < 10; i++ {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()
		out, err := echoClient.Echo(ctx, &pb.StringMessage{Value: strconv.Itoa(i)}, cli.Opts...)
		if err != nil {
			log.Error("call echo failed", zap.Error(err))
		} else {
			log.Debug("call echo success", zap.Any("result", out))
		}
	}
}

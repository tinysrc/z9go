package main

import (
	"strconv"

	"github.com/tinysrc/z9go/examples/echo/pb"
	"github.com/tinysrc/z9go/pkg/log"
	"github.com/tinysrc/z9go/pkg/svr"
	"go.uber.org/zap"
)

func main() {
	defer log.Close()
	cli := svr.NewClient()
	conn, err := cli.Dial("echo")
	if err != nil {
		return
	}
	defer conn.Close()
	impl := pb.NewEchoClient(conn)
	for i := 0; i < 10; i++ {
		out, err := impl.Echo(cli.NewCallCtx(), &pb.StringMessage{Value: strconv.Itoa(i)})
		if err != nil {
			log.Error("call echo failed", zap.Error(err))
		} else {
			log.Debug("call echo success", zap.Any("result", out))
		}
	}
}

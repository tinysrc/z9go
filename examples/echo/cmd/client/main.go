package main

import (
	"context"
	"strconv"
	"time"

	"github.com/tinysrc/z9go/examples/echo/pb"
	"github.com/tinysrc/z9go/pkg/conf"
	"github.com/tinysrc/z9go/pkg/log"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func main() {
	defer log.Close()
	addr := conf.Global.GetString("service.addr")
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		log.Fatal("grpc dial failed", zap.Error(err))
		return
	}
	defer conn.Close()
	service := pb.NewEchoServiceClient(conn)
	for i := 0; i < 10; i++ {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()
		out, err := service.Echo(ctx, &pb.StringMessage{Value: strconv.Itoa(i)})
		if err != nil {
			log.Error("call echo failed", zap.Error(err))
		} else {
			log.Debug("call echo success", zap.Any("result", out))
		}
	}
}

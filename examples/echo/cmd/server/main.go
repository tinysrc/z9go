package main

import (
	"context"
	"os"
	"os/signal"
	"sync"

	"github.com/tinysrc/z9go/examples/echo/pb"
	"github.com/tinysrc/z9go/pkg/conf"
	"github.com/tinysrc/z9go/pkg/log"
	"github.com/tinysrc/z9go/pkg/rpc"
	"go.uber.org/zap"
)

type service struct {
	pb.UnimplementedEchoServiceServer
}

func (s *service) Echo(ctx context.Context, in *pb.StringMessage) (*pb.StringMessage, error) {
	log.Debug("Echo", zap.Any("in", in))
	return &pb.StringMessage{Value: in.Value}, nil
}

func main() {
	defer log.Close()
	// init rpc server
	svr := rpc.NewServer(nil, nil)
	svc := &service{}
	pb.RegisterEchoServiceServer(svr.Server, svc)
	// run
	wg := sync.WaitGroup{}
	// run rpc server
	wg.Add(1)
	go func() {
		defer wg.Done()
		addr := conf.Global.GetString("service.listen")
		log.Info("rpc server start", zap.String("addr", addr))
		if err := svr.Run(addr); err != nil {
			log.Error("rpc server run failed", zap.Error(err))
		}
		log.Info("rpc server stopped", zap.String("addr", addr))
	}()
	// signal
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	<-signalChan
	svr.Stop()
	wg.Wait()
}

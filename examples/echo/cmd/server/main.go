package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/tinysrc/z9go/examples/echo/pb"
	"github.com/tinysrc/z9go/pkg/conf"
	"github.com/tinysrc/z9go/pkg/log"
	"github.com/tinysrc/z9go/pkg/naming/registry"
	"github.com/tinysrc/z9go/pkg/naming/registry/etcd"
	"github.com/tinysrc/z9go/pkg/rpc"
	etcd3 "go.etcd.io/etcd/client/v3"
	"go.uber.org/zap"
	"google.golang.org/grpc/metadata"
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
	// init etcd
	cfg := etcd3.Config{
		Endpoints: conf.Global.GetStringSlice("etcd.endpoints"),
	}
	si := &registry.ServiceInfo{
		Id:       conf.Global.GetString("service.id"),
		Name:     conf.Global.GetString("service.name"),
		Version:  conf.Global.GetString("service.version"),
		Address:  conf.Global.GetString("service.addr"),
		Metadata: metadata.Pairs("weight", conf.Global.GetString("service.weight")),
	}
	reg, err := etcd.NewRegistrar(&etcd.Config{
		Config: cfg,
		Dir:    "/z9/backend/services",
		TTL:    10 * time.Second,
	})
	if err != nil {
		log.Fatal("new registrar failed", zap.Error(err))
		return
	}
	log.Info("init etcd success")
	// init rpc server
	svr := rpc.NewServer(nil)
	svc := &service{}
	pb.RegisterEchoServiceServer(svr.Server, svc)
	// init rpc gateway
	ctx := context.Background()
	mux := runtime.NewServeMux()
	gwsvr := http.Server{
		Addr:    conf.Global.GetString("service.gwlisten"),
		Handler: mux,
	}
	err = pb.RegisterEchoServiceHandlerServer(ctx, mux, svc)
	if err != nil {
		log.Fatal("register gateway protocol failed", zap.Error(err))
		return
	}
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
	// run rpc gateway
	wg.Add(1)
	go func() {
		defer wg.Done()
		log.Info("rpc gateway start", zap.String("addr", gwsvr.Addr))
		if err := gwsvr.ListenAndServe(); err != nil {
			log.Error("rpc gateway run failed", zap.Error(err))
		}
		log.Info("rpc gateway stopped", zap.String("addr", gwsvr.Addr))
	}()
	// run etcd registrar
	wg.Add(1)
	go func() {
		defer wg.Done()
		log.Info("etcd register start")
		reg.Register(si)
		log.Info("etcd register stopped")
	}()
	// signal
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	<-signalChan
	// etcd unregister
	reg.Unregister(si)
	gwsvr.Shutdown(ctx)
	svr.Stop()
	wg.Wait()
}

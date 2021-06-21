package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
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
		Metadata: metadata.Pairs("weight", conf.Global.GetString("service.weight")),
	}
	reg, err := etcd.NewRegistrar(&etcd.Config{
		Config: cfg,
		Dir:    "/z9/backend/services",
		TTL:    10 * time.Second,
	})
	if err != nil {
		panic(err)
	}
	// init rpc server
	svr := rpc.NewServer(nil)
	svc := &service{}
	pb.RegisterEchoServiceServer(svr.Server, svc)
	// init rpc gateway
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	mux := runtime.NewServeMux()
	pb.RegisterEchoServiceHandlerServer(ctx, mux, svc)
	// run
	wg := sync.WaitGroup{}
	// run rpc server
	wg.Add(1)
	go func() {
		defer wg.Done()
		addr := conf.Global.GetString("service.addr")
		svr.Run(addr)
	}()
	// run rpc gateway
	go func() {
		defer wg.Done()
		addr := conf.Global.GetString("service.gwaddr")
		http.ListenAndServe(addr, mux)
	}()
	// run etcd registrar
	wg.Add(1)
	go func() {
		defer wg.Done()
		reg.Register(si)
	}()
	// signal
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGTERM, syscall.SIGINT)
	<-signalChan
	// etcd unregister
	reg.Unregister(si)
	cancel()
	svr.Stop()
	wg.Wait()
}

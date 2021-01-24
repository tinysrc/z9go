package main

import (
	"context"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/tinysrc/z9go/examples/echo/pb"
	"github.com/tinysrc/z9go/pkg/conf"
	"github.com/tinysrc/z9go/pkg/log"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func main() {
	defer log.Close()
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	mux := runtime.NewServeMux()
	dest := conf.Global.GetString("service.dest")
	opts := []grpc.DialOption{grpc.WithInsecure()}
	err := pb.RegisterEchoHandlerFromEndpoint(ctx, mux, dest, opts)
	if err != nil {
		log.Fatal("register grpc failed", zap.Error(err))
		return
	}
	addr := conf.Global.GetString("service.addr")
	if err = http.ListenAndServe(addr, mux); err != nil {
		log.Fatal("http listen and serve failed", zap.Error(err))
	}
}

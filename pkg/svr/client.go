package svr

import (
	"fmt"

	"github.com/tinysrc/z9go/pkg/conf"
	"github.com/tinysrc/z9go/pkg/log"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

// Dial impl
func Dial(serviceName string) (conn *grpc.ClientConn, err error) {
	gateway := conf.Global.GetString("service.gateway")
	addr := fmt.Sprintf("%s/%s.grpc", gateway, serviceName)
	conn, err = grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		log.Fatal("grpc dial failed", zap.Error(err), zap.String("addr", addr))
	} else {
		log.Info("grpc dial success", zap.String("addr", addr))
	}
	return
}

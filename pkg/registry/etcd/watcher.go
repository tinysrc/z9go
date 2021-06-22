package etcd

import (
	"encoding/json"
	"sync"

	"github.com/tinysrc/z9go/pkg/registry"
	"go.etcd.io/etcd/api/v3/mvccpb"
	etcd3 "go.etcd.io/etcd/client/v3"
	"golang.org/x/net/context"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/resolver"
)

type Watcher struct {
	key    string
	cli    *etcd3.Client
	ctx    context.Context
	cancel context.CancelFunc
	wg     sync.WaitGroup
	addrs  []resolver.Address
}

func newWatcher(key string, cli *etcd3.Client) *Watcher {
	ctx, cancel := context.WithCancel(context.Background())
	return &Watcher{
		key:    key,
		cli:    cli,
		ctx:    ctx,
		cancel: cancel,
	}
}

func (w *Watcher) Close() {
	w.cancel()
}

func (w *Watcher) GetAllAddresses() []resolver.Address {
	ret := []resolver.Address{}
	resp, err := w.cli.Get(w.ctx, w.key, etcd3.WithPrefix())
	if err == nil {
		addrs := extractAddrs(resp)
		if len(addrs) > 0 {
			for _, v := range addrs {
				ret = append(ret, resolver.Address{
					Addr:     v.Address,
					Metadata: &v.Metadata,
				})
			}
		}
	}
	return ret
}

func extractAddrs(resp *etcd3.GetResponse) []registry.ServiceInfo {
	ret := []registry.ServiceInfo{}
	if resp == nil || resp.Kvs == nil {
		return ret
	}
	for i := range resp.Kvs {
		if v := resp.Kvs[i].Value; v != nil {
			si := registry.ServiceInfo{}
			err := json.Unmarshal(v, &si)
			if err != nil {
				grpclog.Infof("parse service info failed error=%s", err.Error())
				continue
			}
			ret = append(ret, si)
		}
	}
	return ret
}

func (w *Watcher) Watch() chan []resolver.Address {
	out := make(chan []resolver.Address, 10)
	w.wg.Add(1)
	go func() {
		defer func() {
			close(out)
			w.wg.Done()
		}()
		wc := w.cli.Watch(w.ctx, w.key, etcd3.WithPrefix())
		for resp := range wc {
			for _, ev := range resp.Events {
				switch ev.Type {
				case mvccpb.PUT:
					si := registry.ServiceInfo{}
					err := json.Unmarshal([]byte(ev.Kv.Value), &si)
					if err != nil {
						grpclog.Errorf("parse service info failed error=%s", err.Error())
						continue
					}
					addr := resolver.Address{
						Addr:     si.Address,
						Metadata: &si.Metadata,
					}
					if w.addAddr(addr) {
						out <- w.cloneAddrs(w.addrs)
					}
				case mvccpb.DELETE:
					si := registry.ServiceInfo{}
					err := json.Unmarshal([]byte(ev.Kv.Value), &si)
					if err != nil {
						grpclog.Errorf("pase service info failed error=%s", err.Error())
						continue
					}
					addr := resolver.Address{
						Addr:     si.Address,
						Metadata: &si.Metadata,
					}
					if w.removeAddr(addr) {
						out <- w.cloneAddrs(w.addrs)
					}
				}
			}
		}
	}()
	return out
}

func (w *Watcher) cloneAddrs(in []resolver.Address) []resolver.Address {
	out := make([]resolver.Address, len(in))
	for i := 0; i < len(in); i++ {
		out[i] = in[i]
	}
	return out
}

func (w *Watcher) addAddr(addr resolver.Address) bool {
	for _, v := range w.addrs {
		if addr.Addr == v.Addr {
			return false
		}
	}
	w.addrs = append(w.addrs, addr)
	return true
}

func (w *Watcher) removeAddr(addr resolver.Address) bool {
	for i, v := range w.addrs {
		if addr.Addr == v.Addr {
			w.addrs = append(w.addrs[:i], w.addrs[i+1:]...)
			return true
		}
	}
	return false
}

package etcd

import (
	"fmt"
	"sync"

	etcd3 "go.etcd.io/etcd/client/v3"
	"google.golang.org/grpc/resolver"
)

type etcdResolver struct {
	scheme  string
	config  etcd3.Config
	path    string
	watcher *Watcher
	cc      resolver.ClientConn
	wg      sync.WaitGroup
}

func (r *etcdResolver) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {
	cli, err := etcd3.New(r.config)
	if err != nil {
		return nil, err
	}
	r.cc = cc
	r.watcher = newWatcher(r.path, cli)
	func() {
		r.wg.Add(1)
		go func() {
			defer r.wg.Done()
			out := r.watcher.Watch()
			for addrs := range out {
				r.cc.UpdateState(resolver.State{Addresses: addrs})
			}
		}()
	}()
	return r, nil
}

func (r *etcdResolver) Scheme() string {
	return r.scheme
}

func (r *etcdResolver) ResolveNow(o resolver.ResolveNowOptions) {
}

func (r *etcdResolver) Close() {
	r.watcher.Close()
	r.wg.Wait()
}

func RegisterResolver(scheme string, config etcd3.Config, dir, name, version string) {
	resolver.Register(&etcdResolver{
		scheme: scheme,
		config: config,
		path:   fmt.Sprintf("%s/%s/%s", dir, name, version),
	})
}

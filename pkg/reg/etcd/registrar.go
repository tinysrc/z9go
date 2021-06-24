package etcd

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/tinysrc/z9go/pkg/reg"
	etcd3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/etcdserver/api/v3rpc/rpctypes"
	"golang.org/x/net/context"
	"google.golang.org/grpc/grpclog"
)

type Registrar struct {
	sync.RWMutex
	config  *Config
	client  *etcd3.Client
	cancels map[string]context.CancelFunc
}

type Config struct {
	Config etcd3.Config
	Dir    string
	TTL    time.Duration
}

func NewRegistrar(cfg *Config) (*Registrar, error) {
	cli, err := etcd3.New(cfg.Config)
	if err != nil {
		return nil, err
	}
	return &Registrar{
		config:  cfg,
		client:  cli,
		cancels: make(map[string]context.CancelFunc),
	}, nil
}

func (r *Registrar) Register(si *reg.ServiceInfo) error {
	val, err := json.Marshal(si)
	if err != nil {
		return err
	}
	key := fmt.Sprintf("%s/%s/%s/%s", r.config.Dir, si.Name, si.Version, si.Id)
	value := string(val)
	ctx, cancel := context.WithCancel(context.Background())
	r.Lock()
	r.cancels[si.Id] = cancel
	r.Unlock()
	insertFunc := func() error {
		resp, err := r.client.Grant(ctx, int64(r.config.TTL/time.Second))
		if err != nil {
			return err
		}
		_, err = r.client.Get(ctx, key)
		if err != nil {
			if err == rpctypes.ErrKeyNotFound {
				_, err := r.client.Put(ctx, key, value, etcd3.WithLease(resp.ID))
				if err != nil {
					grpclog.Infof("etcd put failed key=%s error=%s", key, err.Error())
				}
			} else {
				grpclog.Infof("etcd get failed key=%s error=%s", key, err.Error())
			}
			return err
		} else {
			_, err = r.client.Put(ctx, key, value, etcd3.WithLease(resp.ID))
			if err != nil {
				grpclog.Infof("etcd refresh failed key=%s error=%s", key, err.Error())
				return err
			}
		}
		return nil
	}
	err = insertFunc()
	if err != nil {
		return err
	}
	ticker := time.NewTicker(r.config.TTL / 5)
	for {
		select {
		case <-ticker.C:
			insertFunc()
		case <-ctx.Done():
			ticker.Stop()
			_, err := r.client.Delete(context.Background(), key)
			if err != nil {
				grpclog.Infof("etcd delete failed key=%s error=%s", key, err.Error())
			}
			return nil
		}
	}
}

func (r *Registrar) Unregister(si *reg.ServiceInfo) error {
	r.RLock()
	cancel, ok := r.cancels[si.Id]
	r.RUnlock()
	if ok {
		cancel()
	}
	return nil
}

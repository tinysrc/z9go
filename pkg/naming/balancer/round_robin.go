package balancer

import (
	"math/rand"
	"sync"

	"google.golang.org/grpc/balancer"
	"google.golang.org/grpc/balancer/base"
	"google.golang.org/grpc/grpclog"
)

const RoundRobin = "round_robin"

func init() {
	balancer.Register(newRoundRobinBuilder())
}

func newRoundRobinBuilder() balancer.Builder {
	return base.NewBalancerBuilder(RoundRobin, &roundRobinPickerBuilder{}, base.Config{HealthCheck: true})
}

type roundRobinPickerBuilder struct{}

func (*roundRobinPickerBuilder) Build(buildInfo base.PickerBuildInfo) balancer.Picker {
	grpclog.Infof("roundRobinPickerBuilder build info=%v", buildInfo)
	if len(buildInfo.ReadySCs) == 0 {
		return base.NewErrPicker(balancer.ErrNoSubConnAvailable)
	}
	var scs []balancer.SubConn
	for subConn, subConnInfo := range buildInfo.ReadySCs {
		weight := getWeight(subConnInfo.Address)
		for i := 0; i < weight; i++ {
			scs = append(scs, subConn)
		}
	}
	return &roundRobinPicker{
		subConns: scs,
		next:     rand.Intn(len(scs)),
	}
}

type roundRobinPicker struct {
	subConns []balancer.SubConn
	mu       sync.Mutex
	next     int
}

func (p *roundRobinPicker) Pick(balancer.PickInfo) (balancer.PickResult, error) {
	ret := balancer.PickResult{}
	p.mu.Lock()
	ret.SubConn = p.subConns[p.next]
	p.next = (p.next + 1) % len(p.subConns)
	p.mu.Unlock()
	return ret, nil
}

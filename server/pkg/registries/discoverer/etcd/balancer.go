package etcd

import (
	"github.com/pkg/errors"
	"go.luxshare-ict.com/pkg/apeol.com/sort"
	"go.luxshare-ict.com/pkg/config"
	"google.golang.org/grpc/balancer"
	"google.golang.org/grpc/balancer/base"
	"sync"
)

type Picker struct {
	name    string
	config  *config.Service
	pick    int
	total   int
	subConn []balancer.SubConn
	bridge  []int
	mu      sync.Mutex
}

func NewBuilder(cfg *config.Service) *Picker {
	bp := &Picker{
		config: cfg,
		name:   cfg.Name,
	}
	balancer.Register(base.NewBalancerBuilderV2(cfg.Name, bp, base.Config{HealthCheck: false}))
	return bp
}

func (t *Picker) Build(info base.PickerBuildInfo) balancer.V2Picker {
	t.mu.Lock()
	defer t.mu.Unlock()
	l := len(info.ReadySCs)
	if l == 0 {
		return base.NewErrPickerV2(balancer.ErrNoSubConnAvailable)
	}

	t.total = 0
	t.pick = 0
	cfg := t.config
	t.subConn = make([]balancer.SubConn, l)
	ws := make([]float64, l)
	i := 0
	var total float64
	for sc, sci := range info.ReadySCs {
		node := GetAddrInfo(sci.Address)
		w := float64(node.Weight)
		ws[i] = w
		total += w
		t.subConn[i] = sc
		i++

		if cfg.Version == "" {
			// 使用最新版本

		} else if node.Version == cfg.Version {

		}
	}
	list := &sort.Store{}
	for m, w := range ws {
		j := total / w
		var k float64 = 0
		for n := w; n > 0; n-- {
			k++
			list.Append(&bridge{
				w: j * k,
				t: m,
			})
		}
	}
	t.bridge = make([]int, 0)
	list.Sort().Foreach(func(key int, val sort.Sort) {
		sc := val.(*bridge)
		t.bridge = append(t.bridge, sc.t)
	})
	t.total = int(total)
	return t
}
func (t *Picker) Pick(balancer.PickInfo) (res balancer.PickResult, err error) {
	t.mu.Lock()
	defer t.mu.Unlock()

	if t.total <= 0 {
		err = errors.Errorf("无可用服务")
		return
	}
	if t.pick >= t.total {
		t.pick = 0
	}
	res = balancer.PickResult{SubConn: t.subConn[t.bridge[t.pick]]}
	t.pick++
	return
}

type bridge struct {
	w float64
	t int
}

func (s *bridge) Val() (v sort.Value) {
	return sort.Value{Float64: s.w}
}

func (s *bridge) Less(l sort.Sort) bool {
	return s.w-l.Val().Float64 > 0
}

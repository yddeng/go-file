package etcd

import (
	"github.com/pkg/errors"
	"go.etcd.io/etcd/clientv3"
	"go.luxshare-ict.com/pkg/config"
	"go.luxshare-ict.com/pkg/registries/discoverer"
	"go.luxshare-ict.com/pkg/registries/tools"
	"google.golang.org/grpc/attributes"
	"google.golang.org/grpc/resolver"
	"sync"
)

type Discover struct {
	name        string
	registry    string
	client      *clientv3.Client
	conn        resolver.ClientConn
	serviceList sync.Map
	schema      string
	prefix      string
	total       int
	index       int
	cache       []resolver.Address
	config      *config.Service
	watcher     discoverer.Watcher
}

func NewDiscover(watcher discoverer.Watcher, cli *clientv3.Client, cfg *config.Service, prefix string) discoverer.Discoverer {
	NewBuilder(cfg)
	d := &Discover{
		name:   cfg.Name,
		client: cli,
		//serviceList: sync.Map{},
		config:  cfg,
		prefix:  prefix,
		schema:  cfg.Schema,
		watcher: watcher,
	}
	resolver.Register(d)
	return d
}

func ToEtcdDiscoverer(m discoverer.Discoverer) *Discover {
	return m.(*Discover)
}

type attributeKey struct{}

//Build 为给定目标创建一个新的`resolver`，当调用`grpc.Dial()`时执行
func (m *Discover) Build(target resolver.Target, conn resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {
	m.conn = conn
	//registry.Get(m.registry).Watch(m)
	m.watcher.Watch(m)
	return m, nil
}

func (m *Discover) Name() string {
	return m.name
}
func (m *Discover) GetConfig() *config.Service {
	return m.config
}
func (m *Discover) Scheme() string {
	return m.schema
}

func (m *Discover) Export() *discoverer.Export {
	return &discoverer.Export{Etcd: m.client}
}

func (m *Discover) Prefix() string {
	return m.prefix
}

func (m *Discover) Client() *clientv3.Client {
	return m.client
}

func (m *Discover) Conn() resolver.ClientConn {
	return m.conn
}

func (m *Discover) UpdateState() {
	ss := m.GetServices()
	m.conn.UpdateState(resolver.State{Addresses: ss})
}

// ResolveNow 监视目标更新
func (m *Discover) ResolveNow(rn resolver.ResolveNowOptions) {
	//logger.Infof("ResolveNow: %+v, %+v", rn, m)
}

//Close 关闭
func (m *Discover) Close() {
	m.client.Close()
}
func (m *Discover) SetServiceList(key string, val interface{}) (err error) {
	addr, ok := val.(*tools.Addr)
	if !ok {
		if bt, ok := val.([]byte); ok {
			addr, _ = tools.DecodeAddr(bt)
		}
	}
	if addr == nil {
		return errors.Errorf("val is invalid")
	}
	m.serviceList.Store(key, addr)
	m.cache = nil
	//m.UpdateState()
	return nil
}
func (m *Discover) DelServiceList(key string) {
	m.serviceList.Delete(key)
	m.cache = nil
	//m.UpdateState()
}
func GetAddrInfo(addr resolver.Address) *tools.Addr {
	v := addr.Attributes.Value(attributeKey{})
	ai, _ := v.(*tools.Addr)
	return ai
}

//GetServices 获取服务地址
func (m *Discover) GetServices() []resolver.Address {
	if m.cache == nil {
		addrs := make([]resolver.Address, 0)
		m.serviceList.Range(func(k, v interface{}) bool {
			if addr, ok := v.(*tools.Addr); ok {
				addrs = append(addrs, resolver.Address{Addr: addr.Addr, Attributes: attributes.New(attributeKey{}, addr)})
			}
			return true
		})
		m.cache = addrs
	}
	return m.cache
}

package registry

import (
	"context"
	"go.etcd.io/etcd/clientv3"
	"go.etcd.io/etcd/mvcc/mvccpb"
	"go.luxshare-ict.com/pkg/config"
	"go.luxshare-ict.com/pkg/logger"
	"go.luxshare-ict.com/pkg/registries/discoverer"
	discover "go.luxshare-ict.com/pkg/registries/discoverer/etcd"
	"go.luxshare-ict.com/pkg/registries/tools"
	"sync"
	"time"
)

type etcd struct {
	name          string
	active        bool
	config        *config.Registry
	client        *clientv3.Client
	lease         clientv3.Lease
	leaseResp     *clientv3.LeaseGrantResponse
	canclefunc    func()
	keepAliveChan <-chan *clientv3.LeaseKeepAliveResponse
	services      sync.Map
}

type cluster struct {
	list []string
}

func init() {
	new(&etcd{name: "etcd"})
}

func (p *etcd) Name() string {
	return p.name
}

func (p *etcd) New(cfg *config.Registry) (err error) {
	p.config = cfg
	if cfg.DialTimeout <= 0 {
		cfg.DialTimeout = time.Second * 5
	}
	p.client, err = clientv3.New(clientv3.Config{
		Endpoints:   cfg.Cluster,
		DialTimeout: cfg.DialTimeout,
	})
	if err != nil {
		return err
	}
	p.active = true
	if err = p.setLease(cfg.TTL); err != nil {
		logger.Errorf("\n\netcd registry set lease failed: %+v\n\n", err)
		return nil
	}
	go p.listenLeaseRespChan()
	return nil
}

func (p *etcd) Close() {
	if p != nil && p.client != nil {
		p.client.Close()
	}
}
func (p *etcd) IsActive() bool {
	return p.active
}
func (p *etcd) kvco() (kv clientv3.KV, ctx context.Context, opt []clientv3.OpOption) {
	kv = clientv3.NewKV(p.client)
	ctx = context.TODO()
	if p.leaseResp != nil {
		opt = append(opt, clientv3.WithLease(p.leaseResp.ID))
	}
	return
}
func (p *etcd) Register(k string, addr *tools.Addr) (err error) {
	kv, ctx, opt := p.kvco()
	_, err = kv.Put(ctx, k, addr.ToString(), opt...)
	return
}

func (p *etcd) NewDiscover(cfg *config.Service) (d discoverer.Discoverer, err error) {
	prefix := GenPrefix(cfg.Region, cfg.Name, cfg.Type)
	if cfg.Schema == "" {
		cfg.Schema = cfg.Name
	}
	if p.config.DialTimeout <= 0 {
		p.config.DialTimeout = time.Second * 3
	}
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   p.config.Cluster,
		DialTimeout: p.config.DialTimeout,
	})
	if err != nil {
		return nil, err
	}
	if len(cfg.Cluster) > 0 {
		d = discover.NewDiscover(&cluster{list: cfg.Cluster}, cli, cfg, prefix)
	} else {
		d = discover.NewDiscover(p, cli, cfg, prefix)
	}
	return
}

func (p *cluster) Watch(d discoverer.Discoverer) error {
	m := discover.ToEtcdDiscoverer(d)
	for _, item := range p.list {
		m.SetServiceList(item, &tools.Addr{
			Weight: 1,
			Addr:   item,
		})
	}
	m.UpdateState()
	return nil
}
func (p *etcd) Watch(d discoverer.Discoverer) (err error) {
	m := discover.ToEtcdDiscoverer(d)
	opt := clientv3.WithPrefix()
	cli := m.Client()
	resp, err := cli.Get(context.Background(), m.Prefix(), opt)
	if err != nil {
		return err
	}
	for _, ev := range resp.Kvs {
		m.SetServiceList(string(ev.Key), ev.Value)
	}
	m.UpdateState()

	//监视前缀，修改变更的server
	go func() {
		rch := cli.Watch(context.Background(), m.Prefix(), opt)
		for wc := range rch {
			for _, ev := range wc.Events {
				switch ev.Type {
				case mvccpb.PUT:
					m.SetServiceList(string(ev.Kv.Key), ev.Kv.Value)
				case mvccpb.DELETE:
					m.DelServiceList(string(ev.Kv.Key))
				}
			}
			m.UpdateState()
		}
	}()
	return nil
}

func (p *etcd) Deregister(key string) (err error) {
	kv, ctx, opt := p.kvco()
	_, err = kv.Delete(ctx, key, opt...)
	return
}

//设置租约
func (p *etcd) setLease(ttl int64) error {
	lease := clientv3.NewLease(p.client)

	ctx, cancel := context.WithTimeout(context.TODO(), time.Minute)
	if ttl <= 0 {
		ttl = 10
	}
	leaseResp, err := lease.Grant(ctx, ttl)
	if err != nil {
		cancel()
		return err
	}

	ctx, cancelFunc := context.WithCancel(context.TODO())
	leaseRespChan, err := lease.KeepAlive(ctx, leaseResp.ID)

	if err != nil {
		return err
	}

	p.lease = lease
	p.leaseResp = leaseResp
	p.canclefunc = cancelFunc
	p.keepAliveChan = leaseRespChan

	go func() {
		for {
			<-leaseRespChan
			//logger.Print("ttl:", ka.TTL)
		}
	}()
	return nil
}

//撤销租约
func (p *etcd) revokeLease() error {
	logger.Warning("撤销租约")
	p.canclefunc()
	time.Sleep(2 * time.Second)
	_, err := p.lease.Revoke(context.TODO(), p.leaseResp.ID)
	return err
}

//监听续租情况
func (p *etcd) listenLeaseRespChan() {
	/*
		for {
			select {
			case leaseKeepResp := <-p.keepAliveChan:
				if leaseKeepResp == nil {
					// 未关闭也有可能出现这种情况，原因暂未分析
					logger.Warning("已经关闭续租功能")
					return
				} else {
					logger.Info("续租成功")
				}
			}
		}

	*/
}

package discoverer

import (
	"go.etcd.io/etcd/clientv3"
	"go.luxshare-ict.com/pkg/config"
)

type Discoverer interface {
	Name() string
	GetConfig() *config.Service
	Scheme() string
	Prefix() string
	Export() *Export
	Close()
}

type Watcher interface {
	Watch(discoverer Discoverer) error
}

type Export struct {
	Etcd *clientv3.Client
}

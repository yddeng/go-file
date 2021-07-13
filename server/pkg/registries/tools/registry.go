package tools

import (
	"go.luxshare-ict.com/pkg/config"
	"go.luxshare-ict.com/pkg/registries/discoverer"
)

type Registry interface {
	New(cfg *config.Registry) error
	IsActive() bool
	Close()
	Register(key string, addr *Addr) error
	Deregister(key string) error
	NewDiscover(cfg *config.Service) (discoverer.Discoverer, error)
	Watch(d discoverer.Discoverer) error
	Name() string
}

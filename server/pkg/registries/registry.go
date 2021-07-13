package registry

import (
	"github.com/pkg/errors"
	"go.luxshare-ict.com/pkg/config"
	"go.luxshare-ict.com/pkg/registries/tools"
	"strings"
	"sync"
)

var registries = make(map[string]tools.Registry)
var appConfig *config.App
var sep = "/"
var lock sync.Mutex

func new(p tools.Registry) {
	lock.Lock()
	registries[p.Name()] = p
	lock.Unlock()
}
func Init(cfg *config.App) {
	appConfig = cfg
}
func GenPrefix(region, name, _type string) string {
	return strings.Join([]string{region, name, _type, ""}, sep)
}

func Get(name string) tools.Registry {
	lock.Lock()
	defer lock.Unlock()
	return registries[name]
}
func New(name string, cfg *config.Registry) (p tools.Registry, err error) {
	p = Get(name)
	if p == nil {
		err = errors.Errorf("%v 不存在", name)
	} else {
		err = p.New(cfg)
	}
	return
}

func Close() {
	for _, v := range registries {
		if v != nil {
			v.Close()
		}
	}
}

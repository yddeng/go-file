package server

import (
	"go.luxshare-ict.com/app/service/account/internal/service"
	"go.luxshare-ict.com/pkg/config"
	protocol "go.luxshare-ict.com/pkg/protocols"
)

var servers = make(map[string]server)
var svc *service.Service

type server interface {
	Name() string
	Start(protocol protocol.Protocol) (err error)
}

func Init(cfg *config.Conf) {
	svc = service.New(cfg)
}

func Close() {
	svc.Close()
}

func register(p server) {
	servers[p.Name()] = p
}

func Get(name string) server {
	return servers[name]
}

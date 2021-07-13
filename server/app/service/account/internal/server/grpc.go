package server

import (
	v1 "go.luxshare-ict.com/app/service/account/api/v1"
	protocol "go.luxshare-ict.com/pkg/protocols"
)

type grpc struct {
	name string
}

func init() {
	register(&grpc{name: "grpc"})
}
func (p *grpc) Name() string {
	return p.name
}

func (p *grpc) Start(protocol protocol.Protocol) (err error) {
	e := protocol.Server().Grpc
	v1.RegisterAccountServer(e, svc.Account)
	return protocol.Mount()
}

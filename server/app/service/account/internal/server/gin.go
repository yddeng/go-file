package server

import (
	gonicGin "github.com/gin-gonic/gin"
	ginServer "go.luxshare-ict.com/app/service/account/internal/server/http"
	protocol "go.luxshare-ict.com/pkg/protocols"
)

type gin struct {
	name string
}

func init() {
	register(&gin{name: "gin"})
}
func (p *gin) Name() string {
	return p.name
}

func (p *gin) Start(protocol protocol.Protocol) (err error) {
	ginServer.Init(svc)
	route(protocol.Server().Gin)
	return protocol.Mount()
}
func route(e *gonicGin.Engine) {
	gV1 := e.Group("/v1")
	gV1.POST("/checkToken", ginServer.CheckToken)
}

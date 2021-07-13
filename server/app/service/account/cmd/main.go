package main

import (
	"go.luxshare-ict.com/app/service/account/internal/server"
	"go.luxshare-ict.com/pkg/app"
	bi "go.luxshare-ict.com/pkg/app/buildInfo"
	"go.luxshare-ict.com/pkg/config"
)

var buildTime, version string

type App struct{}

func (a *App) Init(cfg *config.Conf) {
	server.Init(cfg)
}
func (a *App) Close() {
	server.Close()
}
func (a *App) BuildInfo() *bi.BuildInfo {
	return &bi.BuildInfo{
		Time:    buildTime,
		Version: version,
	}
}
func (a *App) GetServer(name string) app.Server {
	return server.Get(name)
}

func main() {
	app.Default(&App{}, "/Users/luxshare/project/luxshare/dev/LuxGoMES/server/app/service/account/cmd/configs/dev.yaml")
}

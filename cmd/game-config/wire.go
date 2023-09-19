//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.

package main

import (
	"game-config/internal/biz"
	"game-config/internal/conf"
	"game-config/internal/data"
	"game-config/internal/middleware"
	"game-config/internal/server"
	"game-config/internal/service"
	"game-config/internal/task"
	"game-config/internal/client"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
)

// initApp init kratos application.
func initApp(*conf.Bootstrap, log.Logger) (*kratos.App, func(), error) {
	panic(wire.Build(server.ProviderSet, data.ProviderSet, biz.ProviderSet, service.ProviderSet, task.ProviderSet, middleware.ProviderSet, client.ProviderSet, newApp))
}

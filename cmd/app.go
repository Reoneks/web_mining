package cmd

import (
	"test/config"
	"test/server"
	"test/server/handlers"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/pkgerrors"
	"go.uber.org/fx"
)

func Exec() fx.Option {
	return fx.Options(
		fx.Provide(
			config.Get,
			handlers.NewHandler,
			server.NewHTTPServer,
		),
		fx.Invoke(
			prepareLogger,
			prepareHooks,
		),
	)
}

func prepareLogger() {
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
}

func prepareHooks(server server.HTTPServer, lc fx.Lifecycle) {
	lc.Append(fx.Hook{OnStart: server.Start, OnStop: server.Stop})
}

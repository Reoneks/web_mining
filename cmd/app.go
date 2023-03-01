package cmd

import (
	"test/config"
	"test/internal/crawler"
	"test/pkg/postgres"
	"test/pkg/whois"
	"test/server"
	"test/server/handlers"

	"go.uber.org/fx"
)

func Exec() fx.Option {
	return fx.Options(
		fx.Provide(
			config.Get,
			postgres.NewPostgres,
			fx.Annotate(annotationDupl[postgres.Postgres], fx.As(new(crawler.Postgres))),
			fx.Annotate(whois.NewWhoIS, fx.As(new(crawler.WhoIS))),
			fx.Annotate(crawler.NewCrawlerBase, fx.As(new(handlers.Crawler))),
			handlers.NewHandler,
			server.NewHTTPServer,
		),
		fx.Invoke(
			prepareHooks,
		),
	)
}

func annotationDupl[T any](v *T) *T {
	return v
}

func prepareHooks(server server.HTTPServer, postgres *postgres.Postgres, lc fx.Lifecycle) {
	lc.Append(fx.Hook{OnStop: postgres.Stop})
	lc.Append(fx.Hook{OnStart: server.Start, OnStop: server.Stop})
}

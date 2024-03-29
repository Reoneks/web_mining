package cmd

import (
	"dyploma/config"
	"dyploma/internal/crawler"
	"dyploma/internal/cron"
	"dyploma/pkg/postgres"
	"dyploma/pkg/whois"
	"dyploma/server"
	"dyploma/server/handlers"

	"go.uber.org/fx"
)

func Exec() fx.Option {
	return fx.Options(
		fx.Provide(
			config.Get,
			postgres.NewPostgres,
			fx.Annotate(
				annotationDupl[postgres.Postgres],
				fx.As(new(crawler.Postgres)),
				fx.As(new(handlers.Postgres)),
				fx.As(new(cron.Postgres)),
			),
			fx.Annotate(whois.NewWhoIS, fx.As(new(crawler.WhoIS))),
			fx.Annotate(crawler.NewCrawlerBase,
				fx.As(new(handlers.Crawler)),
				fx.As(new(cron.Crawler)),
			),
			cron.NewCron,
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

func prepareHooks(server *server.HTTPServer, postgres *postgres.Postgres, _ *cron.Cron, lc fx.Lifecycle) {
	// lc.Append(fx.Hook{OnStart: cron.Start, OnStop: cron.Stop})
	lc.Append(fx.Hook{OnStop: postgres.Stop})
	lc.Append(fx.Hook{OnStart: server.Start, OnStop: server.Stop})
}

package cron

import (
	"test/settings"
	"test/structs"

	"github.com/rs/zerolog/log"
	"github.com/sourcegraph/conc/pool"
)

func (c *Cron) UpdateSites() {
	sites, err := c.postgres.GetSites()
	if err != nil {
		log.Error().Str("function", "UpdateSites").Err(err).Msg("cron.UpdateSites postgres.GetSites error")
		return
	}

	p := pool.New().WithMaxGoroutines(settings.MaxCronGoroutines)
	for _, site := range sites {
		p.Go(func(site structs.SiteStruct) func() {
			return func() {
				_, err := c.crawler.PageWalker(site.BaseURL+"/", site.Exclude, false, true, site.Headers)
				if err != nil {
					log.Error().Str("function", "UpdateSites").Err(err).Msg("cron.UpdateSites crawler.PageWalker error")
				}
			}
		}(site))
	}
}

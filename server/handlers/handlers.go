package handlers

import (
	"dyploma/structs"

	"github.com/sourcegraph/conc/pool"
)

type Crawler interface {
	PageWalker(page string, exclude []string, onlyThisPage, forceCollect bool, headers map[string]string) (siteStruct structs.SiteStruct, err error)
}

type Postgres interface {
	GetCrawlerData(link string) (structs.CrawlerData, error)
}

type Handler struct {
	crawler  Crawler
	postgres Postgres

	p *pool.Pool
}

func NewHandler(crawler Crawler, postgres Postgres) *Handler {
	p := pool.New()
	return &Handler{crawler: crawler, postgres: postgres, p: p}
}

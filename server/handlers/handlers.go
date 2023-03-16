package handlers

import (
	"test/structs"

	"github.com/sourcegraph/conc/pool"
)

type Crawler interface {
	PageWalker(page string, exclude []string, onlyThisPage, forceCollect bool, headers map[string]string) (siteStruct structs.SiteStruct, err error)
}

type Postgres interface {
	GetCrawlerData(link string) (structs.CrawlerData, error)
}

type WSManager interface {
	Send(eventID string, data any) error
}

type Handler struct {
	crawler  Crawler
	postgres Postgres
	ws       WSManager

	p *pool.Pool
}

func NewHandler(crawler Crawler, postgres Postgres, ws WSManager) *Handler {
	p := pool.New()
	return &Handler{crawler: crawler, postgres: postgres, ws: ws, p: p}
}

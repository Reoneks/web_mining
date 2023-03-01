package handlers

import (
	"test/structs"
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
}

func NewHandler(crawler Crawler, postgres Postgres) *Handler {
	return &Handler{crawler: crawler, postgres: postgres}
}

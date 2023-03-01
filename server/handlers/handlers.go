package handlers

import (
	"test/structs"
)

type Crawler interface {
	PageWalker(page string, exclude []string, onlyThisPage, forceCollect bool, headers map[string]string) (siteStruct structs.SiteStruct, err error)
}

type Handler struct {
	crawler Crawler
}

func NewHandler(crawler Crawler) *Handler {
	return &Handler{crawler: crawler}
}

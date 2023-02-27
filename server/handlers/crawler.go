package handlers

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"test/internal/crawler"
	"test/structs"
	"test/tools"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

func (h *Handler) GetSiteStruct(ctx echo.Context) error {
	var siteParseReq structs.SiteParseReq
	if err := ctx.Bind(&siteParseReq); err != nil {
		log.Error().Str("function", "GetSiteStruct").Err(err).Msg("Failed to bind to siteParseReq")
		return ctx.JSON(http.StatusBadRequest, newHTTPError(ErrBind))
	}

	u, err := url.Parse(siteParseReq.URL)
	if err != nil {
		log.Error().Str("function", "GetSiteStruct").Err(err).Msg(ErrParseUrl.Error())
		return ctx.JSON(http.StatusBadRequest, newHTTPError(ErrParseUrl))
	}

	headers := make(map[string]string)
	for key, values := range ctx.Request().Header {
		headers[key] = strings.Join(values, ",")
	}

	crawler := crawler.NewCrawler(siteParseReq.URL, siteParseReq.Exclude)
	hierarchy, err := crawler.PageWalker(siteParseReq.URL, siteParseReq.OnlyThisPage, headers)
	if err != nil {
		log.Error().Str("function", "GetSiteStruct").Err(err).Msg(ErrGetSiteStruct.Error())
		return ctx.JSON(http.StatusInternalServerError, newHTTPError(ErrGetSiteStruct))
	}

	resp := structs.SiteStruct{
		Url:                 siteParseReq.URL,
		BaseURL:             fmt.Sprintf("%s://%s", u.Scheme, u.Host),
		ProcessedHyperlinks: 1,
		StatusCodesCounter:  make(map[int]int64),
	}

	resp.Hierarchy = tools.HierarchyProcess(&resp, &hierarchy, make(map[string]*structs.LinkHierarchy))
	return ctx.JSON(http.StatusOK, resp)
}

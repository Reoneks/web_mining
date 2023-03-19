package handlers

import (
	"net/http"
	"strings"
	"test/structs"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

func (h *Handler) GetSiteStruct(ctx echo.Context) error {
	var siteParseReq structs.SiteParseReq
	if err := ctx.Bind(&siteParseReq); err != nil {
		log.Error().Str("function", "GetSiteStruct").Err(err).Msg("Failed to bind to siteParseReq")
		return ctx.JSON(http.StatusBadRequest, newHTTPError(ErrBind))
	}

	headers := make(map[string]string)
	for key, values := range ctx.Request().Header {
		headers[key] = strings.Join(values, ",")
	}

	resp, err := h.crawler.PageWalker(siteParseReq.URL, siteParseReq.Exclude, siteParseReq.OnlyThisPage, siteParseReq.ForceCollect, headers)
	if err != nil {
		log.Error().Str("function", "GetSiteStruct").Err(err).Msg(ErrGetSiteStruct.Error())
		return ctx.JSON(http.StatusBadRequest, newHTTPError(ErrGetSiteStruct))
	}

	return ctx.JSON(http.StatusOK, resp)
}

func (h *Handler) GetDetails(ctx echo.Context) error {
	resp, err := h.postgres.GetCrawlerData(ctx.QueryParam("link"))
	if err != nil {
		log.Error().Str("function", "GetDetails").Err(err).Msg(ErrGetDetails.Error())
		return ctx.JSON(http.StatusInternalServerError, newHTTPError(ErrGetDetails))
	}

	return ctx.JSON(http.StatusOK, resp)
}

package handlers

import (
	"net/http"
	"strings"
	"test/structs"

	"github.com/google/uuid"
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

	eventID := uuid.NewString()
	h.p.Go(func(siteParseReq structs.SiteParseReq, headers map[string]string, eventID string) func() {
		return func() {
			resp, err := h.crawler.PageWalker(siteParseReq.URL, siteParseReq.Exclude, siteParseReq.OnlyThisPage, siteParseReq.ForceCollect, headers)
			if err != nil {
				log.Error().Str("function", "GetSiteStruct").Err(err).Msg(ErrGetSiteStruct.Error())
				err := h.ws.Send(eventID, "Server internal error")
				if err != nil {
					log.Error().Str("function", "GetSiteStruct").Err(err).Msg("Failed to send ws message")
				}
			}

			err = h.ws.Send(eventID, resp)
			if err != nil {
				log.Error().Str("function", "GetSiteStruct").Err(err).Msg("Failed to send ws message")
			}
		}
	}(siteParseReq, headers, eventID))

	return ctx.String(http.StatusOK, eventID)
}

func (h *Handler) GetDetails(ctx echo.Context) error {
	resp, err := h.postgres.GetCrawlerData(ctx.QueryParam("link"))
	if err != nil {
		log.Error().Str("function", "GetDetails").Err(err).Msg(ErrGetDetails.Error())
		return ctx.JSON(http.StatusInternalServerError, newHTTPError(ErrGetDetails))
	}

	return ctx.JSON(http.StatusOK, resp)
}

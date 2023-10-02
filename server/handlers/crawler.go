package handlers

import (
	"net/http"
	"sort"
	"strings"

	"dyploma/structs"

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

	wordsCounter := make(map[string]int64)
	for _, word := range strings.Split(strings.ReplaceAll(resp.Text, "\n", " "), " ") {
		if len(word) > 3 {
			wordsCounter[strings.ToLower(word)]++
		}
	}

	for k, v := range wordsCounter {
		resp.WordsCounter = append(resp.WordsCounter, structs.WordCount{
			Word:  k,
			Count: v,
		})
	}

	sort.Slice(resp.WordsCounter, func(i, j int) bool {
		return resp.WordsCounter[i].Count > resp.WordsCounter[j].Count
	})

	if len(resp.WordsCounter) > 50 {
		resp.WordsCounter = resp.WordsCounter[:50]
	}

	return ctx.JSON(http.StatusOK, resp)
}

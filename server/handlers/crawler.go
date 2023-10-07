package handlers

import (
	"net/http"
	"strings"

	"dyploma/structs"

	textrank "github.com/DavidBelicza/TextRank/v2"
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

	tr := textrank.NewTextRank()
	rule := textrank.NewDefaultRule()
	language := textrank.NewDefaultLanguage()
	algorithmDef := textrank.NewDefaultAlgorithm()
	tr.Populate(resp.Text, language, rule)
	tr.Ranking(algorithmDef)

	phrases := textrank.FindPhrases(tr)
	p := make([]structs.Phrase, 0, len(phrases))
	for _, phrase := range phrases {
		p = append(p, structs.Phrase{
			Left:   phrase.Left,
			Right:  phrase.Right,
			Weight: phrase.Weight,
			Qty:    phrase.Qty,
		})
	}

	resp.Phrases = p

	words := textrank.FindSingleWords(tr)
	w := make([]structs.Word, 0, len(words))
	for _, word := range words {
		w = append(w, structs.Word{
			Word:   word.Word,
			Weight: word.Weight,
			Qty:    word.Qty,
		})
	}

	resp.Words = w

	sentences := textrank.FindSentencesByRelationWeight(tr, 20)
	for _, sentence := range sentences {
		resp.Sentences = append(resp.Sentences, sentence.Value)
	}

	return ctx.JSON(http.StatusOK, resp)
}

package crawler

import (
	"errors"
	"fmt"
	"net/url"
	"reflect"

	"dyploma/structs"
	"dyploma/tools"

	"github.com/jtarchie/pagerank"
	"github.com/lib/pq"
	whoisparser "github.com/likexian/whois-parser"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

type Postgres interface {
	SaveSiteStruct(siteStruct structs.SiteStruct, onlyThisPage, force bool) error
	GetFullData(link, url string, onlyThisPage bool) (structs.SiteStruct, error)
	GetCrawlerData(link string) (structs.CrawlerData, error)
}

type WhoIS interface {
	WhoIS(site *url.URL) (whoisparser.WhoisInfo, error)
}

type Base struct {
	postgres Postgres
	whois    WhoIS
}

func (cb *Base) PageWalker(page string, exclude []string, onlyThisPage, forceCollect bool, headers map[string]string) (structs.SiteStruct, error) {
	pageParsedURL, err := url.Parse(page)
	if err != nil {
		return structs.SiteStruct{}, fmt.Errorf("CrawlerBase.PageWalker url parse error: %w", err)
	}

	if page == fmt.Sprintf("%s://%s", pageParsedURL.Scheme, pageParsedURL.Host) {
		page += "/"
	}

	var siteStruct structs.SiteStruct
	if !forceCollect {
		siteStruct, err = cb.postgres.GetFullData(fmt.Sprintf("%s://%s", pageParsedURL.Scheme, pageParsedURL.Host), page, onlyThisPage)
		if onlyThisPage && err == nil {
			data, err := cb.postgres.GetCrawlerData(page)
			if err == nil {
				siteStruct.Hierarchy = &structs.Hierarchy{CrawlerData: data}
			}
		}
	}

	if err != nil || reflect.DeepEqual(siteStruct, structs.SiteStruct{}) || forceCollect {
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			log.Error().Str("function", "PageWalker").Err(err).Msg("CrawlerBase.PageWalker postgres GetFullData error")
		}

		whoisParsed, err := cb.whois.WhoIS(pageParsedURL)
		if err != nil {
			log.Error().Str("function", "PageWalker").Err(err).Msg("CrawlerBase.PageWalker whois error")
		} else if whoisParsed.Domain == nil {
			whoisParsed.Domain = new(whoisparser.Domain)
			log.Error().Str("function", "PageWalker").Msg("CrawlerBase.PageWalker whoisParsed Domain empty")
		}

		crawler := newCrawler(page, exclude)
		hierarchy, err := crawler.PageWalker(page, onlyThisPage, headers)
		if err != nil {
			log.Error().Str("function", "PageWalker").Err(err).Msg("CrawlerBase.PageWalker url parse error")
		}

		hierarchy.ParentLink = fmt.Sprintf("%s://%s", pageParsedURL.Scheme, pageParsedURL.Host)
		siteStruct = structs.SiteStruct{
			DomainID:       whoisParsed.Domain.ID,
			BaseURL:        fmt.Sprintf("%s://%s", pageParsedURL.Scheme, pageParsedURL.Host),
			Punycode:       whoisParsed.Domain.Punycode,
			DNSSec:         whoisParsed.Domain.DNSSec,
			NameServers:    whoisParsed.Domain.NameServers,
			Status:         whoisParsed.Domain.Status,
			WhoisServer:    whoisParsed.Domain.WhoisServer,
			CreatedDate:    whoisParsed.Domain.CreatedDate,
			UpdatedDate:    whoisParsed.Domain.UpdatedDate,
			ExpirationDate: whoisParsed.Domain.ExpirationDate,
			Registrar:      whoisParsed.Registrar,
			Registrant:     whoisParsed.Registrant,
			Administrative: whoisParsed.Administrative,
			Technical:      whoisParsed.Technical,
			Billing:        whoisParsed.Billing,
			Hierarchy:      &hierarchy,
			Exclude:        pq.StringArray(exclude),
			Headers:        headers,
		}

		if err := cb.postgres.SaveSiteStruct(siteStruct, onlyThisPage, forceCollect); err != nil {
			log.Error().Str("function", "PageWalker").Err(err).Msg("CrawlerBase.PageWalker postgres SaveSiteStruct error")
		}
	} else {
		tools.SetParents(siteStruct.Hierarchy)
	}

	siteStruct.URL = page
	siteStruct.ProcessedHyperlinks = 1
	siteStruct.StatusCodesCounter = make(map[int]int64)

	graph := pagerank.NewGraph[string]()
	siteStruct.LinkHierarchy = tools.HierarchyProcess(&siteStruct, siteStruct.Hierarchy, make(map[string][]structs.LinkHierarchy), graph)

	ranks := make(map[string]float64)
	graph.Rank(0.85, 0.000001, func(node string, rank float64) {
		ranks[node] = rank
	})

	tools.SetPageRank(&siteStruct.LinkHierarchy, ranks)
	siteStruct.UniqueHyperlinks = tools.UniqueHyperlinks(siteStruct.Hierarchy)
	return siteStruct, nil
}

func NewCrawlerBase(postgres Postgres, whois WhoIS) *Base {
	return &Base{postgres: postgres, whois: whois}
}

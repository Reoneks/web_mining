package crawler

import (
	"fmt"
	"net/url"
	"reflect"
	"test/structs"
	"test/tools"

	whoisparser "github.com/likexian/whois-parser"
	"github.com/rs/zerolog/log"
)

type Postgres interface {
	SaveSiteStruct(siteStruct structs.SiteStruct) error
	GetFullData(link string) (structs.SiteStruct, error)
}

type WhoIS interface {
	WhoIS(site *url.URL) (whoisparser.WhoisInfo, error)
}

type CrawlerBase struct {
	postgres Postgres
	whois    WhoIS
}

func (cb *CrawlerBase) PageWalker(page string, exclude []string, onlyThisPage bool, headers map[string]string) (structs.SiteStruct, error) {
	siteStruct, err := cb.postgres.GetFullData(page)
	if err != nil || reflect.DeepEqual(siteStruct, structs.SiteStruct{}) {
		if err != nil {
			log.Error().Str("function", "PageWalker").Err(err).Msg("CrawlerBase.PageWalker postgres GetFullData error")
		}

		pageParsedUrl, err := url.Parse(page)
		if err != nil {
			return structs.SiteStruct{}, fmt.Errorf("CrawlerBase.PageWalker url parse error: %w", err)
		}

		whoisParsed, err := cb.whois.WhoIS(pageParsedUrl)
		if err != nil {
			return structs.SiteStruct{}, fmt.Errorf("CrawlerBase.PageWalker whois get error: %w", err)
		}

		crawler := newCrawler(page, exclude, cb.postgres)
		siteStruct.Hierarchy, err = crawler.PageWalker(page, onlyThisPage, headers)
		if err != nil {
			return structs.SiteStruct{}, fmt.Errorf("CrawlerBase.PageWalker url parse error: %w", err)
		}

		siteStruct = structs.SiteStruct{
			DomainID:       whoisParsed.Domain.ID,
			Url:            page,
			BaseURL:        fmt.Sprintf("%s://%s", pageParsedUrl.Scheme, pageParsedUrl.Host),
			Punycode:       whoisParsed.Domain.Punycode,
			DNSSec:         whoisParsed.Domain.DNSSec,
			NameServers:    whoisParsed.Domain.NameServers,
			Status:         whoisParsed.Domain.Status,
			WhoisServer:    whoisParsed.Domain.WhoisServer,
			CreatedDate:    whoisParsed.Domain.CreatedDate,
			UpdatedDate:    whoisParsed.Domain.UpdatedDate,
			ExpirationDate: whoisParsed.Domain.ExpirationDate,
		}

		if err := cb.postgres.SaveSiteStruct(siteStruct); err != nil {
			log.Error().Str("function", "PageWalker").Err(err).Msg("CrawlerBase.PageWalker postgres SaveSiteStruct error")
		}
	}

	siteStruct.ProcessedHyperlinks = 1
	siteStruct.StatusCodesCounter = make(map[int]int64)
	siteStruct.LinkHierarchy = tools.HierarchyProcess(&siteStruct, &siteStruct.Hierarchy, make(map[string]*structs.LinkHierarchy))
	return siteStruct, nil
}

func NewCrawlerBase(postgres Postgres, whois WhoIS) *CrawlerBase {
	return &CrawlerBase{postgres: postgres, whois: whois}
}

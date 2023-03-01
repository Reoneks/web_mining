package crawler

import (
	"errors"
	"fmt"
	"net/url"
	"reflect"
	"test/structs"
	"test/tools"

	"github.com/lib/pq"
	whoisparser "github.com/likexian/whois-parser"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
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

func (cb *CrawlerBase) PageWalker(page string, exclude []string, onlyThisPage, forceCollect bool, headers map[string]string) (structs.SiteStruct, error) {
	pageParsedUrl, err := url.Parse(page)
	if err != nil {
		return structs.SiteStruct{}, fmt.Errorf("CrawlerBase.PageWalker url parse error: %w", err)
	}

	siteStruct, err := cb.postgres.GetFullData(fmt.Sprintf("%s://%s", pageParsedUrl.Scheme, pageParsedUrl.Host))
	if err != nil || reflect.DeepEqual(siteStruct, structs.SiteStruct{}) || forceCollect {
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			log.Error().Str("function", "PageWalker").Err(err).Msg("CrawlerBase.PageWalker postgres GetFullData error")
		}

		whoisParsed, err := cb.whois.WhoIS(pageParsedUrl)
		if err != nil {
			return structs.SiteStruct{}, fmt.Errorf("CrawlerBase.PageWalker whois get error: %w", err)
		}

		crawler := newCrawler(page, exclude)
		hierarchy, err := crawler.PageWalker(page, onlyThisPage, headers)
		if err != nil {
			return structs.SiteStruct{}, fmt.Errorf("CrawlerBase.PageWalker url parse error: %w", err)
		}

		hierarchy.ParentLink = fmt.Sprintf("%s://%s", pageParsedUrl.Scheme, pageParsedUrl.Host)
		siteStruct = structs.SiteStruct{
			DomainID:       whoisParsed.Domain.ID,
			BaseURL:        fmt.Sprintf("%s://%s", pageParsedUrl.Scheme, pageParsedUrl.Host),
			Punycode:       whoisParsed.Domain.Punycode,
			DNSSec:         whoisParsed.Domain.DNSSec,
			NameServers:    whoisParsed.Domain.NameServers,
			Status:         whoisParsed.Domain.Status,
			WhoisServer:    whoisParsed.Domain.WhoisServer,
			CreatedDate:    whoisParsed.Domain.CreatedDate,
			UpdatedDate:    whoisParsed.Domain.UpdatedDate,
			ExpirationDate: whoisParsed.Domain.ExpirationDate,
			Hierarchy:      &hierarchy,
			Exclude:        pq.StringArray(exclude),
			Headers:        headers,
		}

		if err := cb.postgres.SaveSiteStruct(siteStruct); err != nil {
			log.Error().Str("function", "PageWalker").Err(err).Msg("CrawlerBase.PageWalker postgres SaveSiteStruct error")
		}
	}

	siteStruct.Url = page
	siteStruct.ProcessedHyperlinks = 1
	siteStruct.StatusCodesCounter = make(map[int]int64)
	siteStruct.LinkHierarchy = tools.HierarchyProcess(&siteStruct, siteStruct.Hierarchy)
	return siteStruct, nil
}

func NewCrawlerBase(postgres Postgres, whois WhoIS) *CrawlerBase {
	return &CrawlerBase{postgres: postgres, whois: whois}
}

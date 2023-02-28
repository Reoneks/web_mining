package whois

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/likexian/whois"
	whoisparser "github.com/likexian/whois-parser"
)

type WhoIS struct{}

func (wi *WhoIS) WhoIS(site *url.URL) (whoisparser.WhoisInfo, error) {
	whoisRaw, err := whois.Whois(strings.TrimPrefix(site.Host, "www."))
	if err != nil {
		return whoisparser.WhoisInfo{}, fmt.Errorf("WhoIS get error: %w", err)
	}

	whoisParsed, err := whoisparser.Parse(whoisRaw)
	if err != nil {
		return whoisparser.WhoisInfo{}, fmt.Errorf("WhoIS parse error: %w", err)
	}

	return whoisParsed, nil
}

func NewWhoIS() *WhoIS {
	return &WhoIS{}
}

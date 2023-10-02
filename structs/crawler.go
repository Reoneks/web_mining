package structs

import (
	"slices"
	"strings"

	"dyploma/settings"

	"github.com/lib/pq"
	whoisparser "github.com/likexian/whois-parser"
)

type LinkHierarchy struct {
	Link       string            `json:"name"`
	Attributes map[string]string `json:"attributes"`
	Children   []LinkHierarchy   `json:"children"`
}

type SiteStruct struct {
	DomainID            string               `json:"domain_id"`
	URL                 string               `json:"url" gorm:"-"`
	BaseURL             string               `json:"base_url" gorm:"primary_key"`
	Punycode            string               `json:"punycode"`
	DNSSec              bool                 `json:"dns_sec"`
	NameServers         pq.StringArray       `json:"name_servers" gorm:"type:text[]"`
	Status              pq.StringArray       `json:"status" gorm:"type:text[]"`
	WhoisServer         string               `json:"whois_server"`
	Registrar           *whoisparser.Contact `json:"registrar" gorm:"serializer:json"`
	Registrant          *whoisparser.Contact `json:"registrant" gorm:"serializer:json"`
	Administrative      *whoisparser.Contact `json:"administrative" gorm:"serializer:json"`
	Technical           *whoisparser.Contact `json:"technical" gorm:"serializer:json"`
	Billing             *whoisparser.Contact `json:"billing" gorm:"serializer:json"`
	Images              int64                `json:"images" gorm:"-"`
	VideoLinks          int64                `json:"video" gorm:"-"`
	AudioLinks          int64                `json:"audio" gorm:"-"`
	Files               int64                `json:"files" gorm:"-"`
	Fonts               int64                `json:"fonts" gorm:"-"`
	Hyperlinks          int64                `json:"hyperlinks" gorm:"-"`
	UniqueHyperlinks    int64                `json:"unique_hyperlinks" gorm:"-"`
	ProcessedHyperlinks int64                `json:"processed_hyperlinks" gorm:"-"`
	InternalLinks       int64                `json:"internal_links" gorm:"-"`
	Symbols             int64                `json:"symbols" gorm:"-"`
	Words               int64                `json:"words" gorm:"-"`
	Paragraphs          int64                `json:"paragraphs" gorm:"-"`
	Errors              int64                `json:"errors" gorm:"-"`
	StatusCodesCounter  map[int]int64        `json:"status_codes" gorm:"-"`
	CreatedDate         string               `json:"created_date"`
	UpdatedDate         string               `json:"updated_date"`
	ExpirationDate      string               `json:"expiration_date"`
	LinkHierarchy       LinkHierarchy        `json:"hierarchy" gorm:"-"`

	Hierarchy *Hierarchy        `json:"-" gorm:"foreignKey:ParentLink;references:BaseURL"`
	Exclude   pq.StringArray    `json:"-" gorm:"type:text[]"`
	Headers   map[string]string `json:"-" gorm:"serializer:json"`
}

func (SiteStruct) TableName() string {
	return "sites"
}

type Hierarchy struct {
	CrawlerData `gorm:"embedded"`
	Childrens   []Hierarchy `json:"childrens" gorm:"foreignKey:ParentLink;references:Link"`
	ParentLink  string      `json:"-"`
	Parent      *Hierarchy  `json:"-" gorm:"-"`
}

func (Hierarchy) TableName() string {
	return "link_data"
}

type CrawlerData struct {
	Link          string              `json:"link" gorm:"primary_key"`
	StatusCode    int                 `json:"status_code"`
	Error         string              `json:"error"`
	Text          string              `json:"text"`
	Ping          float64             `json:"ping"`
	Images        pq.StringArray      `json:"images" gorm:"type:text[]"`
	Audio         pq.StringArray      `json:"audio" gorm:"type:text[]"`
	Video         pq.StringArray      `json:"video" gorm:"type:text[]"`
	Fonts         pq.StringArray      `json:"fonts" gorm:"type:text[]"`
	Files         pq.StringArray      `json:"files" gorm:"type:text[]"`
	Hyperlinks    pq.StringArray      `json:"hyperlinks" gorm:"type:text[]"`
	InternalLinks pq.StringArray      `json:"internal_links" gorm:"type:text[]"`
	Metadata      []map[string]string `json:"metadata" gorm:"serializer:json"`
	WordsCounter  []WordCount         `json:"words_counter" gorm:"-"`
}

func (cd *CrawlerData) Merge(isHTMLBlock bool, toMerge ...CrawlerData) {
	for _, merge := range toMerge {
		cd.mergeText(merge.Text, isHTMLBlock)
		cd.Images = append(cd.Images, merge.Images...)
		cd.Audio = append(cd.Audio, merge.Audio...)
		cd.Video = append(cd.Video, merge.Video...)
		cd.Fonts = append(cd.Fonts, merge.Fonts...)
		cd.Files = append(cd.Files, merge.Files...)
		cd.Hyperlinks = append(cd.Hyperlinks, merge.Hyperlinks...)
		cd.InternalLinks = append(cd.InternalLinks, merge.InternalLinks...)
		cd.Metadata = append(cd.Metadata, merge.Metadata...)
	}
}

func (cd *CrawlerData) mergeText(text string, isHTMLBlock bool) {
	if text != "" {
		if cd.Text != "" && isHTMLBlock && !strings.HasSuffix(cd.Text, "\n") {
			cd.Text += "\n"
		}

		if cd.Text != "" && !strings.HasSuffix(cd.Text, "\n") && !slices.Contains(settings.DontAddSpaceWhereSymbolRight, string(text[0])) && !slices.Contains(settings.DontAddSpaceWhereSymbolLeft, string(cd.Text[len(cd.Text)-1])) {
			cd.Text += " "
		}

		cd.Text += text
		cd.Text = strings.TrimPrefix(cd.Text, " ")

		if isHTMLBlock && !strings.HasSuffix(cd.Text, "\n") {
			cd.Text += "\n"
		}

		cd.Text = strings.ToValidUTF8(cd.Text, "")
	}
}

type WordCount struct {
	Word  string
	Count int64
}

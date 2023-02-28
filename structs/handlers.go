package structs

type SiteParseReq struct {
	URL          string   `query:"url"`
	OnlyThisPage bool     `query:"onlyThisPage"`
	Exclude      []string `query:"exclude"`
}

type LinkHierarchy struct {
	Link       string            `json:"name"`
	Attributes map[string]string `json:"attributes"`
	Children   []LinkHierarchy   `json:"children"`
}

type SiteStruct struct {
	ID                  string        `json:"id"`
	Url                 string        `json:"url"`
	BaseURL             string        `json:"base_url"`
	Punycode            string        `json:"punycode"`
	DNSSec              bool          `json:"dns_sec"`
	NameServers         []string      `json:"name_servers"`
	Status              []string      `json:"status"`
	WhoisServer         string        `json:"whois_server"`
	Images              int64         `json:"images"`
	VideoLinks          int64         `json:"video_links"`
	AudioLinks          int64         `json:"audio_links"`
	Hyperlinks          int64         `json:"hyperlinks"`
	ProcessedHyperlinks int64         `json:"processed_hyperlinks"`
	InternalLinks       int64         `json:"internal_links"`
	Symbols             int64         `json:"symbols"`
	Words               int64         `json:"words"`
	Paragraphs          int64         `json:"paragraphs"`
	Errors              int64         `json:"errors"`
	StatusCodesCounter  map[int]int64 `json:"status_codes"`
	CreatedDate         string        `json:"created_date"`
	UpdatedDate         string        `json:"updated_date"`
	ExpirationDate      string        `json:"expiration_date"`

	Hierarchy LinkHierarchy `json:"hierarchy"`
}

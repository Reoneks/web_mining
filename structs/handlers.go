package structs

type SiteParseReq struct {
	URL          string   `query:"url"`
	OnlyThisPage bool     `query:"only_this_page"`
	ForceCollect bool     `query:"force_collect"`
	Exclude      []string `query:"exclude"`
}

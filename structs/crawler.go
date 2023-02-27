package structs

import (
	"strings"
	"test/settings"

	"golang.org/x/exp/slices"
)

type Hierarchy struct {
	Data      CrawlerData `json:"data"`
	Childrens []Hierarchy `json:"childrens"`
	Parent    *Hierarchy  `json:"parent"`
}

type CrawlerData struct {
	Link          string              `json:"link"`
	StatusCode    int                 `json:"status_code"`
	Error         string              `json:"error"`
	Text          string              `json:"text"`
	Images        []string            `json:"images"`
	Audio         []string            `json:"audio"`
	Video         []string            `json:"video"`
	Hyperlinks    []string            `json:"hyperlinks"`
	InternalLinks []string            `json:"internal_links"`
	Metadata      []map[string]string `json:"metadata"`
}

func (cd *CrawlerData) Merge(isHTMLBlock bool, ToMerge ...CrawlerData) {
	for _, merge := range ToMerge {
		cd.mergeText(merge.Text, isHTMLBlock)
		cd.Images = slices.Compact(append(cd.Images, merge.Images...))
		cd.Audio = slices.Compact(append(cd.Audio, merge.Audio...))
		cd.Video = slices.Compact(append(cd.Video, merge.Video...))
		cd.Hyperlinks = slices.Compact(append(cd.Hyperlinks, merge.Hyperlinks...))
		cd.InternalLinks = slices.Compact(append(cd.InternalLinks, merge.InternalLinks...))
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
	}
}

package tools

import (
	"slices"
	"strings"

	"dyploma/structs"

	textrank "github.com/DavidBelicza/TextRank/v2"
	"github.com/jtarchie/pagerank"
	"github.com/spf13/cast"
)

func PrepareCrawlerText(text string) string {
	text = strings.TrimSpace(strings.ReplaceAll(text, "  ", " "))
	if text == "\n" {
		return ""
	}

	return text
}

func CheckWisited(wisited []string, link string) bool {
	for _, wisited := range wisited {
		if strings.HasSuffix(wisited, "*") {
			wisited = wisited[:len(wisited)-1]
			if strings.HasPrefix(link, wisited) {
				return true
			}
		} else if wisited == link {
			return true
		}
	}

	return false
}

func HierarchyProcess(
	resp *structs.SiteStruct,
	hierarchy *structs.Hierarchy,
	hyperlinks map[string][]structs.LinkHierarchy,
	graph *pagerank.Graph[string],
) (res structs.LinkHierarchy) {
	if resp == nil || hierarchy == nil {
		return structs.LinkHierarchy{}
	}

	resp.Images += int64(len(hierarchy.Images))
	resp.AudioLinks += int64(len(hierarchy.Audio))
	resp.VideoLinks += int64(len(hierarchy.Video))
	resp.Files += int64(len(hierarchy.Files))
	resp.Fonts += int64(len(hierarchy.Fonts))
	resp.Hyperlinks += int64(len(hierarchy.Hyperlinks))
	resp.ProcessedHyperlinks += int64(len(hierarchy.Childrens))
	resp.InternalLinks += int64(len(hierarchy.InternalLinks))
	resp.Paragraphs += int64(len(strings.Split(hierarchy.Text, "\n")))
	resp.Words += int64(len(strings.Split(strings.ReplaceAll(hierarchy.Text, "\n", " "), " ")))
	resp.Symbols += int64(len(strings.ReplaceAll(hierarchy.Text, "\n", "")))
	resp.StatusCodesCounter[hierarchy.StatusCode]++

	tr := textrank.NewTextRank()
	rule := textrank.NewDefaultRule()
	language := textrank.NewDefaultLanguage()
	algorithmDef := textrank.NewDefaultAlgorithm()
	tr.Populate(hierarchy.Text, language, rule)
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

	hierarchy.Phrases = p

	words := textrank.FindSingleWords(tr)
	w := make([]structs.Word, 0, len(words))
	for _, word := range words {
		w = append(w, structs.Word{
			Word:   word.Word,
			Weight: word.Weight,
			Qty:    word.Qty,
		})
	}

	hierarchy.Words = w

	sentences := textrank.FindSentencesByRelationWeight(tr, 20)
	for _, sentence := range sentences {
		hierarchy.Sentences = append(hierarchy.Sentences, sentence.Value)
	}

	if hierarchy.Error != "" {
		resp.Errors++
	}

	res.Link = hierarchy.Link
	res.Attributes = make(map[string]string)

	childrenLinks := make([]string, 0, len(hierarchy.Childrens))
	for i, child := range hierarchy.Childrens {
		childrenLinks = append(childrenLinks, child.Link)
		graph.Link(hierarchy.Link, child.Link, 1.0)
		processed := HierarchyProcess(resp, &hierarchy.Childrens[i], hyperlinks, graph)
		res.Children = append(res.Children, processed)

		if links, ok := hyperlinks[child.Link]; ok {
			var path string
			for root := &hierarchy.Childrens[i]; root != nil; root = root.Parent {
				path = root.Link + " | " + path
			}

			path = strings.TrimSuffix(path, " | ")
			for _, link := range links {
				link.Attributes["Already processed"] = path
			}
		}
	}

	for _, link := range hierarchy.Hyperlinks {
		if strings.Contains(link, resp.BaseURL) && !slices.Contains(childrenLinks, link) {
			linkHierarchy := structs.LinkHierarchy{Link: link, Attributes: make(map[string]string)}

			if links, ok := hyperlinks[link]; ok && len(links) > 0 && len(links[0].Attributes) > 0 {
				linkHierarchy.Attributes["Already processed"] = links[0].Attributes["Already processed"]
			} else {
				hyperlinks[link] = append(hyperlinks[link], linkHierarchy)
			}

			res.Children = append(res.Children, linkHierarchy)
		}
	}

	return res
}

func SetPageRank(hierarchy *structs.LinkHierarchy, ranks map[string]float64) {
	hierarchy.Attributes["PageRank"] = cast.ToString(ranks[hierarchy.Link])

	for i := range hierarchy.Children {
		SetPageRank(&hierarchy.Children[i], ranks)
	}
}

func SetParents(hierarchy *structs.Hierarchy) {
	if hierarchy == nil {
		return
	}

	for i := range hierarchy.Childrens {
		hierarchy.Childrens[i].Parent = hierarchy
		SetParents(&hierarchy.Childrens[i])
	}
}

func UniqueHyperlinks(hierarchy *structs.Hierarchy) int64 {
	if hierarchy == nil {
		return 0
	}

	return int64(len(uniqueHyperlinksProcessor(hierarchy)))
}

func uniqueHyperlinksProcessor(hierarchy *structs.Hierarchy) []string {
	var result []string

	result = slices.Clone(hierarchy.Hyperlinks)
	for i := range hierarchy.Childrens {
		links := uniqueHyperlinksProcessor(&hierarchy.Childrens[i])
		result = append(result, links...)
	}

	return Compact(result)
}

func PrepareLinks(links []string, baseURL string) []string {
	for i, link := range links {
		if !strings.Contains(link, "http://") && !strings.Contains(link, "https://") {
			url := strings.Split(baseURL, "/")
			if len(url) > 0 {
				links[i] = strings.Join(url[:len(url)-1], "/") + link
			}
		}
	}

	return links
}

func Compact[T comparable](array []T) []T {
	res := make([]T, 0, len(array))

	for _, elem := range array {
		if !slices.Contains(res, elem) {
			res = append(res, elem)
		}
	}

	return slices.Clip(res)
}

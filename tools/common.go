package tools

import (
	"strings"
	"test/structs"

	"golang.org/x/exp/slices"
)

func PrepareCrawlerText(text string) string {
	text = strings.TrimSpace(strings.Replace(text, "  ", " ", -1))
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
		} else {
			if wisited == link {
				return true
			}
		}
	}

	return false
}

func HierarchyProcess(resp *structs.SiteStruct, hierarchy *structs.Hierarchy) (res structs.LinkHierarchy) {
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
	resp.StatusCodesCounter[hierarchy.StatusCode] = resp.StatusCodesCounter[hierarchy.StatusCode] + 1

	if hierarchy.Error != "" {
		resp.Errors += 1
	}

	res.Link = hierarchy.Link
	for _, link := range hierarchy.Hyperlinks {
		if strings.Contains(link, resp.BaseURL) {
			linkHierarchy := structs.LinkHierarchy{Link: link}
			res.Children = append(res.Children, linkHierarchy)
		}
	}

	for _, child := range hierarchy.Childrens {
		processed := HierarchyProcess(resp, &child)
		res.Children = append(res.Children, processed)
	}

	return
}

func UniqueHyperlinks(hierarchy *structs.Hierarchy) int64 {
	return int64(len(uniqueHyperlinksProcessor(hierarchy)))
}

func uniqueHyperlinksProcessor(hierarchy *structs.Hierarchy) []string {
	var result []string

	result = slices.Clone(hierarchy.Hyperlinks)
	for _, child := range hierarchy.Childrens {
		links := uniqueHyperlinksProcessor(&child)
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

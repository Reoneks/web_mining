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

func HierarchyProcess(resp *structs.SiteStruct, hierarchy *structs.Hierarchy, hyperlinks map[string]*structs.LinkHierarchy) (res structs.LinkHierarchy) {
	resp.Images += int64(len(hierarchy.Data.Images))
	resp.AudioLinks += int64(len(hierarchy.Data.Audio))
	resp.VideoLinks += int64(len(hierarchy.Data.Video))
	resp.Hyperlinks += int64(len(hierarchy.Data.Hyperlinks))
	resp.ProcessedHyperlinks += int64(len(hierarchy.Childrens))
	resp.InternalLinks += int64(len(hierarchy.Data.InternalLinks))
	resp.Paragraphs += int64(len(strings.Split(hierarchy.Data.Text, "\n")))
	resp.Words += int64(len(strings.Split(strings.ReplaceAll(hierarchy.Data.Text, "\n", " "), " ")))
	resp.Symbols += int64(len(strings.ReplaceAll(hierarchy.Data.Text, "\n", "")))
	resp.StatusCodesCounter[hierarchy.Data.StatusCode] = resp.StatusCodesCounter[hierarchy.Data.StatusCode] + 1

	if hierarchy.Data.Error != "" {
		resp.Errors += 1
	}

	res.Link = hierarchy.Data.Link
	for _, link := range hierarchy.Data.Hyperlinks {
		linkHierarchy := structs.LinkHierarchy{Link: link}
		hyperlinks[link] = &linkHierarchy
		res.Children = append(res.Children, linkHierarchy)
	}

	for _, child := range hierarchy.Childrens {
		processed := HierarchyProcess(resp, &child, hyperlinks)
		res.Children = append(res.Children, processed)

		if link, ok := hyperlinks[child.Data.Link]; ok {
			link.Children = slices.Clone(processed.Children)
		}
	}

	return
}

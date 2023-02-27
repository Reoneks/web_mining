package tools

import (
	"strings"
	"test/structs"
)

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
	res.Attributes = make(map[string]string)
	for _, link := range hierarchy.Data.Hyperlinks {
		linkHierarchy := structs.LinkHierarchy{Link: link, Attributes: make(map[string]string)}
		hyperlinks[link] = &linkHierarchy
		res.Children = append(res.Children, linkHierarchy)
	}

	for _, child := range hierarchy.Childrens {
		res.Children = append(res.Children, HierarchyProcess(resp, &child, hyperlinks))

		if link, ok := hyperlinks[child.Data.Link]; ok {
			var path string
			for root := &child; root != nil; root = root.Parent {
				path += root.Data.Link + " | "
			}

			link.Attributes["Already processed"] = strings.TrimPrefix(path, " | ")
		}
	}

	return
}

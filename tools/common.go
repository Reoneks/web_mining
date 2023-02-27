package tools

import "strings"

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

package crawler

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/url"
	"strings"
	"test/settings"
	"test/structs"
	"test/tools"

	"github.com/go-resty/resty/v2"
	"github.com/rs/zerolog/log"
	"golang.org/x/exp/slices"
	"golang.org/x/net/html"
)

type Crawler struct {
	wisited []string
}

func (c *Crawler) PageWalker(page string, onlyThisPage bool, headers map[string]string) (hierarchy structs.Hierarchy, err error) {
	resp, err := resty.New().SetHeaders(headers).R().Get(page)
	if err != nil {
		hierarchy.Data.Link = page
		hierarchy.Data.StatusCode = resp.StatusCode()
		hierarchy.Data.Error = err.Error()
		return
	}

	u, err := url.Parse(page)
	if err != nil {
		hierarchy.Data.Link = page
		hierarchy.Data.StatusCode = resp.StatusCode()
		hierarchy.Data.Error = err.Error()
		return
	}

	data, err := c.ParsePage(bytes.NewReader(resp.Body()), u)
	if err != nil {
		return
	}

	data.Link = page
	data.StatusCode = resp.StatusCode()

	hierarchy.Data = data
	if !onlyThisPage {
		process := make([]string, 0, len(data.Hyperlinks))
		for _, link := range data.Hyperlinks {
			if strings.Contains(link, fmt.Sprintf("%s://%s", u.Scheme, u.Host)) && !tools.CheckWisited(c.wisited, link) {
				c.wisited = append(c.wisited, link)
				process = append(process, link)
			}
		}

		for _, link := range process {
			child, err := c.PageWalker(link, onlyThisPage, headers)
			if err != nil {
				log.Error().Str("function", "PageWalker").Err(err).Msg("PageWalker error")
			}

			child.Parent = &hierarchy
			hierarchy.Childrens = append(hierarchy.Childrens, child)
		}
	}

	return
}

func (c *Crawler) ParsePage(page io.Reader, baseURL *url.URL) (structs.CrawlerData, error) {
	if page == nil || baseURL == nil {
		return structs.CrawlerData{}, errors.New("Page or URL is null")
	}

	doc, err := html.Parse(page)
	if err != nil {
		return structs.CrawlerData{}, fmt.Errorf("Page parse error: %w", err)
	}

	data := c.crawlerFunc(doc)

	data.Audio = slices.Compact(data.Audio)
	data.Video = slices.Compact(data.Video)
	data.Hyperlinks = slices.Compact(data.Hyperlinks)
	data.Images = slices.Compact(data.Images)
	data.InternalLinks = slices.Compact(data.InternalLinks)

	for i, link := range data.Hyperlinks {
		if strings.HasPrefix(link, "//") {
			data.Hyperlinks[i] = baseURL.Scheme + ":" + link
		} else if strings.HasPrefix(link, "/") {
			data.Hyperlinks[i] = fmt.Sprintf("%s://%s", baseURL.Scheme, baseURL.Host) + link
		} else if strings.HasPrefix(link, "?") {
			baseURL.RawQuery = ""
			data.Hyperlinks[i] = baseURL.String() + link
		}
	}

	for i, video := range data.Video {
		if !strings.Contains(video, "http://") && !strings.Contains(video, "https://") {
			url := strings.Split(baseURL.String(), "/")
			if len(url) > 0 {
				data.Video[i] = strings.Join(url[:len(url)-1], "/") + video
			}
		}
	}

	for i, audio := range data.Audio {
		if !strings.Contains(audio, "http://") && !strings.Contains(audio, "https://") {
			url := strings.Split(baseURL.String(), "/")
			if len(url) > 0 {
				data.Audio[i] = strings.Join(url[:len(url)-1], "/") + audio
			}
		}
	}

	for i, image := range data.Images {
		if !strings.Contains(image, "http://") && !strings.Contains(image, "https://") {
			url := strings.Split(baseURL.String(), "/")
			if len(url) > 0 {
				data.Images[i] = strings.Join(url[:len(url)-1], "/") + image
			}
		}
	}

	return data, nil
}

func (c *Crawler) crawlerFunc(node *html.Node) structs.CrawlerData {
	var data structs.CrawlerData

	if node.Type == html.TextNode && !slices.Contains(settings.NotAllowedTags, node.Parent.Data) {
		text := tools.PrepareCrawlerText(node.Data)
		if text != "" {
			data.Text = text
		}
	} else if node.Type == html.ElementNode {
		if node.Data == "meta" {
			meta := make(map[string]string)

			for _, attr := range node.Attr {
				meta[attr.Key] = attr.Val
			}

			data.Metadata = append(data.Metadata, meta)
		} else if node.Data == "link" || node.Data == "a" {
			for _, attr := range node.Attr {
				if attr.Key == "href" && !strings.Contains(attr.Val, "javascript:void(0)") {
					if strings.HasPrefix(attr.Val, "#") {
						data.InternalLinks = append(data.InternalLinks, attr.Val)
					} else {
						data.Hyperlinks = append(data.Hyperlinks, attr.Val)
					}

					break
				}
			}
		} else if node.Data == "source" {
			for _, attr := range node.Attr {
				if attr.Key == "src" {
					data.Video = append(data.Video, attr.Val)
					break
				}
			}
		} else if node.Data == "img" {
			for _, attr := range node.Attr {
				if attr.Key == "src" && !strings.HasPrefix(attr.Val, "data:image/") {
					data.Images = append(data.Images, attr.Val)
					break
				}
			}
		}
	}

	for child := node.FirstChild; child != nil; child = child.NextSibling {
		toAdd := c.crawlerFunc(child)
		data.Merge(slices.Contains(settings.HTMLBlocks, child.Data), toAdd)
	}

	return data
}

func NewCrawler(firstPage string, exclude []string) *Crawler {
	return &Crawler{wisited: append(exclude, firstPage)}
}

package crawler

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/url"
	"slices"
	"strings"
	"time"

	"dyploma/settings"
	"dyploma/structs"
	"dyploma/tools"

	"github.com/go-resty/resty/v2"
	"github.com/rs/zerolog/log"
	"golang.org/x/net/html"
)

type Crawler struct {
	wisited []string
}

func (c *Crawler) PageWalker(page string, onlyThisPage bool, headers map[string]string) (hierarchy structs.Hierarchy, err error) {
	resp, err := resty.New().SetHeaders(headers).R().Get(page)
	if err != nil {
		hierarchy.Link = page
		hierarchy.StatusCode = resp.StatusCode()
		hierarchy.Error = err.Error()
		return hierarchy, err
	}

	u, err := url.Parse(page)
	if err != nil {
		hierarchy.Link = page
		hierarchy.StatusCode = resp.StatusCode()
		hierarchy.Error = err.Error()
		return hierarchy, err
	}

	data, err := c.ParsePage(bytes.NewReader(resp.Body()), u)
	if err != nil {
		hierarchy.Link = page
		hierarchy.StatusCode = resp.StatusCode()
		hierarchy.Error = err.Error()
		return hierarchy, err
	}

	data.Link = page
	data.StatusCode = resp.StatusCode()

	hierarchy.CrawlerData = data
	hierarchy.Ping = float64(resp.Time()) / float64(time.Second)
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

			child.ParentLink = hierarchy.Link
			child.Parent = &hierarchy
			hierarchy.Childrens = append(hierarchy.Childrens, child)
		}
	}

	return hierarchy, err
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
	for i, link := range data.Hyperlinks {
		switch {
		case strings.HasPrefix(link, "//"):
			data.Hyperlinks[i] = baseURL.Scheme + ":" + link
		case strings.HasPrefix(link, "/"):
			data.Hyperlinks[i] = fmt.Sprintf("%s://%s", baseURL.Scheme, baseURL.Host) + link
		case strings.HasPrefix(link, "?"):
			baseURL.RawQuery = ""
			data.Hyperlinks[i] = baseURL.String() + link

		}
	}

	data.Images = tools.PrepareLinks(data.Images, baseURL.String())
	data.Audio = tools.PrepareLinks(data.Audio, baseURL.String())
	data.Video = tools.PrepareLinks(data.Video, baseURL.String())
	data.Files = tools.PrepareLinks(data.Files, baseURL.String())
	data.Fonts = tools.PrepareLinks(data.Fonts, baseURL.String())

	data.Audio = tools.Compact(data.Audio)
	data.Video = tools.Compact(data.Video)
	data.Hyperlinks = tools.Compact(data.Hyperlinks)
	data.Images = tools.Compact(data.Images)
	data.InternalLinks = tools.Compact(data.InternalLinks)
	data.Files = tools.Compact(data.Files)
	data.Fonts = tools.Compact(data.Fonts)
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
		switch {
		case node.Data == "meta":
			meta := make(map[string]string)

			for _, attr := range node.Attr {
				meta[attr.Key] = attr.Val
			}

			data.Metadata = append(data.Metadata, meta)
		case node.Data == "link" || node.Data == "a":
		LINK_LOOP:
			for _, attr := range node.Attr {
				if attr.Key == "href" && !strings.Contains(attr.Val, "javascript:void(0)") {
					if strings.HasPrefix(attr.Val, "#") {
						data.InternalLinks = append(data.InternalLinks, attr.Val)
					} else {
						for _, imgFormat := range settings.ImageExtensions {
							if strings.HasSuffix(strings.Split(attr.Val, "?")[0], imgFormat) {
								data.Images = append(data.Images, attr.Val)
								break LINK_LOOP
							}
						}

						for _, fontFormat := range settings.FontsExtensions {
							if strings.HasSuffix(strings.Split(attr.Val, "?")[0], fontFormat) {
								data.Fonts = append(data.Fonts, attr.Val)
								break LINK_LOOP
							}
						}

						for _, fileFormat := range settings.FilesExtensions {
							if strings.HasSuffix(strings.Split(attr.Val, "?")[0], fileFormat) {
								data.Files = append(data.Files, attr.Val)
								break LINK_LOOP
							}
						}

						for _, videoFormat := range settings.VideoExtensions {
							if strings.HasSuffix(strings.Split(attr.Val, "?")[0], videoFormat) {
								data.Video = append(data.Video, attr.Val)
								break LINK_LOOP
							}
						}

						for _, audioFormat := range settings.AudioExtensions {
							if strings.HasSuffix(strings.Split(attr.Val, "?")[0], audioFormat) {
								data.Audio = append(data.Audio, attr.Val)
								break LINK_LOOP
							}
						}

						data.Hyperlinks = append(data.Hyperlinks, attr.Val)
					}

					break
				}
			}
		case node.Data == "source":
			var links []string
			for _, attr := range node.Attr {
				if attr.Key == "src" {
					links = append(links, attr.Val)
					break
				}
			}

			if node.Parent.Type == html.ElementNode {
				if node.Parent.Data == "video" {
					data.Video = append(data.Video, links...)
				} else if node.Parent.Data == "audio" {
					data.Audio = append(data.Audio, links...)
				}
			} else if node.Parent.Parent.Type == html.ElementNode {
				if node.Parent.Parent.Data == "video" {
					data.Video = append(data.Video, links...)
				} else if node.Parent.Parent.Data == "audio" {
					data.Audio = append(data.Audio, links...)
				}
			}
		case node.Data == "img":
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

func newCrawler(firstPage string, exclude []string) *Crawler {
	return &Crawler{wisited: append(exclude, firstPage)}
}

package nanoserp

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type SERP struct {
	Results []SearchResult `json:"results"`
}

type SearchResult struct {
	Title   string   `json:"title"`
	URL     string   `json:"url"`
	Content string   `json:"content"`
	Engines []string `json:"engines"`
}

func SearchSearXNG(client *http.Client, endpoint, keyword string, engines []string) ([]SearchResult, error) {
	u, err := url.Parse(endpoint)
	if err != nil {
		return nil, err
	}

	q := u.Query()
	q.Set("q", keyword)
	q.Set("engines", strings.Join(engines, ","))
	u.RawQuery = q.Encode()

	fmt.Println(u.String())

	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/125.0.0.0 Safari/537.36")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")
	req.Header.Set("Priority", "u=0, i")
	req.Header.Set("Sec-Ch-Ua", "\"Google Chrome\";v=\"125\", \"Chromium\";v=\"125\", \"Not.A/Brand\";v=\"24\"")
	req.Header.Set("Sec-Ch-Ua-Mobile", "?0")
	req.Header.Set("Sec-Ch-Ua-Platform", "\"macOS\"")
	req.Header.Set("Sec-Fetch-Dest", "document")
	req.Header.Set("Sec-Fetch-Mode", "navigate")
	req.Header.Set("Sec-Fetch-Site", "none")
	req.Header.Set("Sec-Fetch-User", "?1")
	req.Header.Set("Upgrade-Insecure-Requests", "1")

	if client == nil {
		client = http.DefaultClient
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("status code error: %s", resp.Status)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}

	var results []SearchResult
	doc.Find("article.result").Each(func(i int, s *goquery.Selection) {
		title := strings.TrimSpace(s.Find("h3 a").Text())

		url, _ := s.Find("a.url_wrapper").Attr("href")
		url = strings.TrimSpace(url)

		content := strings.TrimSpace(s.Find("p.content").Text())

		var engines []string
		s.Find("div.engines > span").Each(func(i int, s *goquery.Selection) {
			engines = append(engines, strings.TrimSpace(s.Text()))
		})

		results = append(results, SearchResult{
			Title:   title,
			URL:     url,
			Content: content,
			Engines: engines,
		})
	})

	return results, nil
}

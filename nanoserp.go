package nanoserp

type SERP struct {
	Results []SearchResult `json:"results"`
}

type SearchResult struct {
	Title   string   `json:"title"`
	URL     string   `json:"url"`
	Content string   `json:"content"`
	Engines []string `json:"engines"`
}

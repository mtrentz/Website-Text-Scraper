package scraper

import "sync"

type PageResult struct {
	Url       string `json:"url"`
	Header    string `json:"header"`
	Body      string `json:"text"`
	Footer    string `json:"footer"`
	VisitedAt string `json:"visited_at"`
}

type WebsiteResult struct {
	Url        string       `json:"url"`
	PageAmount int          `json:"page_amount"`
	VisitedAt  string       `json:"visited_at"`
	Pages      []PageResult `json:"pages"`
}

var resultsMutex sync.Mutex

func (wr *WebsiteResult) AddPage(page PageResult) {
	resultsMutex.Lock()
	wr.Pages = append(wr.Pages, page)
	wr.PageAmount++
	resultsMutex.Unlock()
}

package model

type PageCache struct {
	ID        int64  `json:"-"`
	URL       string `json:"url"`
	Title     string `json:"title,omitempty"`
	Markdown  string `json:"markdown"`
	FetchedAt int64  `json:"-"`
}

type ConvertResponse struct {
	Title    string `json:"title,omitempty"`
	URL      string `json:"url_source"`
	Markdown string `json:"markdown_content"`
}

type ConvertRequest struct {
	URL string `json:"url"`
}

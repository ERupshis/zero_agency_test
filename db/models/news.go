package models

/*//go:generate reform*/

// News represents a row in news table.
//
//go:generate easyjson -all news.go
//reform:news
type News struct {
	ID         int64   `json:"Id" reform:"id,pk"`
	Title      string  `json:"Title" reform:"title"`
	Content    string  `json:"Content" reform:"content"`
	Categories []int64 `json:"Categories" reform:"-"`
}

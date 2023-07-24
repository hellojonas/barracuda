package news

type Article struct {
	Title       string
	Description string
	Date        string
	Link        string
	Image		string
}

type  NewsPage interface {
	FindNews () ([]Article, error)
}

func IsEmpty(a Article) bool {
	return a.Title == "" && a.Link == ""
}

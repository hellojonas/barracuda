package news

type Article struct {
	Title       string
	Description string
	Date        string
	Link        string
	Image       string
	Category    string
}

type NewsPage interface {
	FindNews() ([]Article, error)
}

func ValidArticle(a Article) bool {
	return a.Title != "" && a.Link != ""
}

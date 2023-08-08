package news

type Article struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Date        string `json:"date"`
	Link        string `json:"link"`
	Image       string `json:"image"`
	Category    string `json:"category"`
}

type NewsPage interface {
	FindNews() ([]Article, error)
}

func ValidArticle(a Article) bool {
	return a.Title != "" && a.Link != ""
}

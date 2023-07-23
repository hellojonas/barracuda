package barracuda

import (
	"github.com/hellojonas/barracuda/internal/sources/opais"
	"github.com/hellojonas/barracuda/pkg/news"
)

type barracuda struct {
    sources []news.NewsPage
}

func NewBarracuda() *barracuda {
    sources := []news.NewsPage{
        opais.NewPage(),
    }
    return &barracuda{
        sources: sources,
    }
}

func (b * barracuda) getArticles() ([]news.Article, error) {
    var articles []news.Article

    for _, s := range b.sources {
	as, err := s.FindNews()
	if err != nil {
	    continue
	}
	articles = append(articles, as...)
    }

    return articles, nil
}

package barracuda

import (
	"errors"
	"fmt"

	"github.com/hellojonas/barracuda/internal/sources/opais"
	"github.com/hellojonas/barracuda/pkg/news"
)

type barracuda struct {
    sources map[string]string
}

func NewBarracuda() *barracuda {
    return &barracuda{
	sources: map[string]string{
	    "opais": "O País",
	},
    }
}

// func saveArticles(souce string, article []news.Article) error {
// }

// func (b * barracuda) refreshArticles() error {
//     // refresh saved articles
// }

func (b * barracuda) getArticles(source string) ([]news.Article, error) {
    // find articles by sorouce
    _, ok := b.sources[source]

    if !ok {
	return nil, fmt.Errorf("could not locate source %s", source)
    }

    switch source {
    case "opais":
	articles, err := opais.NewPage().FindNews()
	if err != nil {
	    return nil, err
	}
	return articles, nil
    }

    return nil, errors.New("could not load articles. source not found")
}

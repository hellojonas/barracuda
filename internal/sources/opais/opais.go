package opais

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/hellojonas/barracuda/internal/berror"
	"github.com/hellojonas/barracuda/pkg/news"
)

type opais struct{}

func NewPage() news.NewsPage {
	return opais{}
}

func (o opais) FindNews() ([]news.Article, error) {
	domain := "https://opais.co.mz"
	postSelector := ".elementor-widget-container > .elementor-posts-container.elementor-posts.ecs-posts .elementor-post .elementor-section-wrap > .elementor-section .elementor-row"
	descSelector := ".elementor-column:last-child .elementor-element > .elementor-widget-container > p:first-child"
	titleSelector := ".elementor-column:last-child .elementor-heading-title a"
	imageSelector := ".elementor-column:first-child .elementor-image a img"
	dateSelector := ".elementor-column:last-child .elementor-element > .elementor-widget-container > ul.elementor-post-info.elementor-inline-items > li.elementor-inline-item:nth-child(2)"

	res, err := http.Get(domain)

	if err != nil {
		return nil, berror.New(berror.ErrSourceNotReachable, "failed to load page")
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, berror.New(berror.ErrSourceStatusNotOK, "page responded with non OK status")
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)

	if err != nil {
		return nil, berror.New(berror.ErrSourceParseFailed, fmt.Sprintf("error loading page document. %v", err))
	}

	var articles []news.Article
	doc.Find(postSelector).Each(func(_ int, s *goquery.Selection) {
		article := news.Article{}
		titleSelection := s.Find(titleSelector)
		article.Title = titleSelection.Text()
		article.Link = titleSelection.AttrOr("href", "")
		article.Description = s.Find(descSelector).Text()
		article.Image = s.Find(imageSelector).First().AttrOr("data-src", "")
		article.Date = strings.TrimSpace(s.Find(dateSelector).First().Text())

		if news.ValidArticle(article) {
			articles = append(articles, article)
		}
	})

	return articles, nil
}

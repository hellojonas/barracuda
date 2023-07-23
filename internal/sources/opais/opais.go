package opais

import (
	"fmt"
	"strings"

	"github.com/gocolly/colly"
	"github.com/hellojonas/barracuda/pkg/news"
)

type opais struct{}

func NewPage() news.NewsPage {
	return opais{}
}

func (o opais) FindNews() ([]news.Article, error) {

	c := colly.NewCollector(
		colly.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/115.0.0.0 Safari/537.36"),
	)

	articles := make([]news.Article, 0)

	c.OnHTML(".elementor-widget-container > .elementor-posts", func(e *colly.HTMLElement) {
		fmt.Println("-----> Hit")
		articles = getArticles(e)
	})

	var err error
	c.OnError(func(r *colly.Response, e error) {
		err = e
	})

	if err != nil {
		return nil, err
	}

	c.Visit("https://opais.co.mz")

	return articles, nil
}

func getArticles(h *colly.HTMLElement) []news.Article {
	var articles []news.Article

	h.ForEach(".elementor-post .elementor-row", func(i int, h *colly.HTMLElement) {
		article := news.Article{}
		h.ForEach(".elementor-column:last-child .elementor-heading-title a", func(_ int, h *colly.HTMLElement) {
			article.Title = h.Text
			article.Link = h.Attr("href")
		})
		h.ForEach(".elementor-column:last-child .elementor-element > .elementor-widget-container > p:first-child", func(_ int, h *colly.HTMLElement) {
			article.Description = h.Text
		})
		h.ForEach(".elementor-column:last-child .elementor-element > .elementor-widget-container > ul > li:nth-child(2)", func(_ int, h *colly.HTMLElement) {
			article.Date = strings.TrimSpace(h.Text)
		})
		articles = append(articles, article)
	})

	return articles
}

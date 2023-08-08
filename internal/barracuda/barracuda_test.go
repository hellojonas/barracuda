package barracuda

import (
	"database/sql"
	"log"
	"testing"
	"time"

	"github.com/hellojonas/barracuda/pkg/news"
)

func _TestLoadArticles(t *testing.T) {
    source := "opais"
    category := "news"
    date, err := time.Parse("2006-01-02", "2023-08-06")

    if err != nil {
	t.Fatalf("invalid date. %v", err)
    }

    datasource := ""
    db, err := sql.Open("postgres", datasource)

    if err != nil {
	log.Fatalf("error opening database connection. %v", err)
    }

    if err = db.Ping(); err != nil {
	log.Fatalf("could not establish connection to database %v", err)
    }

    b := New(db)

    articles, err := b.FindArticles(source, category, date)

    if err != nil {
	t.Fatalf("error finding articles. %v", err)
    }

    if len(articles) == 0 {
	t.Fatal("no articles fetched")
    }
}

func _TestSaveArticles(t *testing.T) {
    source := "opais"
    datasource := ""
    db, err := sql.Open("postgres", datasource)

    if err != nil {
	log.Fatalf("error opening database connection. %v", err)
    }

    if err = db.Ping(); err != nil {
	log.Fatalf("could not establish connection to database %v", err)
    }

    b := New(db)

    articles := []news.Article {
	{
	    Title:       "Test title",
	    Link:        "https://link-to-page",
	},
    }
    if _, err = b.saveArticles(source, "news", articles); err != nil {
	t.Fatalf("should have saved 1 article, but saved none. %v", err)
    }
}

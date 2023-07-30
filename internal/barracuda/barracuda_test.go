package barracuda

import (
	"database/sql"
	"log"
	"testing"

	"github.com/hellojonas/barracuda/pkg/news"
)

func _TestLoadArticles(t *testing.T) {
    sourceSuccess := "opais"
    sourceFail := "test"

    b := New(nil)

    _, err := b.getArticles(sourceFail)

    if err == nil {
	t.Fatal("should have failed on non exitent source")
    }

    articles, err := b.getArticles(sourceSuccess)

    if err != nil {
	t.Fatalf("failed loading articles %v\n", err)
    }

    if len(articles) == 0 {
	t.Fatal("no articles fetched")
    }
}

func TestSaveArticles(t *testing.T) {
    source := "opais"
    datasource := "user=postgres password=LvQ7qm4smtsrSQ host=db.cvvinaictvfkqxzaxkkc.supabase.co port=5432 dbname=postgres"
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
    if err = b.saveArticles(source, "news", articles); err != nil {
	t.Fatalf("should have saved 1 article, but saved none. %v", err)
    }
}

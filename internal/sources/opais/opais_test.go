package opais

import (
	"testing"
)

func TestFindNews(t *testing.T) {
	var page = NewPage()

	articles, err := page.FindNews()

	if err != nil {
		t.Fatal(err)
	}

	if len(articles) == 0 {
		t.Fatalf("expected a list of articles, got nothing.")
	}

	for _, a := range articles {
		emptyArticle := a.Title == "" && a.Description == "" && a.Link == "" && a.Date == "" 
		if emptyArticle {
			t.Fatalf("empty article found")
		}
	}
}

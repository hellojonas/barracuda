package opais

import (
	"testing"

	"github.com/hellojonas/barracuda/pkg/news"
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
		if news.IsEmpty(a) {
			t.Fatalf("empty article found")
		}
	}
}

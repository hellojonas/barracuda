package barracuda

import "testing"

func TestLoadArticles(t *testing.T) {
    sourceSuccess := "opais"
    sourceFail := "test"

    b := NewBarracuda()

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

package barracuda

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
)

func Router(b *barracuda) chi.Router {
	r := chi.NewRouter()

	r.Get("/articles/{source}", func(w http.ResponseWriter, r *http.Request) {
		params := r.URL.Query()
		category := params.Get("category")
		date := params.Get("dateStart")
		dateStart := time.Now().AddDate(0, 0, -1)
		source := chi.URLParam(r, "source")
		encoder := json.NewEncoder(w)

		if date != "" {
			d, err := time.Parse("2006-01-02", date)
			if err != nil {
				resErr := map[string]string{
					"message": "invalid date. date layout should be yyyy-mm-dd",
				}
				b.logger.Error("error fetching articles from %s", source)
				w.WriteHeader(http.StatusInternalServerError)
				_ = encoder.Encode(resErr)
				return
			}
			dateStart = d
		}

		articles, err := b.FindArticles(source, category, dateStart)

		if err != nil {
			resErr := map[string]string{
				"message": "error fetching articles",
			}
			b.logger.Error("error fetching articles from %s", source)
			w.WriteHeader(http.StatusInternalServerError)
			_ = encoder.Encode(resErr)
			return
		}

		_ = encoder.Encode(articles)
		return
	})

	return r
}

package barracuda

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func Routes() chi.Router {
    b := NewBarracuda()
    r := chi.NewRouter()

    r.Get("/articles", func(w http.ResponseWriter, r *http.Request) {
	a, err := b.getArticles()
	res := json.NewEncoder(w)

	if err != nil {
	    e := map[string]string {
		"message": "could not load articles",
	    }
	    // Set status
	    res.Encode(e)
	    return
	}

	res.Encode(a)
    })

    return r
}

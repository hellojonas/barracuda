package barracuda

import (
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"github.com/hellojonas/barracuda/internal/sources/opais"
	"github.com/hellojonas/barracuda/pkg/news"
	_ "github.com/lib/pq"
)

type barracuda struct {
	db      *sql.DB
	sources map[string]string
}

type article struct {
	Id          string         `sql:"id"`
	Title       string         `sql:"title"`
	Link        string         `sql:"link"`
	Source      string         `sql:"source"`
	Description sql.NullString `sql:"description"`
	Date        sql.NullString `sql:"date"`
	Image       sql.NullString `sql:"image"`
	Category    string         `sql:"cateogry"`
	createdAt   time.Time      `sql:"created_at"`
}

func New(db *sql.DB) *barracuda {
	return &barracuda{
		sources: map[string]string{
			"opais": "O País",
		},
		db: db,
	}
}

// func (b * barracuda) refreshArticles() error {
//     // refresh saved articles
// }

func hash(title string) (string, error) {
	sha := sha256.New()
	_, err := sha.Write([]byte(title))

	if err != nil {
		return "", fmt.Errorf("error creating hash of string. %v", err)
	}

	return hex.EncodeToString(sha.Sum(nil)), nil
}

func (b *barracuda) exists(id string) (bool, error) {
	db := b.db
	tx, err := db.Begin()

	if err != nil {
		return false, fmt.Errorf("error starting transaction")
	}

	stmt, err := tx.Prepare(`SELECT COUNT(id) > 0 FROM articles WHERE id = $1;`)

	if err != nil {
		return false, fmt.Errorf("error preparing query %v", err)
	}

	row := stmt.QueryRow(id)
	var exist bool

	if err = row.Scan(&exist); err != nil {
		return false, fmt.Errorf("error extracting sql result. %v", err)
	}

	return exist, nil
}

func (b *barracuda) saveArticles(source string, category string, articles []news.Article) error {
	db := b.db
	tx, err := db.Begin()

	if err != nil {
		return fmt.Errorf("error starting transaction")
	}

	stmt, err := tx.Prepare(`INSERT INTO articles (id, title, link, source, description, date, image, created_at, category)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9);`)

	if err != nil {
		return fmt.Errorf("error preparing query %v", err)
	}

	for _, a := range articles {
		if !news.ValidArticle(a) {
			return fmt.Errorf("failed trying to save an empty article")
		}

		hashedTitle, err := hash(a.Title)

		if err != nil {
			return err
		}

		exist, err := b.exists(hashedTitle)

		if err != nil {
			return fmt.Errorf("error executing sql query. %v", err)
		}

		if exist {
			continue
		}

		_, err = stmt.Exec(hashedTitle, a.Title, a.Link, source, a.Description, a.Date, a.Image, time.Now(), category)

		if err != nil {
			tx.Rollback()
			return fmt.Errorf("error saving an article. %v", err)
		}
	}

	return tx.Commit()
}

func (b *barracuda) getArticles(source string) ([]news.Article, error) {
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

package barracuda

import (
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/hellojonas/barracuda/internal/berror"
	"github.com/hellojonas/barracuda/internal/sources/opais"
	"github.com/hellojonas/barracuda/pkg/logs"
	"github.com/hellojonas/barracuda/pkg/news"
	_ "github.com/lib/pq"
)

type barracuda struct {
	db      *sql.DB
	logger  *logs.BLogger
	sources []bSource
}

type bSource struct {
	id       string
	category string
	page     news.NewsPage
}

type bArticle struct {
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
		db: db,
		sources: []bSource{
			{id: "opais", category: "news", page: opais.NewPage()},
		},
	}
}

func (b *barracuda) SetLogger(logger *logs.BLogger) {
	b.logger = logger
}

func (b *barracuda) RefreshArticles() {
	for _, s := range b.sources {
		b.logger.Info("fetching articles from %s", s.id)
		articles, err := s.page.FindNews()

		if _, ok := err.(berror.BError); ok {
			b.logger.Error("failed fetching articles on %s", s.id)
			continue
		}

		b.logger.Info("persisting %d articles from %s", len(articles), s.id)
		saved, err := b.saveArticles(s.id, s.category, articles)

		if err != nil {
			b.logger.Error("failed persisting articles from %s", s.id)
			continue
		}

		b.logger.Info("%d articles persisted to %s", saved, s.id)
	}
}

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
		return false, berror.New(berror.ErrDBTxInintFailed, "error starting transaction")
	}

	stmt, err := tx.Prepare(`SELECT COUNT(id) > 0 FROM articles WHERE id = $1;`)

	if err != nil {
		return false, berror.New(berror.ErrDBInvalidQuery, fmt.Sprintf("error preparing query %v", err))
	}

	row := stmt.QueryRow(id)
	var exist bool

	if err = row.Scan(&exist); err != nil {
		return false, berror.New(berror.ErrDBResultExtractFailed, fmt.Sprintf("error extracting sql result. %v", err))
	}

	return exist, nil
}

func (b *barracuda) saveArticles(source string, category string, articles []news.Article) (int, error) {
	db := b.db
	tx, err := db.Begin()

	if err != nil {
		b.logger.Error("error starting db transaction")
		return 0, berror.New(berror.ErrDBTxInintFailed, "error starting transaction")
	}

	stmt, err := tx.Prepare(`INSERT INTO articles (id, title, link, source, description, date, image, created_at, category)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9);`)

	if err != nil {
		b.logger.Error("error preparing db statement")
		return 0, berror.New(berror.ErrDBInvalidQuery, fmt.Sprintf("error preparing query %v", err))
	}

	saved := 0
	for _, a := range articles {
		if !news.ValidArticle(a) {
			continue
		}

		hashedTitle, err := hash(a.Title)

		if err != nil {
			b.logger.Error("error hasing article title, skipping")
			continue
		}

		exist, err := b.exists(hashedTitle)

		if err != nil {
			b.logger.Error("error checking if article exists, skipping")
			continue
		}

		if exist {
			b.logger.Warn("article already exists, skipping.")
			continue
		}

		_, err = stmt.Exec(hashedTitle, a.Title, a.Link, source, a.Description, a.Date, a.Image, time.Now(), category)

		if err != nil {
			b.logger.Warn("failed saving article, skipping.")
			continue
		}
		saved++
	}

	return saved, tx.Commit()
}

func (b *barracuda) FindArticles(source string, category string, dateStart time.Time) ([]news.Article, error) {
	db := b.db
	tx, err := db.Begin()

	if err != nil {
		return nil, berror.New(berror.ErrDBTxInintFailed, "error opening db transaction")
	}

	stmt, err := tx.Prepare(`SELECT title, description, link, image, date, category 
		FROM articles WHERE source = $1 and category = $2 and created_at >= $3`)

	if err != nil {
		return nil, berror.New(berror.ErrDBInvalidQuery, "error opening db transaction")
	}

	rows, err := stmt.Query(source, category, dateStart)

	if err != nil {
		return nil, berror.New(berror.ErrDBQueryFailed, "error executing query")
	}

	var articles []news.Article
	for rows.Next() {
		var title, description, link, image, date, category string
		rows.Scan(&title, &description, &link, &image, &date, &category)
		article := news.Article{
			Title:       title,
			Description: description,
			Link:        link,
			Image:       image,
			Date:        date,
			Category:    category,
		}
		articles = append(articles, article)
	}

	return articles, nil
}

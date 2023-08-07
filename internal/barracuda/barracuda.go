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

		b.logger.Info("%d articles persisted to %s", saved,  s.id)
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

func (b *barracuda) getArticles(source string) ([]news.Article, error) {
	var src news.NewsPage
	var sourceFound bool

	for _, s := range b.sources {
		if s.id == source {
			src = s.page
			sourceFound = true
		}
	}

	if !sourceFound {
		return nil, berror.New(berror.ErrSourceNotFound, fmt.Sprintf("source %s not found", source))
	}

	articles, err := src.FindNews()

	if err != nil {
		return nil, err
	}

	return articles, nil
}

package sqlite

import (
	"context"
	"database/sql"

	"github.com/Gipohub/goTgBot/lib/e"
	"github.com/Gipohub/goTgBot/storage"
	_ "github.com/glebarez/go-sqlite"
)

const (
	OpnDbErr       = "cant open database"
	CnnctDbErr     = "cant connect database"
	SavePageErr    = "cant save page"
	PickPageErr    = "cant pick random page"
	RemovePageErr  = "cant remove page"
	ExistsPageErr  = "cant check if page exists"
	CreateTableErr = "cant crate table"
)

type Storage struct {
	db *sql.DB
}

// New creates new SQLite storage.
func New(path string) (*Storage, error) {
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, e.Wrap(OpnDbErr, err)
	}

	if err := db.Ping(); err != nil {
		return nil, e.Wrap(CnnctDbErr, err)
	}

	return &Storage{db: db}, nil
}

// Save saves page to storage.
func (s *Storage) Save(ctx context.Context, p *storage.Page) error {
	q := `INSERT INTO pages (url, user_name) VALUES (?, ?)`

	_, err := s.db.ExecContext(ctx, q, p.URL, p.UserName)
	if err != nil {
		return e.Wrap(SavePageErr, err)
	}
	return nil
}

// PickRandom picks random page from storage.
func (s *Storage) PickRandom(ctx context.Context, userName string) (*storage.Page, error) {
	q := `SELECT url FROM pages WHERE user_name = ? ORDER BY RANDOM() LIMIT 1`

	var url string

	err := s.db.QueryRowContext(ctx, q, userName).Scan(&url)
	if err == sql.ErrNoRows {
		return nil, storage.ErrNoSavedPages
	}
	if err != nil {
		return nil, e.Wrap(PickPageErr, err)
	}

	return &storage.Page{
		URL:      url,
		UserName: userName,
	}, nil //craft page object, and nil error
}

// PickRandom picks pages from storage in random order.
func (s *Storage) PickAllList(ctx context.Context, userName string) ([]*storage.Page, error) {
	q := `SELECT url FROM pages WHERE user_name = ? ORDER BY RANDOM()`

	rows, err := s.db.QueryContext(ctx, q, userName)
	if err != nil {
		return nil, e.Wrap(PickPageErr, err)
	}
	defer rows.Close()

	var pages []*storage.Page

	for rows.Next() {
		var url string
		if err := rows.Scan(&url); err != nil {
			return nil, e.Wrap(PickPageErr, err)
		}
		pages = append(pages, &storage.Page{
			URL:      url,
			UserName: userName,
		})
	}

	if err := rows.Err(); err != nil {
		return nil, e.Wrap(PickPageErr, err)
	}

	if len(pages) == 0 {
		return nil, storage.ErrNoSavedPages
	}

	return pages, nil
}

// Remove removes page from storage.
func (s *Storage) Remove(ctx context.Context, page *storage.Page) error {
	q := `DELETE FROM pages WHERE url = ? AND user_name = ?`

	if _, err := s.db.ExecContext(ctx, q, page.URL, page.UserName); err != nil {
		return e.Wrap(RemovePageErr, err)
	}
	return nil
}

// IsExists checks if page exists in storage.
func (s *Storage) IsExists(ctx context.Context, page *storage.Page) (bool, error) {
	q := `SELECT COUNT(*) FROM pages WHERE url = ? AND user_name = ?`

	var count int

	err := s.db.QueryRowContext(ctx, q, page.URL, page.UserName).Scan(&count)
	if err != nil {
		return false, e.Wrap(ExistsPageErr, err)
	}
	return count > 0, nil
}

func (s *Storage) Init(ctx context.Context) error {
	q := `CREATE TABLE IF NOT EXISTS pages (url TEXT, path TEXT,user_name TEXT)`
	_, err := s.db.ExecContext(ctx, q)
	if err != nil {
		return e.Wrap(CreateTableErr, err)
	}

	return nil
}

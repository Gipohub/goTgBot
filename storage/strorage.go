package storage

import (
	"context"
	"crypto/sha1"
	"errors"
	"fmt"
	"io"

	"github.com/Gipohub/goTgBot/lib/e"
)

//import "time"

type Storage interface {
	Save(ctx context.Context, p *Page) error
	PickRandom(ctx context.Context, UserName string) (*Page, error)
	Remove(ctx context.Context, p *Page) error
	IsExists(ctx context.Context, p *Page) (bool, error)
	//Pars()
}

var ErrNoSavedPages = errors.New("no saved pages")

// основной тип хранения содержит ссылку на статью и ?юзера кому отдать?
type Page struct {
	URL      string
	UserName string
	//не стал добавлять сортировку по сначала новые (старые)
	// потом можно реализовать
	//Created  time.Time
}

func (p Page) Hash() (string, error) {
	h := sha1.New()

	if _, err := io.WriteString(h, p.URL); err != nil {
		return "", e.Wrap("cant calculate hash", err)
	}
	if _, err := io.WriteString(h, p.UserName); err != nil {
		return "", e.Wrap("cant calculate hash", err)
	}

	return fmt.Sprintf("%x", h.Sum(nil)), nil
}

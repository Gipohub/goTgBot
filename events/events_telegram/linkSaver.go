package events_telegram

import (
	"context"
	"errors"

	"github.com/Gipohub/goTgBot/lib/e"
	"github.com/Gipohub/goTgBot/storage"
)

func (p *Processor) savePage(pageURL string, chatID int, username string) (err error) {
	defer func() { err = e.Wrap("cant do command: save page", err) }()

	page := &storage.Page{
		URL:      pageURL,
		UserName: username,
	}
	isExists, err := p.storage.IsExists(context.Background(), page)
	if err != nil {
		return err
	}
	if isExists {
		return p.tg.SendMesages(chatID, msgAlreadyExists)
	}

	if err := p.storage.Save(context.Background(), page); err != nil {
		return err
	}

	if err := p.tg.SendMesages(chatID, msgSaved); err != nil {
		return err
	}
	return nil

}

func (p *Processor) SendRandom(chatID int, username string) (err error) {
	defer func() { err = e.Wrap("cant do command:cant send random", err) }()

	page, err := p.storage.PickRandom(context.Background(), username)
	if err != nil && !errors.Is(err, storage.ErrNoSavedPages) {
		return err
	}
	if errors.Is(err, storage.ErrNoSavedPages) {
		return p.tg.SendMesages(chatID, msgNoSavedPages)
	}
	if err := p.tg.SendMesages(chatID, page.URL); err != nil {
		return err
	}
	return err //p.storage.Remove(context.Background(), page)
}

func (p *Processor) SendList(chatID int, username string) (err error) {
	defer func() { err = e.Wrap("cant do command:cant send list", err) }()

	pages, err := p.storage.PickAllList(context.Background(), username)
	if err != nil && !errors.Is(err, storage.ErrNoSavedPages) {

		return err
	}
	if errors.Is(err, storage.ErrNoSavedPages) {
		return p.tg.SendMesages(chatID, msgNoSavedPages)
	}
	for _, page := range pages {
		if err := p.tg.SendMesages(chatID, page.URL); err != nil {
			return err
		}
	}
	return err //p.storage.Remove(context.Background(), page)
}

package events_telegram

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/url"
	"strings"

	"github.com/Gipohub/goTgBot/clients/ytClient"
	"github.com/Gipohub/goTgBot/lib/e"
	"github.com/Gipohub/goTgBot/storage"
)

const (
	RndCmd   = "/rnd"
	HelpCmd  = "/help"
	StartCmd = "/start"
	Parser   = "/pars"
	ListCmd  = "/list"
	Exit     = "/exit"

	//huinya = context.Background()
)

func (p *Processor) doCmd(text string, chatID int, username string) error {
	text = strings.TrimSpace(text)

	log.Printf("got new comma '%s' from '%s'", text, username)

	if isSaveCmd(text) {
		return p.savePage(text, chatID, username)
	}

	switch text {
	case RndCmd:
		fmt.Println("rnd msg")

		return p.SendRandom(chatID, username)
	case HelpCmd:
		fmt.Println("help msg")

		return p.SendHelp(chatID)
	//case SaveCmd:
	case StartCmd:
		fmt.Println("start msg")

		return p.SendHello(chatID)
	case ListCmd:
		fmt.Println("list msg")
		return p.SendList(chatID, username)
	//case Parser:
	//	fmt.Println("pars msg")
	//	return p.SendParsRes(chatID)
	case Exit:
		fmt.Println("exit msg")
		if p.isOwner(username) {
			if err := p.tg.SendMesages(chatID, msgTurnedOff); err != nil {
				return err
			}
			log.Fatal("service is stopped")
		}
		return p.tg.SendMesages(chatID, msgAccessDenied)

	default:
		fmt.Println("unknown msg")
		return p.tg.SendMesages(chatID, msgUnknownCommand)
	}

}

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

func (p *Processor) SendHelp(chatID int) error {
	return p.tg.SendMesages(chatID, msgHelp)
}
func (p *Processor) SendHello(chatID int) error {
	return p.tg.SendMesages(chatID, msgHello)
}

func (p *Processor) SendParsRes(chatID int) error {
	ytClient.Pars()
	return p.tg.SendMesages(chatID, "pars")
}

func (p *Processor) isOwner(username string) bool {
	return p.tg.GetOwner() == username
}

func isSaveCmd(text string) bool {
	return isURL(text)
}

func isURL(text string) bool {
	u, err := url.Parse(text)

	return err == nil && u.Host != ""
}

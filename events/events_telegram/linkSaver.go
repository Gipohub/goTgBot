package events_telegram

import (
	"context"
	//"errors"
	"log"

	//"time"

	"github.com/Gipohub/goTgBot/lib/e"
	//"github.com/Gipohub/goTgBot/storage"
	linksaver "github.com/Gipohub/linksaver/proto"
	//pb "github.com/Gipohub/linksaver/proto"
	//"google.golang.org/grpc"
)

func (p *Processor) savePage(pageURL string, chatID int, username string) (err error) {
	defer func() { err = e.Wrap("cant do command: save page", err) }()

	ctx := context.Background()

	existsResp, err := p.linker.IsExists(ctx, &linksaver.Page{
		Url:      pageURL,
		Username: username,
	})
	if err != nil {
		return err
	}

	if existsResp.Exists {
		return p.tg.SendMesages(chatID, msgAlreadyExists)
	}
	log.Printf("i am in savepage ")
	_, err = p.linker.Save(ctx, &linksaver.SaveRequest{
		Url:      pageURL,
		Username: username,
	})
	if err != nil {
		return err
	}

	return p.tg.SendMesages(chatID, msgSaved)

	// page := &storage.Page{
	// 	URL:      pageURL,
	// 	UserName: username,
	// }
	// isExists, err := p.storage.IsExists(context.Background(), page)
	// if err != nil {
	// 	return err
	// }
	// if isExists {
	// 	return p.tg.SendMesages(chatID, msgAlreadyExists)
	// }

	// if err := p.storage.Save(context.Background(), page); err != nil {
	// 	return err
	// }

	// if err := p.tg.SendMesages(chatID, msgSaved); err != nil {
	// 	return err
	// }
	// return nil

}

func (p *Processor) SendRandom(chatID int, username string) (err error) {
	defer func() { err = e.Wrap("cant do command: send random", err) }()

	ctx := context.Background()
	page, err := p.linker.PickRandom(ctx, &linksaver.User{
		Username: username,
	})
	if err != nil {
		return err
	}
	if page.Url == "" {
		return p.tg.SendMesages(chatID, msgNoSavedPages)
	}
	return p.tg.SendMesages(chatID, page.Url)

	// defer func() { err = e.Wrap("cant do command:cant send random", err) }()

	// page, err := p.storage.PickRandom(context.Background(), username)
	// if err != nil && !errors.Is(err, storage.ErrNoSavedPages) {
	// 	return err
	// }
	// if errors.Is(err, storage.ErrNoSavedPages) {
	// 	return p.tg.SendMesages(chatID, msgNoSavedPages)
	// }
	// if err := p.tg.SendMesages(chatID, page.URL); err != nil {
	// 	return err
	// }
	// return err //p.storage.Remove(context.Background(), page)
}

func (p *Processor) SendList(chatID int, username string) (err error) {
	defer func() { err = e.Wrap("cant do command: send list", err) }()

	ctx := context.Background()
	pageList, err := p.linker.PickAll(ctx, &linksaver.User{
		Username: username,
	})
	if err != nil {
		return err
	}
	if len(pageList.Pages) == 0 {
		return p.tg.SendMesages(chatID, msgNoSavedPages)
	}

	for _, page := range pageList.Pages {
		if err := p.tg.SendMesages(chatID, page.Url); err != nil {
			return err
		}
	}
	return nil

	// defer func() { err = e.Wrap("cant do command:cant send list", err) }()

	// pages, err := p.storage.PickAllList(context.Background(), username)
	// if err != nil && !errors.Is(err, storage.ErrNoSavedPages) {

	// 	return err
	// }
	// if errors.Is(err, storage.ErrNoSavedPages) {
	// 	return p.tg.SendMesages(chatID, msgNoSavedPages)
	// }
	// for _, page := range pages {
	// 	if err := p.tg.SendMesages(chatID, page.URL); err != nil {
	// 		return err
	// 	}
	// }
	// return err //p.storage.Remove(context.Background(), page)
}

package events_telegram

import (
	"errors"
	"fmt"

	"github.com/Gipohub/goTgBot/clients/tgClient"
	"github.com/Gipohub/goTgBot/events"
	"github.com/Gipohub/goTgBot/lib/e"
	"github.com/Gipohub/goTgBot/storage"
)

type Processor struct {
	tg      *tgClient.Client
	offset  int
	storage storage.Storage
}

type Meta struct {
	ChatID   int
	UserName string
}

var (
	ErrUnknownEventType = errors.New("unknown event type")
	ErrUnknownMetaType  = errors.New("unknown meta type")
)

func New(client *tgClient.Client, storage storage.Storage) *Processor {
	return &Processor{
		tg:      client,
		storage: storage,
	}
}

// обрабатываем апдейты от телеграма, кастим их в эвенты,
// что является общей формой событий от разных например апи
func (p *Processor) Fetch(limit int) ([]events.Event, error) {
	//получаем апдей
	updates, err := p.tg.Updates(p.offset, limit)
	if err != nil {
		return nil, e.Wrap("can't get events", err)
	}
	if len(updates) != 0 {
		fmt.Println("get some updates: ", updates, "events_telegram;p.Fetch")
	} else { //если updates нет то завершаем
		fmt.Print("`")
		return nil, nil
	}

	//алоцируем память для переменной результата
	res := make([]events.Event, 0, len(updates))
	//апдейты в эвенты
	for _, u := range updates {
		res = append(res, event(u))
	}

	p.offset = updates[len(updates)-1].ID + 1
	return res, nil
}

// разные сценарии действия в зависимости от типа эвента
func (p *Processor) Process(event events.Event) error {
	switch event.Type {
	case events.Message:
		return p.processMessage(event)
	case events.Callback:
		return p.processMessage(event)
	default:
		return e.Wrap("cant process mesage", ErrUnknownEventType)
	}
}

func (p *Processor) processMessage(event events.Event) error {
	meta, err := meta(event)
	if err != nil {
		return err
	}

	if err := p.doCmd(event.Text, meta.ChatID, meta.UserName); err != nil {
		return e.Wrap("cnt prcss mssage", err)
	}

	return nil
}

func meta(event events.Event) (Meta, error) {
	res, ok := event.Meta.(Meta)
	if !ok {
		return Meta{}, e.Wrap("cnt get meta", ErrUnknownMetaType)
	}
	return res, nil
}

func event(upd tgClient.Update) events.Event {
	updType := fetchType(upd)

	res := events.Event{
		Type: updType,
		Text: fetchText(updType, upd),
	}

	if updType == events.Message {
		res.Meta = Meta{
			ChatID:   upd.Message.Chat.ID,
			UserName: upd.Message.From.Username,
		}
	}
	if updType == events.Callback {
		res.Meta = Meta{
			ChatID:   upd.Callback.Message.Chat.ID,
			UserName: upd.Callback.From.Username,
		}
		fmt.Println(res)
	}
	return res
}
func fetchText(updType events.Type, upd tgClient.Update) string {
	switch updType {
	//case 0: return ""
	case 1:
		return upd.Message.Text
	case 2:
		return upd.Callback.Data
	default:
		return ""
	}
}
func fetchType(upd tgClient.Update) events.Type {

	if upd.Message == nil {
		if upd.Callback == nil {
			return events.Unknown
		}
		return events.Callback
	}
	return events.Message
}

package events_telegram

import (
	"fmt"

	"github.com/Gipohub/goTgBot/clients/tgClient"
	"github.com/Gipohub/goTgBot/events"
	"github.com/Gipohub/goTgBot/lib/e"
)

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
		fmt.Print("`") //показатель работы цикла
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

func event(upd tgClient.Update) events.Event {
	updType := fetchType(upd)

	res := events.Event{
		Type: updType,
		Text: fetchText(updType, upd),
	}

	if updType == events.Message {
		res.Meta = events.Meta{
			ChatID:   upd.Message.Chat.ID,
			UserName: upd.Message.From.Username,
		}
	}
	if updType == events.Callback {
		res.Meta = events.Meta{
			ChatID:   upd.Callback.Message.Chat.ID,
			UserName: upd.Callback.From.Username,
		}
		fmt.Println(res)
	}
	return res
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

func fetchText(updType events.Type, upd tgClient.Update) string {
	switch updType {
	case 1:
		return upd.Message.Text
	case 2:
		return upd.Callback.Data
	default:
		return ""
	}
}

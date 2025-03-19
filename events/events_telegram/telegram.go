package events_telegram

import (
	"context"
	"errors"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/Gipohub/goTgBot/clients/tgClient"
	"github.com/Gipohub/goTgBot/events"
	"github.com/Gipohub/goTgBot/lib/e"
	"github.com/Gipohub/goTgBot/storage"
)

type Processor struct {
	tg      *tgClient.Client
	offset  int
	storage storage.Storage
	monitor map[string]events.RoutineData
	mu      sync.Mutex
}

// type Meta struct {
// 	ChatID   int
// 	UserName string
// }

var (
	ErrUnknownEventType = errors.New("unknown event type")
	ErrUnknownMetaType  = errors.New("unknown meta type")
)

func New(client *tgClient.Client, storage storage.Storage) *Processor {
	// d := make(map[string]events.RoutineData)
	// c := make(chan events.Event)
	// ctx, cancel:= context.WithCancel(context.Background())
	// d["sem"] =  events.RoutineData{channel: c,
	// context: &ctx}
	p := &Processor{
		tg:      client,
		storage: storage,
		monitor: make(map[string]events.RoutineData),
		//data: d,
	}
	//go p.Start(p)

	return p

}

// func (p *Processor) Start(processor *Processor){
// 	p.data.context.
// }

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

// разные сценарии действия в зависимости от типа эвента
func (p *Processor) Process(event events.Event) error {
	switch event.Type {
	case events.Message, events.Callback:
		return p.processMessage(event)
	default:
		return e.Wrap("cant process mesage", ErrUnknownEventType)
	}
}

func (p *Processor) processMessage(event events.Event) error {
	// meta, err := meta(event)
	// if err != nil {
	// 	return err
	// }
	p.mu.Lock()
	rData, exists := p.monitor[event.Meta.UserName]
	if !exists {
		ch := make(chan events.Event, 10)
		initState := events.RoutineData{
			Channel: ch,
			Context: context.TODO(), ////TODO manage ctx with db request
		}
		p.monitor[event.Meta.UserName] = initState
		go p.userState(initState)
		ch <- event
	} else {
		rData.Channel <- event
	}
	p.mu.Unlock()
	return nil
}
func (p *Processor) userState(data events.RoutineData) {
	timer := time.NewTimer(5 * time.Minute) // Запускаем таймер
	var userName string
	defer func() {
		timer.Stop()
		if len(userName) > 0 {
			p.mu.Lock()
			delete(p.monitor, userName) // Удаляем юзера из монитора
			p.mu.Unlock()
		}
		close(data.Channel)
	}()
	i := 0
	for {
		i++
		fmt.Println(i)
		select {
		case nE := <-data.Channel:
			timer.Reset(5 * time.Minute) // Продлеваем на 5 минут
			userName = nE.Meta.UserName  //для удаления юзера из мапы в defer()

			if err := p.doCmd(nE.Text, nE.Meta.ChatID, nE.Meta.UserName); err != nil {
				log.Printf("cnt prcss mssage in processMessage: %s", err)
			}

		case <-data.Context.Done():
			log.Println("Stopping user process goroutine")
			return // Выход из горутины
		case <-timer.C: // Если таймер истек — завершаем контекст
			log.Println("Timeout reached, stopping process")
			return
		}
	}
}

// func meta(event events.Event) (events.Meta, error) {
// 	res, ok := event.Meta.(events.Meta)
// 	if !ok {
// 		return Meta{}, e.Wrap("cnt get meta", ErrUnknownMetaType)
// 	}
// 	return res, nil
// }

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
	//case 0: return ""
	case 1:
		return upd.Message.Text
	case 2:
		return upd.Callback.Data
	default:
		return ""
	}
}

package events_telegram

import (
	"context"
	"log"
	"time"

	"github.com/Gipohub/goTgBot/events"
	"github.com/Gipohub/goTgBot/lib/e"
)

// разные сценарии действия в зависимости от типа эвента
func (p *Processor) Process(event events.Event) error {
	switch event.Type {
	case events.Message, events.Callback:
		return p.processMessage(event)
	default:
		return e.Wrap("cant process mesage", ErrUnknownEventType)
	}
}

// check or create active user session and them do message process
func (p *Processor) processMessage(event events.Event) error {
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
		log.Printf("%v'nd event on session", i)
		select {
		case nE := <-data.Channel:
			// Продлеваем на 5 минут
			timer.Reset(5 * time.Minute)
			//для удаления юзера из мапы в defer()
			userName = nE.Meta.UserName

			if err := p.doCmd(nE.Text, nE.Meta.ChatID, nE.Meta.UserName); err != nil {
				log.Printf("cnt doCmd mssage in userState: %s", err)
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

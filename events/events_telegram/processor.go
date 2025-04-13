package events_telegram

import (
	"errors"
	"fmt"
	"log"
	"strings"
	"sync"

	"github.com/Gipohub/goTgBot/clients/tgClient"
	"github.com/Gipohub/goTgBot/events"
	"github.com/Gipohub/goTgBot/storage"
)

type Processor struct {
	tg                *tgClient.Client
	offset            int
	storage           storage.Storage
	activeUserSession map[string]events.RoutineData
	commandList       map[string]func(text string, chatID int, username string) error
	mu                sync.Mutex
}

var (
	ErrUnknownEventType = errors.New("unknown event type")
	ErrUnknownMetaType  = errors.New("unknown meta type")
)

func New(client *tgClient.Client, storage storage.Storage) *Processor {
	p := &Processor{
		tg:                client,
		storage:           storage,
		activeUserSession: make(map[string]events.RoutineData),
	}

	return p

}

const (
	RndCmd   = "/rnd"
	HelpCmd  = "/help"
	StartCmd = "/start"
	//Parser   = "/pars"
	ListCmd = "/list"
	Exit    = "/exit"
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

	case StartCmd:
		fmt.Println("start msg")

		return p.SendHello(chatID)

	case ListCmd:
		fmt.Println("list msg")
		return p.SendList(chatID, username)

	case Exit:
		fmt.Println("exit msg")
		if p.isOwner(username) {
			if err := p.tg.SendMesages(chatID, msgTurnedOff); err != nil {
				return err
			}
			//turned off
			log.Fatal("service is stopped")
		}
		return p.tg.SendMesages(chatID, msgAccessDenied)

	default:
		fmt.Println("unknown msg")
		return p.tg.SendMesages(chatID, msgUnknownCommand)
	}

}

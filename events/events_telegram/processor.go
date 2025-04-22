package events_telegram

import (
	"errors"
	"fmt"
	"log"
	"strings"
	"sync"

	"github.com/Gipohub/goTgBot/clients/tgClient"
	"github.com/Gipohub/goTgBot/events"
	"github.com/Gipohub/goTgBot/lib/grpcCommand"
	"github.com/Gipohub/goTgBot/storage"

	linksaver "github.com/Gipohub/linksaver/proto"
)

type Processor struct {
	tg                *tgClient.Client
	offset            int
	storage           storage.Storage
	activeUserSession map[string]events.RoutineData
	commandList       map[string]func(text string, chatID int, username string) error
	mu                sync.Mutex
	linker            linksaver.LinkSaverClient
}

var (
	ErrUnknownEventType = errors.New("unknown event type")
	ErrUnknownMetaType  = errors.New("unknown meta type")
)

const connectionAddress = "localhost:50051"

func New(client *tgClient.Client, storage storage.Storage) (*Processor, error) {

	conn, err := grpcCommand.ConnectToServer(connectionAddress)
	if err != nil {
		return nil, err
	}

	rpcClient := linksaver.NewLinkSaverClient(conn)

	p := &Processor{
		tg:                client,
		storage:           storage,
		activeUserSession: make(map[string]events.RoutineData),
		linker:            rpcClient,
	}

	return p, nil

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
		log.Printf("i am in issavecmd ")
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

		if com, exists := p.commandList[text]; !exists {
			fmt.Println("unknown msg")
			return p.tg.SendMesages(chatID, msgUnknownCommand)
		} else {
			if err := com(text, chatID, username); err != nil {
				log.Printf("cnt do message in doCmd: %s", err)
				return err
			}
			return nil
		}

	}

}

package main

import (
	"context"
	"flag"
	"log"

	//"context"
	"github.com/Gipohub/goTgBot/clients/tgClient"
	"github.com/Gipohub/goTgBot/consumer/event_consumer"
	"github.com/Gipohub/goTgBot/events/events_telegram"
	"github.com/Gipohub/goTgBot/storage/sqlite"
)

const (
	tgBotHost = "api.telegram.org"
	//storagePath = "storage"
	sqliteStoragePath = `data/sqlite/base.db`
	batchSize         = 100
)

func main() {

	//tgClient := tgClient.New(tgBotHost, mustToken())
	s, err := sqlite.New(sqliteStoragePath)
	if err != nil {
		log.Fatal("cant connect storage", err)
	}

	err = s.Init(context.TODO())
	if err != nil {
		log.Fatal("can't init storage: ", err)
	}

	eventsProcessor := events_telegram.New(
		tgClient.New(tgBotHost, mustToken()),
		s,
		//files.New(storagePath),
	)

	log.Print("service started")

	consumer := event_consumer.New(eventsProcessor, eventsProcessor, batchSize)

	if err := consumer.Start(); err != nil {
		log.Fatal("service is stopped", err)
	}
}

// приставка маст говорит о том,
// что функция должна завершаться успешно и,
// если это не так, вызывается падение
func mustToken() string {
	token := flag.String(
		"tg-bot-token",
		"",
		"token for access to telegram bot",
	)

	flag.Parse()

	if *token == "" {
		log.Fatal("token is not specified")
	}

	return *token
}

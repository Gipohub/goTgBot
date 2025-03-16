package main

import (
	"context"
	"flag"
	"log"
	"path/filepath"

	//"context"
	"github.com/Gipohub/goTgBot/clients/tgClient"
	"github.com/Gipohub/goTgBot/consumer/event_consumer"
	"github.com/Gipohub/goTgBot/events/events_telegram"
	"github.com/Gipohub/goTgBot/storage/sqlite"
)

const (
	tgBotHost = "api.telegram.org"
	batchSize = 100
)
const (
	dataBasePathPart1 = "data"
	dataBasePathPart2 = "sqlite"
	filename          = "base.db"
)

func main() {
	sqliteStoragePath := filepath.Join(dataBasePathPart1, dataBasePathPart2, filename)

	sqliteStorage, err := sqlite.New(sqliteStoragePath)
	if err != nil {
		log.Fatal("cant connect storage", err)
	}

	err = sqliteStorage.Init(context.TODO())
	if err != nil {
		log.Fatal("can't init storage: ", err)
	}

	token, owner := mustToken()

	eventsProcessor := events_telegram.New(
		tgClient.New(tgBotHost, token, owner),
		sqliteStorage,
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
func mustToken() (string, string) {
	token := flag.String(
		"tg-bot-token",
		"",
		"token for access to telegram bot",
	)
	owner := flag.String(
		"tg-bot-owner",
		"",
		"owner name for access to bot owner commands",
	)

	flag.Parse()

	if *token == "" {
		log.Fatal("token is not specified")
	}

	if *owner == "" {
		log.Fatal("owner is not specified")
	}

	return *token, *owner
}

package events_telegram

import (
	"net/url"

	"github.com/Gipohub/goTgBot/clients/tgClient"
)

func (p *Processor) SendHello(chatID int) error {
	return p.tg.SendMesages(chatID, msgHello)
}

func (p *Processor) SendHelp(chatID int) error {

	res := []tgClient.InlineButton{{Text: "send1 random", CallbackData: "/rnd"},
		{Text: "send list", CallbackData: "/list"},
		{Text: "send2 random", CallbackData: "/rnd"},
		{Text: "send3 random", CallbackData: "/rnd"},
		{Text: "send4 random", CallbackData: "/rnd"},
		{Text: "send5 random", CallbackData: "/rnd"},
		{Text: "send6 random", CallbackData: "/rnd"}}

	return p.tg.SendButtons(chatID, msgHelp, res, []int{1, 2, 1, 3})
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

package sharebot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type messageHandler interface {
	handle(bot *Bot, update tgbotapi.Update) request
}

type defaultMessageHandler struct {
	handler func(bot *Bot, update tgbotapi.Update) request
}

func (d defaultMessageHandler) handle(bot *Bot, update tgbotapi.Update) request {
	return d.handler(bot, update)
}

package sharebot

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Bot struct {
	tgApi *tgbotapi.BotAPI
	db    *ShopsDataBase

	maxinlineheight int
	maxinlinewidth  int
}

func NewBot(botContext Context) Bot {
	api, err := tgbotapi.NewBotAPI(botContext.TGToken)

	if err != nil {
		log.Fatal(err)
	}

	return Bot{
		tgApi:           api,
		db:              &botContext.DataBase,
		maxinlineheight: 1,
		maxinlinewidth:  1,
	}
}

func (bot *Bot) SetInlineKeyboardSize(height, width int) {
	bot.maxinlineheight = height
	bot.maxinlinewidth = width
}

func (bot *Bot) Run() {
	log.Printf("Authorized on account %s", bot.tgApi.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := bot.tgApi.GetUpdatesChan(u)

	for update := range updates {
		var handler func(bot *Bot, update tgbotapi.Update) request
		if update.Message != nil {
			log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
			handler = getNewMessageHandler(update)
		} else if update.CallbackQuery != nil {
			log.Printf("[%s] %s", update.CallbackQuery.From.UserName, update.CallbackQuery.Data)
			handler = getCallbckQueryHandler(update)
		}

		res := handler(bot, update)
		res.send()
	}
}

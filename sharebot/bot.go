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

	userTable UserTable
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
		userTable:       NewUserTable(),
	}
}

func (bot *Bot) SetInlineKeyboardSize(height, width int) {
	bot.maxinlineheight = height
	bot.maxinlinewidth = width
}

func updateLoger(update tgbotapi.Update) {
	if update.Message != nil {
		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
	} else if update.CallbackQuery != nil {
		log.Printf("[%s] %s", update.CallbackQuery.From.UserName, update.CallbackQuery.Data)
	}
}

func getUsername(update tgbotapi.Update) string {
	if update.Message != nil {
		return update.Message.From.UserName
	} else if update.CallbackQuery != nil {
		return update.CallbackQuery.From.UserName
	} else {
		log.Panic("getUsername() went wrong")
		panic("getUsername() went wrong")
	}
}

func (bot *Bot) Run() {
	log.Printf("Authorized on account %s", bot.tgApi.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := bot.tgApi.GetUpdatesChan(u)

	for update := range updates {
		session := bot.userTable.GetSession(getUsername(update))
	
		updateLoger(update)
		
		handler := session.GetHandler(update) // методы для session
		handler.handle(bot, update).send()
	}
}

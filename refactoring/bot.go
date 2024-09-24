package sharebot

import (
	// "log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Bot struct {
	tgApi *tgbotapi.BotAPI
	users *UserDataBase
	// logger
}

func NewBot(context Context) (Bot, error) {
	api, err := tgbotapi.NewBotAPI(context.TGToken)

	if err != nil {
		return Bot{}, err
	}

	return Bot{
			tgApi: api,
			users: NewUserDataBase(func(id *tgbotapi.User) Session { return Session{} }),
		},
		nil
}

func (bot *Bot) SetInlineKeyboardSize(height, width uint) {}

func (bot *Bot) Run() {
	go bot.userProcess()
}

func (bot *Bot) userProcess() {
	// start
	updatecfg := tgbotapi.NewUpdate(0)
	updatecfg.Timeout = 60
	updates := bot.tgApi.GetUpdatesChan(updatecfg)
	// log.Printf("Authorized on account %s", bot.tgApi.Self.UserName)

	for update := range updates {
		//logUpdate(update)
		bot.stateProcess(update)
	}
}

func (bot *Bot) stateProcess(update tgbotapi.Update) {
	user := getUserFromUpdate(update)

	session := bot.users.СaptureSession(user)
	defer bot.users.ReleaseSession(user)

	session.Handle(bot, update)
}

func getUserFromUpdate(update tgbotapi.Update) *tgbotapi.User {
	if update.Message != nil {
		return update.Message.From
	} else if update.CallbackQuery != nil {
		return update.CallbackQuery.From
	} else {
		// no panic
		panic("getUsername() went wrong")
	}
}


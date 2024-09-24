package sharebot

import (
	"logger"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// Bot struct
type Bot struct {
	tgAPI *tgbotapi.BotAPI

	updates tgbotapi.UpdatesChannel

	users *userDataBase

	logAPI  *logger.Logger
	context Context
}

// NewBot construct new bot from context
func NewBot(context Context, logAPI *logger.Logger) (Bot, error) {
	api, err := tgbotapi.NewBotAPI(context.TGToken)

	if err != nil {
		return Bot{}, err
	}

	return Bot{
			tgAPI:   api,
			users:   newUserDataBase(func(id *tgbotapi.User) session { return session{} }),
			context: context,
			logAPI:  logAPI,
		},
		nil
}

// SetInlineKeyboardSize sets inline keyboard setting
func (bot *Bot) SetInlineKeyboardSize(height, width uint) {}

// Run starts workflow execution
// TODO : async save
func (bot *Bot) Run() {
	bot.authorize()
	bot.load()
	bot.mainPipline()
}

func (bot *Bot) load() {}

func (bot *Bot) authorize() {
	updatecfg := tgbotapi.NewUpdate(0)
	updatecfg.Timeout = 60
	bot.updates = bot.tgAPI.GetUpdatesChan(updatecfg)

	bot.logAPI.Infof("Authorized on account %s", bot.tgAPI.Self.UserName)
}

func (bot *Bot) mainPipline() {
	for update := range bot.updates {
		bot.logAPI.Info(update)

		bot.process(update)
	}
}

func (bot *Bot) process(update tgbotapi.Update) {
	user := bot.getUserFromUpdate(update)

	session := bot.users.СaptureSession(user)
	session.Handle(bot, update)
	bot.users.ReleaseSession(user)
}

func (bot *Bot) getUserFromUpdate(update tgbotapi.Update) *tgbotapi.User {
	if update.Message != nil {
		return update.Message.From
	} else if update.CallbackQuery != nil {
		return update.CallbackQuery.From
	} else {
		bot.logAPI.Panic("can't see user in update")
		return nil
	}
}

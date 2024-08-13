package sharebot

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

type UserState interface {
	NextState(update tgbotapi.Update) (messageHandler, UserState)
}

type NotAuthorisedState struct{}

type DefaultState struct{}

type ErrorState struct{}

func (NotAuthorisedState) NextState(update tgbotapi.Update) (messageHandler, UserState) {
	if update.Message != nil {
		return authGetNewMessageHandler(update)
	} else if update.CallbackQuery != nil {
		return authGetCallbackQueryHandler(update)
	} else {
		panic("incorrect message")
	}
}

func (DefaultState) NextState(update tgbotapi.Update) (messageHandler, UserState) {
	if update.Message != nil {
		return defaultGetNewMessageHandler(update)
	} else if update.CallbackQuery != nil {
		return defaultGetCallbackQueryHandler(update)
	} else {
		panic("incorrect message")
	}
}

func (ErrorState) NextState(update tgbotapi.Update) (messageHandler, UserState) {
	panic("Error state")
}

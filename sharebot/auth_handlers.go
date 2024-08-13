package sharebot

import (
	"encoding/json"
	"math/rand"
	"reflect"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func authGetNewMessageHandler(update tgbotapi.Update) (messageHandler, UserState) {
	if reflect.TypeOf(update.Message.Text).Kind() == reflect.String && update.Message.Text != "" {
		return defaultMessageHandler{authStartHandler}, NotAuthorisedState{}
	}

	return defaultMessageHandler{authStartHandler}, NotAuthorisedState{}
}

func authGetCallbackQueryHandler(update tgbotapi.Update) (messageHandler, UserState) {
	var input inlineKeyboardData
	json.NewDecoder(strings.NewReader(update.CallbackQuery.Data)).Decode(&input)

	switch input.Type {
	case kbdTypeAuthRequest:
		if rand.Int()%2 == 0 {
			return defaultMessageHandler{notSuccesfulAuthHandler}, NotAuthorisedState{}
		} else {
			return defaultMessageHandler{succesfulAuthHandler}, DefaultState{}
		}
	default:
		return defaultMessageHandler{defaultCallbackHandler}, NotAuthorisedState{}
	}
}

var authStartHandler = func(bot *Bot, update tgbotapi.Update) request {
	reply := "Для старта работы авторизуйтесь"

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, reply)

	data := inlineKeyboardData{
		Type: kbdTypeAuthRequest,
		Data: "",
	}
	databuilder := strings.Builder{}
	json.NewEncoder(&databuilder).Encode(data)

	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Авторизоваться", databuilder.String()),
		),
	)

	return sendRequest{
		msg:          &msg,
		upd:          &update,
		removeUpd:    false,
		showKeyboard: false,
		api:          bot.tgApi,
	}
}

var notSuccesfulAuthHandler = func(bot *Bot, update tgbotapi.Update) request {
	reply := "Ошибка при попытке авторизации"

	msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, reply)

	data := inlineKeyboardData{
		Type: "AuthRequest",
		Data: "",
	}
	databuilder := strings.Builder{}
	json.NewEncoder(&databuilder).Encode(data)

	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Авторизоваться", databuilder.String()),
		),
	)

	return sendRequest{
		msg:          &msg,
		upd:          &update,
		removeUpd:    false,
		showKeyboard: false,
		api:          bot.tgApi,
	}
}

var succesfulAuthHandler = func(bot *Bot, update tgbotapi.Update) request {
	reply := "Вот какие магазины доступны на данный момент"

	msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, reply)
	msg.ReplyMarkup = shopListKeybord(bot, 0)

	return sendRequest{
		msg:          &msg,
		upd:          &update,
		removeUpd:    false,
		showKeyboard: false,
		api:          bot.tgApi,
	}
}

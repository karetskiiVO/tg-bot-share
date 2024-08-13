package sharebot

import (
	"encoding/json"
	"log"
	"reflect"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func defaultGetNewMessageHandler(update tgbotapi.Update) (messageHandler, UserState) {
	if reflect.TypeOf(update.Message.Text).Kind() == reflect.String && update.Message.Text != "" {
		if description, ok := commands[update.Message.Text]; ok {
			return defaultMessageHandler{description.Handler}, DefaultState{}
		} else {
			return defaultMessageHandler{defaultHandler}, DefaultState{}
		}
	}

	return defaultMessageHandler{errorHandler}, DefaultState{}
}

func defaultGetCallbackQueryHandler(update tgbotapi.Update) (messageHandler, UserState) {
	var input inlineKeyboardData
	json.NewDecoder(strings.NewReader(update.CallbackQuery.Data)).Decode(&input)

	switch input.Type {
	case kbdTypeShopRequest:
		return defaultMessageHandler{shopCallbackHandler}, DefaultState{}
	case kbdTypeScrollRequest:
		return defaultMessageHandler{scrollCallbackHandler}, DefaultState{}
	default:
		return defaultMessageHandler{defaultCallbackHandler}, DefaultState{}
	}
}

var startHandler = func(bot *Bot, update tgbotapi.Update) request {
	reply := "Этот бот создан для cardshare, для получения карты, которая есть в наличии просто напиши название магазина."

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, reply)

	return sendRequest{
		msg:          &msg,
		upd:          &update,
		removeUpd:    true,
		showKeyboard: true,
		api:          bot.tgApi,
	}
}

var errorHandler = func(bot *Bot, update tgbotapi.Update) request {
	reply := "К сожалению данный формат не поддерживается"

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, reply)

	return sendRequest{
		msg:          &msg,
		upd:          &update,
		removeUpd:    false,
		showKeyboard: true,
		api:          bot.tgApi,
	}
}

var defaultHandler = func(bot *Bot, update tgbotapi.Update) request {
	return getRequestFromShopName(bot, update.Message.From.ID, update.Message.Text)
}

var shopsHandler = func(bot *Bot, update tgbotapi.Update) request {
	reply := "Вот какие магазины доступны на данный момент"

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, reply)
	msg.ReplyMarkup = shopListKeybord(bot, 0)

	return sendRequest{
		msg:          &msg,
		upd:          &update,
		removeUpd:    true,
		showKeyboard: false,
		api:          bot.tgApi,
	}
}

var defaultCallbackHandler = func(bot *Bot, update tgbotapi.Update) request {
	reply := update.CallbackQuery.Data

	msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, reply)

	return sendRequest{
		msg:          &msg,
		upd:          &update,
		removeUpd:    false,
		showKeyboard: false,
		api:          bot.tgApi,
	}
}

var scrollCallbackHandler = func(bot *Bot, update tgbotapi.Update) request {
	var input inlineKeyboardData
	json.NewDecoder(strings.NewReader(update.CallbackQuery.Data)).Decode(&input)
	startPos, err := strconv.Atoi(input.Data)

	if err != nil {
		log.Panic(err)
	}

	return keyboardUpdateRequest{
		replyMarkup: tgbotapi.NewEditMessageReplyMarkup(
			update.CallbackQuery.Message.Chat.ID,
			update.CallbackQuery.Message.MessageID,
			shopListKeybord(bot, startPos),
		),
		api: bot.tgApi,
	}
}

var shopCallbackHandler = func(bot *Bot, update tgbotapi.Update) request {
	var input inlineKeyboardData
	json.NewDecoder(strings.NewReader(update.CallbackQuery.Data)).Decode(&input)

	return getRequestFromShopName(bot, update.CallbackQuery.From.ID, input.Data)
}


package sharebot

import (
	"encoding/json"
	"log"
	"reflect"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func getNewMessageHandler(update tgbotapi.Update) func(bot *Bot, update tgbotapi.Update) request {
	// authorisation check

	if reflect.TypeOf(update.Message.Text).Kind() == reflect.String && update.Message.Text != "" {
		if description, ok := commands[update.Message.Text]; ok {
			return description.Handler
		} else {
			return defaultHandler
		}
	}

	return errorHandler
}

// сделать типа enum
func getCallbckQueryHandler(update tgbotapi.Update) func(bot *Bot, update tgbotapi.Update) request {
	var input inlineKeyboardData
	json.NewDecoder(strings.NewReader(update.CallbackQuery.Data)).Decode(&input)

	switch input.Type {
	case "ShopRequest":
		return shopCallbackHandler
	case "ScrollRequest":
		return scrollCallbackHandler
	default:
		return defaultCallbackHandler
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

type inlineKeyboardData struct {
	Type string `json:"type"`
	Data string `json:"data"`
}

func newInlineButtonFromShop(shopname string) tgbotapi.InlineKeyboardButton {
	data := inlineKeyboardData{
		Type: "ShopRequest",
		Data: shopname,
	}

	databuilder := strings.Builder{}
	json.NewEncoder(&databuilder).Encode(data)

	return tgbotapi.NewInlineKeyboardButtonData(shopname, databuilder.String())
}

func newInlineScrollButton(text string, startPos int) tgbotapi.InlineKeyboardButton {
	data := inlineKeyboardData{
		Type: "ScrollRequest",
		Data: strconv.Itoa(startPos),
	}

	databuilder := strings.Builder{}
	json.NewEncoder(&databuilder).Encode(data)

	return tgbotapi.NewInlineKeyboardButtonData(text, databuilder.String())
}

func shopListKeybord(bot *Bot, startPos int) tgbotapi.InlineKeyboardMarkup {
	startPos = max(startPos, 0)

	keyboardRows := make([][]tgbotapi.InlineKeyboardButton, 0)

	// кнопки магазинов непосредственно
	height, width := bot.maxinlineheight, bot.maxinlinewidth
	for h := 0; h < height; h++ {
		row := make([]tgbotapi.InlineKeyboardButton, 0, width)

		for w := 0; w < width && h*width+w+startPos < len(bot.db.ShopList); w++ {
			row = append(row, newInlineButtonFromShop(bot.db.ShopList[h*width+w+startPos]))
		}

		if len(row) == 0 {
			break
		}

		keyboardRows = append(keyboardRows, tgbotapi.NewInlineKeyboardRow(row...))
	}

	scrollButtons := make([]tgbotapi.InlineKeyboardButton, 0)
	if startPos > 0 {
		newStartPos := max(0, startPos-height*width)
		scrollButtons = append(scrollButtons, newInlineScrollButton("<", newStartPos))
	}
	if height*width+startPos < len(bot.db.ShopList) {
		newStartPos := startPos + height*width
		scrollButtons = append(scrollButtons, newInlineScrollButton(">", newStartPos))
	}
	if len(scrollButtons) > 0 {
		keyboardRows = append(keyboardRows, tgbotapi.NewInlineKeyboardRow(scrollButtons...))
	}

	return tgbotapi.NewInlineKeyboardMarkup(keyboardRows...)
}

func getRequestFromShopName(bot *Bot, chatId int64, shopname string) request {
	photo, err := bot.db.FindShopCard(shopname)

	if err == nil {
		return photoSendUpdateRequest{
			photo: tgbotapi.NewPhoto(
				chatId,
				tgbotapi.FileBytes{Name: "picture", Bytes: photo},
			),
			api: bot.tgApi,
		}
	} else {
		reply := "Данный магазин отсутствует, но мы работаем над этим"

		msg := tgbotapi.NewMessage(chatId, reply)

		return sendRequest{
			msg:          &msg,
			upd:          nil,
			removeUpd:    false,
			showKeyboard: true,
			api:          bot.tgApi,
		}
	}
}

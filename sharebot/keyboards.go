package sharebot

import (
	"encoding/json"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type inlineKeyboardData struct {
	Type string `json:"type"`
	Data string `json:"data"`
}

const (
	kbdTypeAuthRequest   = "AuthRequest"
	kbdTypeShopRequest   = "ShopRequest"
	kbdTypeScrollRequest = "ScrollRequest"
)

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

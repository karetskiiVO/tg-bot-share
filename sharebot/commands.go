package sharebot

import (
	"sort"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type commandDescriptor struct {
	Commandname string
	Description string
	Handler     func(bot *Bot, update tgbotapi.Update) request
	Public      bool
}

var commands = map[string]commandDescriptor{
	"/start": {
		Handler: startHandler,
		Public:  false,
	},
	"Начало": {
		Commandname: "Начало",
		Handler:     startHandler,
		Public:      true,
	},
	"Магазины": {
		Commandname: "Магазины",
		Handler:     shopsHandler,
		Public:      true,
	},
	// "Поделиться картой": {
	// 	Commandname: "Поделиться картой",
	// 	Handler:     startHandler,
	// 	Public:      true,
	// },
}

var commandKeyboard = tgbotapi.NewReplyKeyboard(
	func() [][]tgbotapi.KeyboardButton {
		res := make([][]tgbotapi.KeyboardButton, 0)
		commandnames := make([]string, 0)
		for commandname, description := range commands {
			if description.Public {
				commandnames = append(commandnames, commandname)
			}
		}

		sort.Strings(commandnames)

		for _, commandname := range commandnames {
			res = append(res, tgbotapi.NewKeyboardButtonRow(tgbotapi.NewKeyboardButton(commandname)))
		}
		return res
	}()...,
)

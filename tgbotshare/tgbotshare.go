package tgbotshare

import (
	"bufio"
	"encoding/json"
	"log"
	"os"
	"reflect"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type TGbot struct {
	authorisedUsers map[string]struct{}
	allowedCards    map[string]string

	databasePath string
	tgApi        *tgbotapi.BotAPI
}

func (bot *TGbot) init() *TGbot {
	bot.authorisedUsers = make(map[string]struct{})
	bot.allowedCards = make(map[string]string)

	var buf []struct {
		Shop     string `json:"shopname"`
		CardPath string `json:"cardpath"`
	}

	file, err := os.Open(bot.databasePath + "/content.json")
	if err != nil {
		panic("cannot open " + bot.databasePath + "/content.json")
	}
	defer file.Close()

	dec := json.NewDecoder(bufio.NewReader(file))

	for dec.More() {
		err := dec.Decode(&buf)
		if err != nil {
			panic(err)
		}

		for _, elem := range buf {
			bot.allowedCards[elem.Shop] = elem.CardPath
		}
	}

	return bot
}

func NewTGbot(token, databasepath string) *TGbot {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		panic(err)
	}

	return (&TGbot{tgApi: bot, databasePath: databasepath}).init()
}

func (bot *TGbot) Run() {
	log.Printf("Authorized on account %s", bot.tgApi.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := bot.tgApi.GetUpdatesChan(u)

	var commands map[string]struct{
		handler     func(update tgbotapi.Update) tgbotapi.MessageConfig
		description string
	}
	commands = map[string]struct {
		handler     func(update tgbotapi.Update) tgbotapi.MessageConfig
		description string
	}{
		"/start": {
			handler: func(update tgbotapi.Update) tgbotapi.MessageConfig {
				return tgbotapi.NewMessage(update.Message.Chat.ID,
					"Этот бот создан для cardshare, для получения карты, "+
						"которая есть в наличии просто напиши название магазина.\n"+
						"Для повторения этого сообщения напиши /start\n"+
						"Для ознокомления с полным списком команд /help")
			},
			description: "вывод страртового сообщения"},
		"/help": {
			handler: func(update tgbotapi.Update) tgbotapi.MessageConfig {
				var msgTextbilder strings.Builder
				msgTextbilder.WriteString("На данный момент доступны:\n")
				for command, val := range commands {
					msgTextbilder.WriteString("\t")
					msgTextbilder.WriteString(command)
					msgTextbilder.WriteString(" - ")
					msgTextbilder.WriteString(val.description)
					msgTextbilder.WriteString("\n")
				}

				return tgbotapi.NewMessage(update.Message.Chat.ID, msgTextbilder.String())
			},
			description: "вывод текущего списока команд"},
		"/data": {
			handler: func(update tgbotapi.Update) tgbotapi.MessageConfig {
				var msgTextbilder strings.Builder
				msgTextbilder.WriteString("На данный момент доступны:\n")
				for shopname := range bot.allowedCards {
					msgTextbilder.WriteString("\t")
					msgTextbilder.WriteString(shopname)
					msgTextbilder.WriteString("\n")
				}
				return tgbotapi.NewMessage(update.Message.Chat.ID, msgTextbilder.String())
			},
			description: "вывод текущего списока магазинов"},
	}

	for update := range updates {
		if update.Message == nil {
			continue
		}
		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

		if reflect.TypeOf(update.Message.Text).Kind() == reflect.String && update.Message.Text != "" {
			if val, ok := commands[update.Message.Text]; ok {
				bot.tgApi.Send(val.handler(update))
			} else {
				photoname, ok := bot.allowedCards[update.Message.Text]

				if !ok {
					msg := tgbotapi.NewMessage(update.Message.Chat.ID,
						"На данный момент магазин "+update.Message.Text+" отсутствует, но мы скоро это исправим(;")
					bot.tgApi.Send(msg)
					continue
				}

				photo, err := os.ReadFile(bot.databasePath + "/" + photoname)
				if err != nil {
					msg := tgbotapi.NewMessage(update.Message.Chat.ID,
						"что-то пошло не так...(мы уже разбираемся)")
					bot.tgApi.Send(msg)
					continue
				}

				bot.tgApi.Send(tgbotapi.NewPhoto(update.Message.Chat.ID, tgbotapi.FileBytes{Name: "picture", Bytes: photo}))
			}
		}
	}
}

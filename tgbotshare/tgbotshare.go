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

	var buf struct {
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

		bot.allowedCards[buf.Shop] = buf.CardPath
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

	for update := range updates {
		if update.Message == nil {
			continue
		}
		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

		if reflect.TypeOf(update.Message.Text).Kind() == reflect.String && update.Message.Text != "" {
			switch update.Message.Text {
			case "/start":
				msg := tgbotapi.NewMessage(update.Message.Chat.ID,
					"Этот бот создан для cardshare, для получения карты, "+
						"которая есть в наличии просто напиши название магазина.\n"+
						"Для повторения этого сообщения напиши /start\n"+
						"Для ознокомления с полным списком команд /help")
				bot.tgApi.Send(msg)
			case "/help":
				msg := tgbotapi.NewMessage(update.Message.Chat.ID,
					"Вот список актуальных команд:\n"+
						"\t/start - сообщение начала работы\n"+
						"\t/help  - cписок команд\n"+
						"\t/data  - информация о текущих картах в доступе\n")
				bot.tgApi.Send(msg)
			case "/data":
				var msgTextbilder strings.Builder
				msgTextbilder.WriteString("На данный момент доступны:\n")
				for shopname := range bot.allowedCards {
					msgTextbilder.WriteString("\t")
					msgTextbilder.WriteString(shopname)
					msgTextbilder.WriteString("\n")
				}
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, msgTextbilder.String())
				bot.tgApi.Send(msg)
			default:
				photoname, ok := bot.allowedCards[update.Message.Text]

				if !ok {
					msg := tgbotapi.NewMessage(update.Message.Chat.ID,
						"На данный момент магазин "+update.Message.Text+" отсутствует, но мы скоро это исправим(;")
					bot.tgApi.Send(msg)
					break
				}

				photo, err := os.ReadFile(bot.databasePath + "/" + photoname)
				if err != nil {
					msg := tgbotapi.NewMessage(update.Message.Chat.ID,
						"что-то пошло не так...(мы уже разбираемся)")
					bot.tgApi.Send(msg)
					break
				}

				bot.tgApi.Send(tgbotapi.NewPhoto(update.Message.Chat.ID, tgbotapi.FileBytes{Name: "picture", Bytes: photo}))

			}
		}
	}
}

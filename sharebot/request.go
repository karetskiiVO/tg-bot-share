package sharebot

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type request interface {
	send()
}

type sendRequest struct {
	msg          *tgbotapi.MessageConfig
	upd          *tgbotapi.Update
	removeUpd    bool
	showKeyboard bool
	api          *tgbotapi.BotAPI
}

func (sndReq sendRequest) send() {
	if sndReq.upd != nil && sndReq.removeUpd {
		deleteReq := tgbotapi.NewDeleteMessage(sndReq.upd.Message.Chat.ID, sndReq.upd.Message.MessageID)
		sndReq.api.Send(deleteReq)
	}

	if sndReq.showKeyboard {
		sndReq.msg.ReplyMarkup = commandKeyboard
	}

	if _, err := sndReq.api.Send(sndReq.msg); err != nil {
		log.Panic(err)
	}
}

type keyboardUpdateRequest struct {
	replyMarkup tgbotapi.EditMessageReplyMarkupConfig
	api         *tgbotapi.BotAPI
}

func (kbdUpdReq keyboardUpdateRequest) send() {
	if _, err := kbdUpdReq.api.Send(kbdUpdReq.replyMarkup); err != nil {
		log.Panic(err)
	}
}

type photoSendUpdateRequest struct {
	photo tgbotapi.PhotoConfig
	api   *tgbotapi.BotAPI
}

func (phSndReq photoSendUpdateRequest) send() {
	if _, err := phSndReq.api.Send(phSndReq.photo); err != nil {
		log.Panic(err)
	}
}

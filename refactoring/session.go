package sharebot

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

type State interface {
	Handle(bot *Bot, update tgbotapi.Update) State
}

type Session struct {
	// user  tgbotapi.User
	state State
}

func (s *Session) Handle(bot *Bot, update tgbotapi.Update) {
	nextState := s.state.Handle(bot, update)
	s.state = nextState
}

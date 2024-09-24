package sharebot

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

type state interface {
	Handle(bot *Bot, update tgbotapi.Update, env *map[string]string) state
}

type session struct {
	// user  tgbotapi.User
	state state
	Env   *map[string]string
}

func newSession() *session {
	return &session{
		Env:   new(map[string]string),
		state: defaultState{},
	}
}

func (s *session) Handle(bot *Bot, update tgbotapi.Update) {
	nextState := s.state.Handle(bot, update, s.Env)
	s.state = nextState
}

type defaultState struct{}

func (defaultState) Handle(bot *Bot, update tgbotapi.Update, env *map[string]string) state {
	return defaultState{}
}

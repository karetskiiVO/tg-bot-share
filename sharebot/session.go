package sharebot

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

type Session struct {
	username   string
	state      UserState
	authorized bool
}

func (s Session) Username() string {
	return s.username
}

type UserTable struct {
	users map[string](*Session)
}

func NewUserTable() UserTable {
	return UserTable{make(map[string]*Session)}
}

func (ut *UserTable) GetSession(username string) *Session {
	res, ok := ut.users[username]

	if ok {
		return res
	} else {
		newSession := NewSession(username)
		ut.users[username] = &newSession
		return &newSession
	}
}

func NewSession(username string) Session {
	return Session{
		username:   username,
		state:      NotAuthorisedState{},
		authorized: false,
	}
}

func (s *Session) GetHandler(update tgbotapi.Update) messageHandler {
	handler, newstate := s.state.NextState(update)
	s.state = newstate
	return handler
}

package sharebot

import (
	"sync"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// "log"

type UserDataBase struct {
	mu    sync.RWMutex
	users map[tgbotapi.User]*Session

	// logger

	defaultSessionCtor func(id *tgbotapi.User)Session
}

func NewUserDataBase(defaultSessionCtor func(id *tgbotapi.User)Session) *UserDataBase {
	return &UserDataBase{
		users: make(map[tgbotapi.User]*Session),
		defaultSessionCtor: defaultSessionCtor,
	}
}

func (db *UserDataBase) СaptureSession(user *tgbotapi.User) *Session {
	db.mu.Lock()

	res, ok := db.users[*user]
	if ok {
		return res
	} else {
		newSession := new(Session)
		*newSession = db.defaultSessionCtor(user)
		db.users[*user] = newSession
		return newSession
	}
}

func (db *UserDataBase) ReleaseSession(user *tgbotapi.User) {
	db.mu.Unlock()
}

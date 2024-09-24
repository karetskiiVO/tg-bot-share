package sharebot

import (
	"sync"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// "log"

type userDataBase struct {
	mu    sync.RWMutex
	users map[tgbotapi.User]*session

	// logger

	defaultSessionCtor func(id *tgbotapi.User) session
}

func newUserDataBase(defaultSessionCtor func(id *tgbotapi.User) session) *userDataBase {
	return &userDataBase{
		users:              make(map[tgbotapi.User]*session),
		defaultSessionCtor: defaultSessionCtor,
	}
}

func (db *userDataBase) СaptureSession(user *tgbotapi.User) *session {
	db.mu.Lock()

	res, ok := db.users[*user]
	if ok {
		return res
	}
	newSession := new(session)
	*newSession = db.defaultSessionCtor(user)
	db.users[*user] = newSession
	return newSession
}

func (db *userDataBase) ReleaseSession(user *tgbotapi.User) {
	db.mu.Unlock()
}

func (db *userDataBase) LoadFromSQL(filepath string) {
	db.mu.Lock()

	panic("implement me")
	//db.mu.Unlock()
}

func (db *userDataBase) SaveToSQL(filepath string) {
	db.mu.RLock()

	db.mu.RUnlock()
}

package sharebot

import (
	"flag"
	"log"
)

type Context struct {
	TGToken  string
	DataBase ShopsDataBase
}

func GetContext() (Context, error) {
	var telegramBotToken, dbPath string

	flag.StringVar(&telegramBotToken, "telegrambottoken", "", "Telegram Bot Token")
	flag.StringVar(&dbPath, "datasource", "", "Database directory")
	flag.Parse()

	if telegramBotToken == "" {
		log.Fatal("-telegrambottoken is required")
	}
	if dbPath == "" {
		log.Fatal("-datasource is required")
	}

	return Context{
		TGToken:  telegramBotToken,
		DataBase: NewShopsDataBase(dbPath),
	}, nil
}

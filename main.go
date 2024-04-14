package main

import (
	"flag"
	"log"
	"os"
	"tgbotshare"
)

func main() {
	var telegramBotToken, dbPath string

	flag.StringVar(&telegramBotToken, "telegrambottoken", "", "Telegram Bot Token")
	flag.StringVar(&dbPath, "db", "", "Database directory")
	flag.Parse()

	if telegramBotToken == "" {
		log.Print("-telegrambottoken is required")
		os.Exit(1)
	}
	if dbPath == "" {
		log.Print("-db is required")
		os.Exit(1)
	}

	bot := tgbotshare.NewTGbot(telegramBotToken, dbPath)
	
	
	bot.Run()
}

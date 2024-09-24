package main

import "logger"
import "sharebot"

func main() {
	context, _ := sharebot.GetContext()

	bot, _ := sharebot.NewBot(context, new(logger.Logger))
	bot.SetInlineKeyboardSize(5, 2)
	bot.Run()
}

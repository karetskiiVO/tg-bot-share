package main

import "sharebot"

func main() {
	context, _ := sharebot.GetContext()

	bot, _ := sharebot.NewBot(context)
	bot.SetInlineKeyboardSize(5, 2)
	bot.Run()
}

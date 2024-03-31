package main

import "log"

func main() {
	bot, err := setupBot()
	log.Println("Bot is running")

	if err != nil {
		panic(err)
	}

	startWebhookServer(bot)
}

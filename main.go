package main

import (
	"log"
	"os"
)

func main() {
	botToken := getBotToken()

	db := NewDB(getDBDir())
	bot := NewBot(botToken, &db, isPrivate())
	go bot.Start()
	log.Println("Bot is running")

	startWebhookServer(bot)
}

func getDBDir() string {
	volumesDir := os.Getenv("DB_FOLDER")

	if volumesDir == "" {
		log.Fatalf("DB_FOLDER environment variable not set")
	}

	return volumesDir
}

func getBotToken() string {
	botToken := os.Getenv("TELEGRAM_BOT_API_TOKEN")

	if botToken == "" {
		log.Fatalln("TELEGRAM_BOT_API_TOKEN is not set")
	}

	return botToken
}

func isPrivate() bool {
	private := os.Getenv("IS_PRIVATE")

	if private == "" {
		return false
	}

	return true
}

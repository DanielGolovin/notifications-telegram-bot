package main

import (
	"encoding/json"
	"log"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func setupUpdateConfig() tgbotapi.UpdateConfig {
	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 30

	return updateConfig
}

func setupUpdateChannel(bot *tgbotapi.BotAPI) tgbotapi.UpdatesChannel {
	updateConfig := setupUpdateConfig()
	updatesChan := bot.GetUpdatesChan(updateConfig)
	chatsIdMap := loadChatsIdMap()

	for update := range updatesChan {
		handleUpdate(bot, update, chatsIdMap)
	}

	return updatesChan
}

func handleUpdate(bot *tgbotapi.BotAPI, update tgbotapi.Update, chatsIdMap map[int64]bool) {
	chatID := update.Message.Chat.ID
	addChat(chatsIdMap, chatID)
}

func setupBot() (*tgbotapi.BotAPI, error) {
	botToken := getBotToken()
	bot, err := tgbotapi.NewBotAPI(botToken)

	if err != nil {
		return nil, err
	}

	go setupUpdateChannel(bot)

	return bot, nil
}

func getChats() []int64 {
	chatsIdMap := loadChatsIdMap()

	var chats []int64

	for chatID := range chatsIdMap {
		chats = append(chats, chatID)
	}

	return chats
}

func addChat(chatsIdMap map[int64]bool, chatID int64) {
	chatsIdMap[chatID] = true

	saveChatsIdMap(chatsIdMap)
}

func saveChatsIdMap(chatsIdMap map[int64]bool) {
	file, err := os.Create("/data/chats-id-map.json")
	if err != nil {
		log.Fatalf("Error creating file: %v", err)
	}

	defer file.Close()

	jsonData, err := json.Marshal(chatsIdMap)

	if err != nil {
		log.Fatalf("Error stringifying JSON: %v", err)
	}

	_, err = file.Write(jsonData)

	if err != nil {
		log.Fatalf("Error writing to file: %v", err)
	}
}

func loadChatsIdMap() map[int64]bool {
	var chatsIdMap = make(map[int64]bool)

	file, err := os.Open("/data/chats-id-map.json")
	if err == os.ErrNotExist {
		log.Println("File not found, creating a new one")
		saveChatsIdMap(chatsIdMap)
		return chatsIdMap
	}

	if err != nil {
		log.Fatalf("Error opening file: %v", err)
	}

	defer file.Close()

	jsonData := make([]byte, 1024)
	n, err := file.Read(jsonData)

	if err != nil {
		log.Fatalf("Error reading file: %v", err)
	}

	err = json.Unmarshal(jsonData[:n], &chatsIdMap)

	if err != nil {
		log.Fatalf("Error unmarshalling JSON: %v", err)
	}

	return chatsIdMap
}

func getBotToken() string {
	botToken := os.Getenv("TELEGRAM_BOT_API_TOKEN")

	if botToken == "" {
		log.Fatalln("TELEGRAM_BOT_API_TOKEN is not set")
	}

	return botToken
}

func broadcastMessage(bot *tgbotapi.BotAPI, message string) {
	chats := getChats()

	for _, chatID := range chats {
		msg := tgbotapi.NewMessage(chatID, message)
		bot.Send(msg)
	}
}

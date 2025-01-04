package main

import (
	"fmt"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type NotificationBot struct {
	bot     *tgbotapi.BotAPI
	db      *DBInterface
	private bool
}

type DBInterface interface {
	AddChat(chatID int64)
	RemoveChat(chatID int64)
	GetChats() []int64
	GetWhiteList() map[int64]bool
}

func NewBot(token string, db *DBInterface, private bool) *NotificationBot {
	bot := &NotificationBot{}

	_bot, err := tgbotapi.NewBotAPI(token)

	if err != nil {
		log.Fatalf("Error creating bot: %v", err)
	}

	bot.bot = _bot
	bot.db = db
	bot.private = private

	return bot
}

func (bot *NotificationBot) Start() {
	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 30

	updatesChan := bot.bot.GetUpdatesChan(updateConfig)

	for update := range updatesChan {
		bot.handleUpdate(update)
	}
}

func (bot *NotificationBot) handleUpdate(update tgbotapi.Update) {
	(*bot.db).AddChat(update.Message.Chat.ID)
	bot.sendMessage(update.Message.Chat.ID, "Your user_id is: "+fmt.Sprint(update.Message.From.ID))
}

func (bot *NotificationBot) getChatsToNotify() []int64 {
	if bot.private {
		whiteList := (*bot.db).GetWhiteList()
		allowedChats := make([]int64, 0, len(whiteList))
		chats := (*bot.db).GetChats()

		for _, chatID := range chats {
			if whiteList[chatID] {
				allowedChats = append(allowedChats, chatID)
			} else {
				(*bot.db).RemoveChat(chatID)
				bot.sendMessage(chatID, "You are not allowed to receive notifications. Please contact the bot owner.")
			}
		}

		return allowedChats
	}

	return (*bot.db).GetChats()
}

func (bot *NotificationBot) sendMessage(chatID int64, message string) {
	msg := tgbotapi.NewMessage(chatID, message)
	bot.bot.Send(msg)
}

func (bot *NotificationBot) BroadcastMessage(message string) {
	chats := bot.getChatsToNotify()

	for _, chatID := range chats {
		bot.sendMessage(chatID, message)
	}
}

func (bot *NotificationBot) BroadcastFile(fileName string, fileData *[]byte) {
	chats := (*bot.db).GetChats()

	for _, chatID := range chats {
		msg := tgbotapi.NewDocument(chatID, tgbotapi.FileBytes{
			Name:  fileName,
			Bytes: *fileData,
		})
		bot.bot.Send(msg)
	}
}

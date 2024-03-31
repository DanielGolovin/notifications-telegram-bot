package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type NotificationRequestData struct {
	Data   map[string]interface{} `json:"data"`
	Secret string                 `json:"secret"`
}

func startWebhookServer(bot *tgbotapi.BotAPI) {
	mux := http.NewServeMux()

	notificationHandler := createHandlerWithBot(bot)

	mux.HandleFunc("POST /notify", notificationHandler)
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	log.Printf("Starting server on port %s", getServerPort())
	err := http.ListenAndServe(fmt.Sprintf(":%s", getServerPort()), mux)

	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}

func getServerPort() string {
	port := os.Getenv("SERVER_PORT")

	if port == "" {
		log.Fatalln("SERVER_PORT is not set")
	}

	return port
}

func parseNotificationRequest(r *http.Request) (*NotificationRequestData, error) {
	var data NotificationRequestData
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		return nil, err
	}

	return &data, nil
}

func validateSecret(secret string) bool {
	return secret == getSecret()
}

func createHandlerWithBot(bot *tgbotapi.BotAPI) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handleNotificationRequest(w, r, bot)
	}
}

func handleNotificationRequest(w http.ResponseWriter, r *http.Request, bot *tgbotapi.BotAPI) {
	data, err := parseNotificationRequest(r)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if !validateSecret(data.Secret) {
		http.Error(w, "Invalid secret", http.StatusUnauthorized)
		return
	}

	jsonData, err := json.Marshal(data.Data)
	if err != nil {
		log.Fatalf("Error stringifying JSON: %v", err)
	}

	stringData := string(jsonData)

	broadcastMessage(bot, stringData)

	w.WriteHeader(http.StatusOK)
}

func getSecret() string {
	botToken := os.Getenv("NOTIFICATION_SECRET")

	if botToken == "" {
		log.Fatalln("NOTIFICATION_SECRET is not set")
	}

	return botToken
}

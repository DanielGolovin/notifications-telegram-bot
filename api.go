package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

type NotificationRequestData struct {
	Data map[string]interface{} `json:"data"`
}

func startWebhookServer(bot *NotificationBot) {
	mux := http.NewServeMux()

	notificationHandler := createNotificationHandler(bot)
	fileNotificationHandler := createFileNotificationHandler(bot)

	mux.HandleFunc("POST /notify", notificationHandler)
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	mux.HandleFunc("POST /notify/file", fileNotificationHandler)

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
	defer r.Body.Close()
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

func createNotificationHandler(bot *NotificationBot) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handleNotificationRequest(w, r, bot)
	}
}

func createFileNotificationHandler(bot *NotificationBot) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handleFileNotificationRequest(w, r, bot)
	}
}

func handleNotificationRequest(w http.ResponseWriter, r *http.Request, bot *NotificationBot) {
	secret := r.Header.Get("X-Secret")

	if !validateSecret(secret) {
		http.Error(w, "Invalid secret", http.StatusUnauthorized)
		return
	}

	data, err := parseNotificationRequest(r)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	jsonData, err := json.Marshal(data.Data)
	if err != nil {
		log.Fatalf("Error stringifying JSON: %v", err)
	}

	stringData := string(jsonData)

	log.Printf("Broadcasting message: %s", stringData)

	bot.BroadcastMessage(stringData)

	w.WriteHeader(http.StatusOK)
}

func parseFileNotificationRequest(r *http.Request) (*[]byte, error) {
	log.Printf("Received file notification request")
	defer r.Body.Close()
	fileData, err := io.ReadAll(r.Body)

	if err != nil {
		log.Fatalf("Error reading file: %v", err)
		return nil, err
	}

	return &fileData, nil
}

func handleFileNotificationRequest(w http.ResponseWriter, r *http.Request, bot *NotificationBot) {
	data, err := parseFileNotificationRequest(r)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	log.Printf("Broadcasting file: %v", data)

	fileName := r.Header.Get("X-File-Name")
	if fileName == "" {
		fileName = "file"
	}

	bot.BroadcastFile(fileName, data)

	w.WriteHeader(http.StatusOK)
}

func getSecret() string {
	botToken := os.Getenv("NOTIFICATION_SECRET")

	if botToken == "" {
		log.Fatalln("NOTIFICATION_SECRET is not set")
	}

	return botToken
}

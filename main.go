package main

import (
	"log"
	"os"
)

var volumesDir = "/volumes"

func main() {
	ensureVolumesDirExists()

	bot, err := setupBot()
	log.Println("Bot is running")

	if err != nil {
		panic(err)
	}

	startWebhookServer(bot)
}

func ensureVolumesDirExists() {
	err := os.MkdirAll(volumesDir, os.ModePerm)
	if err != nil {
		log.Fatalf("Error creating volumes directory: %v", err)
	}

	log.Printf("Volumes directory created at %s", volumesDir)
}

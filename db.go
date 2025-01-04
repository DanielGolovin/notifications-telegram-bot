package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
)

type DB struct {
	data struct {
		AllowedUsers map[int64]bool `json:"allowedUsers"`
		Chats        map[int64]bool `json:"chats"`
	}
	dbPath string
}

func NewDB(dbPath string) DBInterface {
	db := &DB{}
	db.initDB(dbPath)
	return db
}

func (db *DB) AddChat(chatID int64) {
	db.loadData()
	db.data.Chats[chatID] = true
	db.saveChatsIdMap()
}

func (db *DB) RemoveChat(chatID int64) {
	db.loadData()
	delete(db.data.Chats, chatID)
	db.saveChatsIdMap()
}

func (db *DB) GetWhiteList() map[int64]bool {
	db.loadData()

	return db.data.AllowedUsers
}

func (db *DB) GetChats() []int64 {
	db.loadData()

	chats := make([]int64, 0, len(db.data.Chats))

	for chatID := range db.data.Chats {
		chats = append(chats, chatID)
	}

	return chats
}

func (db *DB) initDB(dbPath string) {
	err := os.MkdirAll(dbPath, os.ModePerm)
	if err != nil {
		log.Fatalf("Error creating volumes directory: %v", err)
	}

	log.Printf("Using DB directory at %s", dbPath)

	db.dbPath = fmt.Sprintf("%s/db.json", dbPath)

	log.Printf("DB file path: %s", db.dbPath)

	db.data = struct {
		AllowedUsers map[int64]bool `json:"allowedUsers"`
		Chats        map[int64]bool `json:"chats"`
	}{}

	_, err = os.Stat(db.dbPath)

	if os.IsNotExist(err) {
		os.Create(db.dbPath)
		os.WriteFile(db.dbPath, []byte("{\"allowedUsers\": {},\"chats\": {}}"), os.ModePerm)
		log.Printf("DB file created at %s", db.dbPath)
	} else {
		log.Printf("DB file found at %s", db.dbPath)
	}
	db.loadData()
}

func (db *DB) loadData() {
	file, err := os.Open(db.dbPath)
	if err != nil {
		log.Fatalf("Error opening file: %v", err)
	}
	defer file.Close()

	jsonData, err := io.ReadAll(file)
	if err != nil {
		log.Fatalf("Error reading file: %v", err)
	}

	err = json.Unmarshal(jsonData, &db.data)
	if err != nil {
		log.Fatalf("Error parsing JSON: %v", err)
	}
}

func (db *DB) saveChatsIdMap() {
	jsonData, err := json.Marshal(db.data)

	if err != nil {
		log.Fatalf("Error stringifying JSON: %v", err)
	}

	err = os.WriteFile(db.dbPath, jsonData, os.ModePerm)

	if err != nil {
		log.Fatalf("Error writing to file: %v", err)
	}
}

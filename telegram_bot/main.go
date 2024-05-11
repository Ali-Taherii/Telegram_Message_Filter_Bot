package main

import (
	"fmt"
	"log"
	"os"
	"telegram_bot/structs"
)

func main() {
	// Read PostgreSQL password from environment variable
	password := os.Getenv("POSTGRES_PASSWORD")
	if password == "" {
		log.Fatal("POSTGRES_PASSWORD environment variable is not set")
	}

	// Open database connection with password
	db, err := structs.NewDB(fmt.Sprintf("postgresql://postgres:%s@31.216.88.247/Telegram_Filter_Bot?sslmode=disable", password))
	if err != nil {
		log.Fatal("Error connecting to database:", err)
	}
	defer db.Close()

	// Get bot token from environment variable
	botToken := os.Getenv("BOT_TOKEN")
	if botToken == "" {
		log.Fatal("BOT_TOKEN environment variable is not set")
	}

	// Initialize the bot
	bot, err := structs.NewBot(botToken, db)
	if err != nil {
		log.Panic(err)
	}

	// Bot runs
	bot.StartListening()
}

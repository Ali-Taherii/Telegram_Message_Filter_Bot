package main

import (
	"fmt"
	"log"
	"os"
	"telegram_bot/structs"

	"golang.org/x/crypto/ssh/terminal"
)

func main() {
	// Prompt for password
	fmt.Print("Enter the PostgreSQL password: ")
	password, err := terminal.ReadPassword(int(os.Stdin.Fd()))
	if err != nil {
		log.Fatal("Error reading password:", err)
	}

	// Open database connection with password
	db, err := structs.NewDB(fmt.Sprintf("postgresql://postgres:%s@localhost:5432/Telegram_Filter_Bot?sslmode=disable", string(password)))
	if err != nil {
		log.Fatal("Error connecting to database:", err)
	}
	defer db.Close()

	// Initialize the bot
	bot, err := structs.NewBot("7031750953:AAGgNTyWcTJbGUVkj6VoYDLfjQ2F3FFmXIo", db)
	if err != nil {
		log.Panic(err)
	}

	// Bot runs
	bot.StartListening()
}

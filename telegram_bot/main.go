package main

import (
	"log"
	"telegram_bot/structs"
)

func main() {

	// Initialize the bot
	bot, err := structs.NewBot("7031750953:AAGgNTyWcTJbGUVkj6VoYDLfjQ2F3FFmXIo")
	if err != nil {
		log.Panic(err)
	}

	// Bot waits for a message
	bot.StartListening()
}

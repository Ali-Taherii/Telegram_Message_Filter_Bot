package structs

import (
	"log"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type TeleBot struct {
	API            *tgbotapi.BotAPI
	StoredWord     string
	WaitingForWord bool // New state variable to track if bot is waiting for filter word
}

func NewBot(token string) (*TeleBot, error) {
	botAPI, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, err
	}
	return &TeleBot{API: botAPI}, nil
}

func (b *TeleBot) StartListening() {
	b.API.Debug = true
	log.Printf("Authorized on account %s", b.API.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := b.API.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil {
			if update.Message.IsCommand() && update.Message.Command() == "start" {
				b.Start(update)
			} else if update.Message.IsCommand() && update.Message.Command() == "filter" {
				b.Filter(update)
				b.WaitingForWord = true // Set WaitingForWord to true when /filter command received
			} else {
				if b.WaitingForWord {
					b.WordReceiver(update)
				}
			} /*else {
				b.ProcessMessage(update)
			}*/
		}
	}
}

func (b *TeleBot) Start(update tgbotapi.Update) {
	reply := "Welcome! This bot will first ask you for a word, and then for a sentence. It will then check if the sentence contains the word or not.\nUse /filter to define the filer word"
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, reply)
	b.API.Send(msg)
}

func (b *TeleBot) Filter(update tgbotapi.Update) {

	// Prompt user to enter a word
	reply := "Write the filter word (one word only)"
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, reply)
	b.API.Send(msg)
}

func (b *TeleBot) WordReceiver(update tgbotapi.Update) (wordReceived bool) {
	// Reset the forward watcher
	b.WaitingForWord = false

	// Split the sentence into words
	words := strings.Fields(update.Message.Text)

	// The user input is more than one word
	if len(words) > 1 {
		reply := "Please provide only one word. Try /filter again."
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, reply)
		b.API.Send(msg)
		return false
	}

	// Store the word to process
	b.StoredWord = words[0]
	reply := "Word received. Please send a sentence containing the word in the next message."
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, reply)
	b.API.Send(msg)

	return true
}

func (b *TeleBot) ProcessMessage(update tgbotapi.Update) {
	if b.StoredWord == "" {
		// Split the sentence into words
		words := strings.Fields(update.Message.Text)

		// The user input is more than one word
		if len(words) > 1 {
			reply := "Please provide only one word. Try again."
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, reply)
			b.API.Send(msg)
			return
		}

		// Store the word to process
		b.StoredWord = words[0]
		reply := "Word received. Please send a sentence containing the word in the next message."
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, reply)
		b.API.Send(msg)
	} else {
		// Split the sentence into words
		words := strings.Fields(update.Message.Text)

		// Check if the stored word is present in the sentence as a whole word
		found := false
		for _, word := range words {
			if strings.EqualFold(word, b.StoredWord) {
				found = true
				break
			}
		}

		// Respond based on whether the word is found or not
		if found {
			reply := "The sentence contains the word!"
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, reply)
			b.API.Send(msg)
		} else {
			reply := "The sentence doesn't contain the word. Please try again."
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, reply)
			b.API.Send(msg)
		}

		// Clear the stored word
		b.StoredWord = ""
	}
}

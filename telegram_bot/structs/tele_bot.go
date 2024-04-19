package structs

import (
	"log"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type TeleBot struct {
	API            *tgbotapi.BotAPI
	StoredWord     string
	WaitingForWord bool // New state variable to track if bot is waiting for filter word
	DB             *DB  // Database connection
}

// Initialize the bot
func NewBot(token string, db *DB) (*TeleBot, error) {
	botAPI, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, err
	}
	return &TeleBot{API: botAPI, DB: db}, nil
}

// Bot runs and gets messages from the user
func (b *TeleBot) StartListening() {
	b.API.Debug = true
	log.Printf("Authorized on account %s", b.API.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := b.API.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil {
			if update.Message.IsCommand() {
				switch update.Message.Command() {
				case "start":
					b.Start(update)
				case "filter":
					b.WaitingForWord = true // Set WaitingForWord to true when /filter command received
					b.Filter(update)
				default:
					b.ProcessMessage(update)
				}
			} else {
				if b.WaitingForWord {
					b.WordReceiver(update)
				} else {
					b.ProcessMessage(update)
				}
			}
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

func (b *TeleBot) WordReceiver(update tgbotapi.Update) {
	// Not waiting for filter word
	b.WaitingForWord = false

	// Split the sentence into words
	words := strings.Fields(update.Message.Text)

	// The user input is more than one word
	if len(words) > 1 {
		reply := "Please provide only one word. Try /filter again."
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, reply)
		b.API.Send(msg)
		return
	}

	// Store the word to process
	b.StoredWord = words[0]
	reply := "Word received. Please send a sentence in the next messages."
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, reply)
	msg.ReplyToMessageID = update.Message.MessageID
	b.API.Send(msg)
}

func (b *TeleBot) ProcessMessage(update tgbotapi.Update) {

	// No filter word yet entered
	if b.StoredWord == "" {
		reply := "No filter word found. Use /filter to enter one"
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, reply)
		b.API.Send(msg)
		return
	}

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

	// Store the message in the appropriate table based on whether the word is found
	if found {
		// Message contains the filter word, store it in messages_with_filter table
		err := b.DB.StoreMessage(update.Message.From.ID, update.Message.Text, time.Now(), b.StoredWord, "messages_with_word")
		if err != nil {
			log.Println("Error storing message with filter word:", err)
		}
	} else {
		// Message does not contain the filter word, store it in messages_without_filter table
		err := b.DB.StoreMessage(update.Message.From.ID, update.Message.Text, time.Now(), b.StoredWord, "messages_without_word")
		if err != nil {
			log.Println("Error storing message without filter word:", err)
		}
	}

	// Respond based on whether the word is found or not
	if found {
		reply := "The sentence contains the word!"
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, reply)
		msg.ReplyToMessageID = update.Message.MessageID
		b.API.Send(msg)
	} else {
		reply := "The sentence doesn't contain the word. Please try again."
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, reply)
		msg.ReplyToMessageID = update.Message.MessageID
		b.API.Send(msg)
	}
}

package structs

import (
	"fmt"
	"log"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type TeleBot struct {
	API            *tgbotapi.BotAPI
	FilterWord     string
	SearchWord     string
	WaitingForWord bool // New state variable to track if bot is waiting for filter word
	IsSearching    bool // New state variable to track if the bot is searching messages
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
				case "stop":
					b.Stop(update)
					return // Stop processing updates
				case "help":
					b.Help(update)
				case "show":
					b.Show(update)
				default:
					b.ProcessMessage(update)
				}
			} else {
				if b.WaitingForWord {
					b.WordReceiver(update)
				} else if !b.WaitingForWord && !b.IsSearching {
					b.ProcessMessage(update)
				}
			}
		} else if update.CallbackQuery != nil {
			// Handle callback query
			b.HandleCallbackQuery(update)
		}
	}
}

func (b *TeleBot) Start(update tgbotapi.Update) {
	reply := "Welcome! This bot will first ask you for a word, and then for a sentence. It will then check if the sentence contains the word or not.\nUse /filter to define the filter word\nUse /show to search for messages"
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

	if !b.IsSearching {
		b.FilterWord = words[0] // Store filter word
		reply := "Word received.\nPlease send a sentence in the next messages."
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, reply)
		msg.ReplyToMessageID = update.Message.MessageID
		b.API.Send(msg)
	} else {
		b.SearchWord = words[0] // Store search word
		reply := "Word received.\nSearching for the messages with this filter word."
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, reply)
		msg.ReplyToMessageID = update.Message.MessageID
		b.API.Send(msg)
		b.SearchMessage(update)
	}

}

// Stop closes the database connection
func (b *TeleBot) Stop(update tgbotapi.Update) {
	reply := "Stopping the bot.\nClosing database connection."
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, reply)
	b.API.Send(msg)
	b.CloseDB()
}

// CloseDB closes the database connection
func (b *TeleBot) CloseDB() {
	if b.DB != nil {
		b.DB.Close()
	}
}

func (b *TeleBot) ProcessMessage(update tgbotapi.Update) {

	// No filter word yet entered
	if b.FilterWord == "" {
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
		if strings.EqualFold(word, b.FilterWord) {
			found = true
			break
		}
	}

	// Store the message in the appropriate table based on whether the word is found
	if found {
		// Message contains the filter word, store it in messages_with_filter table
		err := b.DB.StoreMessage(update.Message.From.ID, update.Message.Text, time.Now(), b.FilterWord, "messages_with_word")
		if err != nil {
			log.Println("Error storing message with filter word:", err)
		}
	} else {
		// Message does not contain the filter word, store it in messages_without_filter table
		err := b.DB.StoreMessage(update.Message.From.ID, update.Message.Text, time.Now(), b.FilterWord, "messages_without_word")
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

// Help command
func (b *TeleBot) Help(update tgbotapi.Update) {
	reply := "Available commands:\n" +
		"/start - Start the bot\n" +
		"/filter - Define a filter word\n" +
		"/stop - Stop the bot\n" +
		"/show - Show the stored messages" +
		"/help - Display this help message"
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, reply)
	b.API.Send(msg)
}

// Show command
func (b *TeleBot) Show(update tgbotapi.Update) {
	// Create an inline keyboard markup
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		// Create a row of buttons
		tgbotapi.NewInlineKeyboardRow(
			// Create the first button for showing messages with a filter word
			tgbotapi.NewInlineKeyboardButtonData("Show messages with filter word", "show_with_filter"),
		),
		// Create another row of buttons
		tgbotapi.NewInlineKeyboardRow(
			// Create the second button for showing messages without a filter word
			tgbotapi.NewInlineKeyboardButtonData("Show messages without filter word", "show_without_filter"),
		),
	)

	// Create a message with the inline keyboard markup
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Choose an option:")
	msg.ReplyMarkup = keyboard

	// Send the message with the menu buttons
	_, err := b.API.Send(msg)
	if err != nil {
		log.Println("Error sending message:", err)
	}
}

// Retrieve messages with a filter word
func (b *TeleBot) SearchMessage(update tgbotapi.Update) {
	b.IsSearching = false

	// Construct the query string by concatenating the SearchWord variable
	query := fmt.Sprintf("SELECT sender_id, message_text, sent_date FROM messages_with_word WHERE filter_word = '%s'", b.SearchWord)
	rows, err := b.DB.QueryRows(query)
	if err != nil {
		log.Println("Error executing query:", err)
		return
	}
	defer rows.Close()

	// Iterate over the rows and build the message
	var message string
	for rows.Next() {
		var senderID int64
		var messageText string
		var sentDate time.Time
		if err := rows.Scan(&senderID, &messageText, &sentDate); err != nil {
			log.Println("Error scanning row:", err)
			continue
		}
		// Format each row into a readable message
		message += fmt.Sprintf("Sender ID: %d\nMessage: %s\nSent Date: %s\n\n", senderID, messageText, sentDate.String())
	}

	// Check if no messages were found
	if message == "" {
		message = "No messages found."
	}

	// Send the message with the retrieved data or the "No messages found" message
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, message)
	_, err = b.API.Send(msg)
	if err != nil {
		log.Println("Error sending message:", err)
	}
}

// HandleCallbackQuery handles callback queries received when a user clicks on the inline keyboard buttons
func (b *TeleBot) HandleCallbackQuery(update tgbotapi.Update) {
	// Extract the callback data from the update
	callbackData := update.CallbackQuery.Data

	// Handle the callback data accordingly
	switch callbackData {
	case "show_with_filter":
		// Prompt the user to enter the filter word
		reply := "Please enter the filter word:"
		msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, reply)
		b.API.Send(msg)

		b.WaitingForWord = true // Set the bot's state to wait for the filter word
		b.IsSearching = true    // Additional state variable to indicate filtering messages with a word

	case "show_without_filter":
		// Execute the SQL query to retrieve messages without a filter word
		query := "SELECT sender_id, message_text, sent_date FROM messages_without_word"
		rows, err := b.DB.QueryRows(query)
		if err != nil {
			log.Println("Error executing query:", err)
			return
		}
		defer rows.Close()

		// Iterate over the rows and build the message
		var message string
		for rows.Next() {
			var senderID int64
			var messageText string
			var sentDate time.Time
			if err := rows.Scan(&senderID, &messageText, &sentDate); err != nil {
				log.Println("Error scanning row:", err)
				continue
			}
			// Format each row into a readable message
			message += fmt.Sprintf("Sender ID: %d\nMessage: %s\nSent Date: %s\n\n", senderID, messageText, sentDate.String())
		}

		// Check if no messages were found
		if message == "" {
			message = "No messages found."
		}

		// Send the message with the retrieved data
		msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, message)
		_, err = b.API.Send(msg)
		if err != nil {
			log.Println("Error sending message:", err)
		}
	}
}

package structs

import (
	"database/sql"
	"log"
	"time"

	_ "github.com/lib/pq"
)

type DB struct {
	*sql.DB
}

func NewDB(dataSourceName string) (*DB, error) {
	db, err := sql.Open("postgres", dataSourceName)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return &DB{db}, nil
}

// StoreMessage stores a message in the appropriate table based on whether it contains the filter word
func (db *DB) StoreMessage(senderID int64, messageText string, sentDate time.Time, filterWord string, tableName string) error {

	query := `
        INSERT INTO ` + tableName + ` (sender_id, message_text, sent_date, filter_word)
        VALUES ($1, $2, $3, $4)
    `
	_, err := db.Exec(query, senderID, messageText, sentDate, filterWord)
	if err != nil {
		log.Printf("Error storing message in %s table: %v\n", tableName, err)
	}
	return err
}

// QueryRows executes a SQL query and returns the result rows
func (db *DB) QueryRows(query string, args ...interface{}) (*sql.Rows, error) {
	rows, err := db.Query(query, args...)
	if err != nil {
		log.Printf("Error executing query: %v\n", err)
		return nil, err
	}
	return rows, nil
}

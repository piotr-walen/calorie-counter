package models

import (
	"database/sql"
	"time"

	"github.com/pkg/errors"
)

// Account struct is used to represent user acc
type Entry struct {
	ID        int
	UserID    int
	ProductID int
	Quantity  float64
}

func MigrateEntries(db *sql.DB) error {
	rows, err := db.Query(`
		CREATE TABLE IF NOT EXISTS entries (
			id SERIAL PRIMARY KEY, 
			user_id INTEGER REFERENCES accounts(id),
			product_id INTEGER REFERENCES products(id),  
			quantity DECIMAL,
			date DATE NOT NULL DEFAULT CURRENT_DATE
		);
	`)
	defer rows.Close()
	return err
}

func CreateEntry(db *sql.DB, entry Entry) error {
	rows, err := db.Query(`
		INSERT INTO entries (user_id, product_id, quantity)
		VALUES ($1, $2, $3);
	`, entry.UserID, entry.ProductID, entry.Quantity)
	defer rows.Close()
	return err
}

func GetEntry(db *sql.DB, id int) (*Entry, error) {
	rows, err := db.Query(`
		SELECT * FROM entries WHERE id = $1
	`, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	entries := []*Entry{}
	for rows.Next() {
		entry := &Entry{}
		rows.Scan(entry.ID, entry.UserID, entry.ProductID, entry.Quantity)
		entries = append(entries, entry)
	}
	if len(entries) > 1 {
		return nil, errors.New("Two entries with the same id")
	}
	return entries[0], err
}

func GetUsersEntries(db *sql.DB, userID int, date time.Time) ([]*Entry, error) {
	rows, err := db.Query(`
		SELECT * FROM entries WHERE user_id=$1
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	entries := []*Entry{}
	for rows.Next() {
		entry := &Entry{}
		rows.Scan(entry.ID, entry.UserID, entry.ProductID, entry.Quantity)
		entries = append(entries, entry)
	}
	return entries, err

}

func DeleteEntry(db *sql.DB, id int) error {
	rows, err := db.Query(`
		DELETE FROM entries WHERE id=$1 
	`, id)
	defer rows.Close()
	if err != nil {
		return errors.Wrap(err, "While deleting entry")
	}
	return nil
}

func UpdateEntry(db *sql.DB, id int, new Entry) error {
	rows, err := db.Query(`
		UPDATE products SET user_id=$2, product_id=$3, quantity=$4 WHERE id=$1;
	`, id, new.UserID, new.ProductID, new.Quantity)
	defer rows.Close()
	if err != nil {
		return errors.Wrap(err, "While updating entry")
	}
	return nil
}

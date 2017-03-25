package migrations

import "time"

// Store defines the methods required for storing
// the state of a database's migrations
type Store interface {
	History() ([]string, error)
	MostRecent() (int, string, error)
	SaveMigration(r Row) error
}

// Row is the row in a miration saved
type Row struct {
	Order   int
	Key     string
	Created time.Time
}

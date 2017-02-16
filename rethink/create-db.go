package rethinkmigrations

import (
	"fmt"

	"github.com/dancannon/gorethink"
	"github.com/dewey4iv/migrations"
)

// CreateDB returns a migration that creates a database
func CreateDB(session *gorethink.Session, database string) migrations.Migration {
	return &createDB{session, database}
}

type createDB struct {
	session  *gorethink.Session
	database string
}

func (m *createDB) Migrate() error {
	result, err := gorethink.DBCreate(m.database).RunWrite(m.session)
	if err != nil {
		return err
	}

	if result.DBsCreated < 1 {
		return fmt.Errorf("database %s not created", m.database)
	}

	return nil
}
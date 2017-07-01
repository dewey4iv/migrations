package rethinkmigrations

import (
	"fmt"

	"github.com/GoRethink/gorethink"
	"github.com/dewey4iv/migrations"
)

// DropDB returns a migration that creates a database
func DropDB(session *gorethink.Session, database string) migrations.Migration {
	return &dropDB{session, database}
}

type dropDB struct {
	session  *gorethink.Session
	database string
}

func (m *dropDB) Migrate() error {
	result, err := gorethink.DBDrop(m.database).RunWrite(m.session)
	if err != nil {
		return err
	}

	if result.DBsCreated < 1 {
		return fmt.Errorf("database %s not dropped", m.database)
	}

	return nil
}

func (m *dropDB) Key(order int) string {
	return fmt.Sprintf(`order: %d | drop database "%s"`, order, m.database)
}

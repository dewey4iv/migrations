package rethinkmigrations

import (
	"fmt"

	"github.com/GoRethink/gorethink"
	"github.com/dewey4iv/migrations"
)

// CreateTable returns a migration that creates a table if it doesn't exist
func CreateTable(session *gorethink.Session, table string) migrations.Migration {
	return &createTable{session, table}
}

type createTable struct {
	session *gorethink.Session
	table   string
}

func (m *createTable) Migrate() error {
	result, err := gorethink.TableCreate(m.table).RunWrite(m.session)
	if err != nil {
		return err
	}

	if result.TablesCreated < 1 {
		return fmt.Errorf("table %s not created", m.table)
	}

	return nil
}

func (m *createTable) Key(order int) string {
	return fmt.Sprintf(`order: %d | create table "%s"`, order, m.table)
}

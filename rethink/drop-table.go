package rethinkmigrations

import (
	"fmt"

	"github.com/dancannon/gorethink"
	"github.com/dewey4iv/migrations"
)

// DropTable returns a migration that drops a table
func DropTable(session *gorethink.Session, table string) migrations.Migration {
	return &dropTable{session, table}
}

type dropTable struct {
	session *gorethink.Session
	table   string
}

func (m *dropTable) Migrate() error {
	result, err := gorethink.TableDrop(m.table).RunWrite(m.session)
	if err != nil {
		return err
	}

	if result.TablesCreated < 1 {
		return fmt.Errorf("table %s not dropped", m.table)
	}

	return nil
}

func (m *dropTable) Key(order int) string {
	return fmt.Sprintf(`order: %d | drop table "%s"`, order, m.table)
}

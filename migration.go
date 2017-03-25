package migrations

import (
	"log"
	"time"
)

// Migration defines a simple interface for anything that can be migrated
type Migration interface {
	Migrate() error
	Key(int) string
}

func New(store Store) *Migrations {
	return &Migrations{store: store}
}

type Migrations struct {
	store      Store
	migrations []Migration
}

func (m *Migrations) Run() error {
	var row Row

	mostRecentOrder, _, err := m.store.MostRecent()
	if err != nil {
		return err
	}

	log.Printf("most recent: %d", mostRecentOrder)

	i := 1
	for _, migration := range m.migrations {
		i++

		if i <= mostRecentOrder {
			log.Printf("Skipping %s -- already migrated", migration.Key(i))
			continue
		}

		if err := migration.Migrate(); err != nil {
			return err
		}

		row = Row{
			Order:   i,
			Key:     migration.Key(i),
			Created: time.Now(),
		}

		if err := m.store.SaveMigration(row); err != nil {
			return err
		}
	}

	return nil
}

func (m *Migrations) Add(migration Migration) {
	m.migrations = append(m.migrations, migration)
}

package migrations

// Migration defines a simple interface for anything that can be migrated
type Migration interface {
	Migrate() error
}

func New() *Migrations {
	return &Migrations{}
}

type Migrations struct {
	migrations []Migration
}

func (m *Migrations) Run() error {
	for _, migration := range m.migrations {
		if err := migration.Migrate(); err != nil {
			return err
		}
	}

	return nil
}

func (m *Migrations) Add(migration Migration) {
	m.migrations = append(m.migrations, migration)
}

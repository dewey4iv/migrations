package rethinkmigrations

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/dancannon/gorethink"
	"github.com/dewey4iv/migrations"
)

func New(database string, session *gorethink.Session) (*Store, error) {
	s := Store{
		db:      gorethink.DB(database).Table("migrations"),
		session: session,
	}

	mostRecent, _, err := s.MostRecent()
	if err != nil && !strings.Contains(err.Error(), fmt.Sprintf("Database `%s` does not exist", database)) {
		return nil, err
	}

	log.Printf("most rec: %d", mostRecent)

	if mostRecent == 0 {
		dbMig := CreateDB(session, database)

		if err := dbMig.Migrate(); err != nil {
			return nil, err
		}

		tableMig := CreateTable(session, "migrations")

		if err := tableMig.Migrate(); err != nil {
			return nil, err
		}

		if err := s.SaveMigration(migrations.Row{Order: 0, Key: dbMig.Key(0), Created: time.Now()}); err != nil {
			return nil, err
		}

		if err := s.SaveMigration(migrations.Row{Order: 1, Key: tableMig.Key(1), Created: time.Now()}); err != nil {
			return nil, err
		}
	}

	return &s, nil
}

// Store implements a migrations.Store
type Store struct {
	db      gorethink.Term
	session *gorethink.Session
}

func (s *Store) History() ([]string, error) {
	var migrations []string

	results, err := s.db.OrderBy("Order").Run(s.session)
	if err != nil {
		return nil, err
	}
	defer results.Close()

	if err := results.All(&migrations); err != nil {
		return nil, err
	}

	return migrations, nil
}

func (s *Store) MostRecent() (int, string, error) {
	var r migrations.Row

	result, err := s.db.OrderBy(gorethink.Desc("Order")).Run(s.session)
	if err != nil {
		return 0, "", err
	}
	defer result.Close()

	if err := result.One(&r); err != nil {
		return 0, "", err
	}

	return r.Order, r.Key, nil
}

func (s *Store) SaveMigration(r migrations.Row) error {
	if _, err := s.db.Insert(r).RunWrite(s.session); err != nil {
		return err
	}

	return nil
}

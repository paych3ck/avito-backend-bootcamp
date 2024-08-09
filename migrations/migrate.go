package migrations

import (
	"database/sql"
	"log"
	"os"
	"path/filepath"
)

func Migrate(db *sql.DB) {
	migrationFIles := []string{
		"create_table_users.sql",
		"create_table_houses.sql",
		"create_table_flats.sql",
	}

	for _, file := range migrationFIles {
		path := filepath.Join("migrations", file)
		sqlCode, err := os.ReadFile(path)

		if err != nil {
			log.Fatalf("Error while reading migration file %s: %v", path, err)
		}

		_, err = db.Exec(string(sqlCode))

		if err != nil {
			log.Fatalf("Error while executing migration %s: %v", path, err)
		}
	}
}

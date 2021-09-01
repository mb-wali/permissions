package test

import (
	"database/sql"
	"fmt"
	"os"
	"testing"

	_ "github.com/lib/pq"
)

const FakeID = "1A960034-969A-46A7-B6B5-3F1866258CAB"

func shouldRun() bool {
	return os.Getenv("RUN_INTEGRATION_TESTS") != ""
}

func dburi() string {
	uri := os.Getenv("DBURI")
	if uri != "" {
		return uri
	}
	return "postgres://de:notprod@dedb:5432/permissions?sslmode=disable"
}

func schema() string {
	dbschema := os.Getenv("DBSCHEMA")
	if dbschema != "" {
		return dbschema
	}
	return "permissions"
}

func truncateTables(db *sql.DB, schema string) error {

	// Truncate all tables.
	tables := []string{"permissions", "subjects", "resources", "resource_types"}
	for _, table := range tables {
		_, err := db.Exec(fmt.Sprintf("DELETE FROM %s.%s", schema, table))
		if err != nil {
			return err
		}
	}

	return nil
}

func initdb(t *testing.T) (*sql.DB, string) {
	db, err := sql.Open("postgres", dburi())
	if err != nil {
		t.Error(err)
	}
	err = db.Ping()
	if err != nil {
		t.Error(err)
	}

	dbschema := schema()

	if err := truncateTables(db, dbschema); err != nil {
		t.Error(err)
	}

	return db, dbschema
}

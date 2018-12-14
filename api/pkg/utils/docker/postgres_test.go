// +build integration

package docker

import (
	"database/sql"
	"errors"
	"testing"

	_ "github.com/lib/pq"
)

func TestRunPostgres(t *testing.T) {
	Configure()

	dsn, close := RunPostgres("")
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		t.Fatal(err)
	}

	// run a query
	var id int
	rows, err := db.Query("select 1 as id;")
	if err != nil {
		t.Error(err)
	}

	// parse returned data
	rows.Next()
	if err := rows.Scan(&id); err != nil {
		t.Error(err)
	}
	if id != 1 {
		t.Error("id should be 1")
	}

	// destroy container
	close()

	// run a query again
	_, err = db.Query("select 1 as id;")
	if err == nil {
		t.Error(errors.New("should be error"))
	}
}

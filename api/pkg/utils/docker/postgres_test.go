// +build integration

package docker

import (
	"database/sql"
	"testing"

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

func TestRunPostgres(t *testing.T) {
	Configure()

	dsn, close := RunPostgres("")
	db, err := sql.Open("postgres", dsn)
	assert.Nil(t, err)

	// run a query
	var id int
	rows, err := db.Query("select 1 as id;")
	assert.Nil(t, err)

	// parse returned data
	rows.Next()
	err = rows.Scan(&id)
	assert.Nil(t, err)
	assert.Equal(t, id, 1)

	// destroy container
	close()

	// run a query again
	_, err = db.Query("select 1 as id;")
	assert.NotNil(t, err)
}

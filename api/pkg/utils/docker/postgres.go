package docker

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq" // postgres driver
)

// RunPostgres create and run postgres server.
// This method waits until the container is ready (container is up and can handle ping).
// It returns a dsn and function to destroy the container.
func RunPostgres(version string) (string, func()) {
	if version == "" {
		version = "latest"
	}
	var (
		dbname     = "somedb"
		dbuser     = "user"
		dbpassword = "password"
	)

	resource, err := startContainer("postgres", version, []string{
		"POSTGRES_DB=" + dbname,
		"POSTGRES_USER=" + dbuser,
		"POSTGRES_PASSWORD=" + dbpassword,
	})
	if err != nil {
		panic(fmt.Sprintf("fail to start postgres container, err: %s", err))
	}

	var dsn = fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		dbuser, dbpassword,
		resource.GetBoundIP("5432/tcp"), resource.GetPort("5432/tcp"),
		dbname,
	)
	checkFunc := func() error {
		db, _ := sql.Open("postgres", dsn)
		_, err := db.Exec("select 1;")
		if err == nil {
			db.Close()
			return nil
		}
		db.Close()
		return fmt.Errorf("fail to ping postgres, err: %s", err)
	}
	err = waitContainer(checkFunc)
	if err != nil {
		panic(fmt.Sprintf("got an error while waiting postgres to be ready, err: %s", err))
	}

	close := func() {
		resource.Close()
	}
	return dsn, close
}

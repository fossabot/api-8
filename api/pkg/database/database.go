package database

import (
	libSql "database/sql"
	"errors"
	"time"

	"github.com/rs/zerolog/log"
)

// error definitions
var (
	ErrTxFinished = errors.New("tx is already finished")
)

// sql exposes methods needed to execute query
type sql interface {
	query(dest interface{}, query string, args ...interface{}) error
	exec(query string, args ...interface{}) error
	newTransaction() Transaction
}

// Transaction represents an sql transaction.
// Transactions are always guaranteed to run in master connection.
type Transaction interface {
	Query(dest interface{}, query string, args ...interface{}) error
	Exec(query string, args ...interface{}) error
	Commit() error
	Rollback() error
}

// DBConf represents a database connection configuration
type DBConf struct {
	URL          string
	MaxOpenConns int
	MaxIdleConns int
	ConnLifetime time.Duration
}

// Config represents a configuration for this package
type Config struct {
	Master *DBConf
	Slaves []*DBConf
}

// Configure configures this package
func Configure(config *Config) error {
	masterConn, err := newConnection(config.Master.URL, config.Master.MaxOpenConns, config.Master.MaxIdleConns, config.Master.ConnLifetime)
	if err != nil {
		log.Error().Err(err).Str("dbURL", config.Master.URL)
		return err
	}

	slaveConns := make([]*connection, len(config.Slaves))
	for i, conf := range config.Slaves {
		conn, err := newConnection(conf.URL, conf.MaxOpenConns, conf.MaxIdleConns, conf.ConnLifetime)
		if err != nil {
			log.Error().Err(err).Str("dbURL", conf.URL)
			return err
		}
		slaveConns[i] = conn
	}

	cs = newConnectionSet(masterConn, slaveConns)
	return nil
}

// ConfigureTest configures this package for testing
func ConfigureTest(dsn string) *libSql.DB {
	master, err := newConnection(dsn, 1, 1, time.Minute)
	if err != nil {
		panic(err)
	}
	cs = newConnectionSet(master, nil)
	return master.db()
}

// Shutdown stops all connections
func Shutdown() {
	cs.quit()
}

package database

import (
	libSql "database/sql"
	"errors"
	"net/url"
	"sync"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres" // postgres dialect
	"github.com/sirupsen/logrus"
)

var (
	connTickDelay = 15 * time.Second
	pingTimeout   = 5 * time.Second
)

var errPingTimeout = errors.New("ping timeout")

// connection abstracts connection to a database server
// it handles closing connection via close() method
// it also handles connection updates by pinging the server every connTickDelay
type connection struct {
	host   string    // connection host for logging purposes
	s      sql       // underlying sql for executing query
	quitCh chan bool // chan for closing this connection

	cMutex    sync.RWMutex
	connected bool // wether we can connect using this connection or not

	pingFn  func() error // func used for pinging the server
	closeFn func()       // func used for closing the connection

	rawDB *libSql.DB // for testing
}

// create a new connection instance
// and start loop in background to update connection status
func newConnection(connStr string, maxOpenConns, maxIdleConns int, connLifetime time.Duration) (*connection, error) {
	u, err := url.Parse(connStr)
	if err != nil {
		logrus.WithError(err).WithField("connectionString", connStr).Errorln("error parsing connection string")
		return nil, err
	}

	db, err := gorm.Open("postgres", connStr)
	if err != nil {
		logrus.WithError(err).WithField("connectionString", connStr).Errorln("error openning database connection")
		return nil, err
	}
	db.DB().SetMaxOpenConns(maxOpenConns)
	db.DB().SetMaxIdleConns(maxIdleConns)
	db.DB().SetConnMaxLifetime(connLifetime)

	pingFn := wrapPingFn(func() error {
		return db.Raw("select 1").Error
	})
	closeFn := func() {
		db.Close()
	}
	conn := &connection{
		host:    u.Host,
		s:       newGormSQL(u.Host, db),
		quitCh:  make(chan bool),
		pingFn:  pingFn,
		closeFn: closeFn,
		rawDB:   db.DB(),
	}

	// updateStatus for the first time
	conn.updateStatus()

	// run loop to update connection status in background
	go conn.loop()

	logrus.WithField("db_host", conn.host).WithField("connected", conn.connected).Infoln("starting database connection")
	return conn, nil
}

func (c *connection) loop() {
	// safety check to prevent multiple ping / updateStatus running at the same time
	var running bool

	ticker := time.NewTicker(connTickDelay)
	for {
		select {
		case <-ticker.C:
			if !running {
				running = true
				c.updateStatus()
				running = false
			}

		case <-c.quitCh:
			logrus.WithField("host", c.host).Warnln("got quit signal, stopping database connection")

			// closeFn could be nil for testing purposes
			if c.closeFn != nil {
				c.closeFn()
			}
			c.connected = false // for now quit means app exists, so no need to secure this with mutex
			ticker.Stop()
			return
		}
	}
}

func (c *connection) updateStatus() {
	old := c.connected

	err := c.pingFn()
	c.cMutex.Lock()
	c.connected = err == nil
	c.cMutex.Unlock()

	if c.connected != old {
		logrus.WithField("host", c.host).WithField("connected", c.connected).Warnln("database connection status changed")
	}
	if err != nil {
		logrus.WithError(err).WithField("host", c.host).WithField("connected", c.connected).Errorln("error while checking database connection status")
	}
}

func (c *connection) quit() {
	close(c.quitCh)
}

func (c *connection) db() *libSql.DB {
	return c.rawDB
}

func wrapPingFn(f func() error) func() error {
	return func() error {
		errCh := make(chan error, 1)
		go func() {
			errCh <- f()
		}()
		select {
		case <-time.After(pingTimeout):
			return errPingTimeout
		case err := <-errCh:
			return err
		}
	}
}

package database

import (
	"errors"
	"fmt"
	"runtime"

	"github.com/jinzhu/gorm"
)

// connectionSet instance used for doing queries
var cs *connectionSet

var (
	errNoCsConfigured    = errors.New("no connection set configured")
	errNoReaderAvailable = errors.New("no reader connection available")
	errNoWriterAvailable = errors.New("no writer connection available")
)

// Query runs query to one of randomly-picked read replica
// If there is not read replica available, the query will run on writer
// Unlike exec, query could retrieve returned rows from query result
// dest can be array or single struct
// If dest is a struct, it will have the value of the first returned row
// If it is an array, it will have the value of all returned rows
// MAKE SURE NOT TO RUN ANY WRITE (INSERT OR UPDATE) QUERY USING THIS METHOD
func Query(dest interface{}, query string, args ...interface{}) error {
	if cs == nil {
		return errNoCsConfigured
	}
	conn := cs.reader()
	if conn == nil {
		return errNoReaderAvailable
	}
	return conn.query(dest, query, args...)
}

// WriterExec runs a query to writer
// It just runs the query without returning any data
func WriterExec(query string, args ...interface{}) error {
	if cs == nil {
		return errNoCsConfigured
	}
	conn := cs.writer()
	if conn == nil {
		return errNoWriterAvailable
	}
	return conn.exec(query, args...)
}

// WriterQuery runs query to writer
// dest can be array or single struct
// If dest is a struct, it will have the value of the first returned row
// If it is an array, it will have the value of all returned rows
// USE THIS METHOD WHEN YOU NEED TO QUERY DATA, RELATED TO JUST INSERTED DATA
// OTHER THAN THIS CASE, USE QUERY METHOD
func WriterQuery(dest interface{}, query string, args ...interface{}) error {
	if cs == nil {
		return errNoCsConfigured
	}
	conn := cs.writer()
	if conn == nil {
		return errNoReaderAvailable
	}
	return conn.query(dest, query, args...)
}

// NewTransaction creates a new database transaction
// This method guaratees that the transaction will be run on writer
func NewTransaction() (Transaction, error) {
	if cs == nil {
		return nil, errNoCsConfigured
	}
	conn := cs.writer()
	if conn == nil {
		return nil, errNoWriterAvailable
	}
	return cs.writer().newTransaction(), nil
}

// IgnoreNoRowsErr ignore if err equals sql.ErrNoRows
func IgnoreNoRowsErr(err error) error {
	if err == gorm.ErrRecordNotFound {
		return nil
	}
	return err
}

// NoRowsErr checks wether err is similar with sql.ErrNoRows
func NoRowsErr(err error) bool {
	return err == gorm.ErrRecordNotFound
}

func injectCallerInfo(sql string) string {
	pc, file, line, ok := runtime.Caller(3)
	details := runtime.FuncForPC(pc)
	if !ok || details == nil {
		msg := "/* failed to get caller info */"
		return fmt.Sprintf("%s\n%s", msg, sql)
	}

	msg := fmt.Sprintf("/* %s at %s:%d */", details.Name(), file, line)
	return fmt.Sprintf("%s\n%s", msg, sql)
}

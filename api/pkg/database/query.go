package database

import (
	"github.com/jinzhu/gorm"
)

// connectionSet instance used for doing queries
var cs *connectionSet

// Query runs query to one of randomly-picked read replica
// If there is not read replica available, the query will run on writer
// Unlike exec, query could retrieve returned rows from query result
// dest can be array or single struct
// If dest is a struct, it will have the value of the first returned row
// If it is an array, it will have the value of all returned rows
// MAKE SURE NOT TO RUN ANY WRITE (INSERT OR UPDATE) QUERY USING THIS METHOD
func Query(dest interface{}, query string, args ...interface{}) error {
	return cs.reader().query(dest, query, args...)
}

// WriterExec runs a query to writer
// It just runs the query without returning any data
func WriterExec(query string, args ...interface{}) error {
	return cs.writer().exec(query, args...)
}

// WriterQuery runs query to writer
// dest can be array or single struct
// If dest is a struct, it will have the value of the first returned row
// If it is an array, it will have the value of all returned rows
// USE THIS METHOD WHEN YOU NEED TO QUERY DATA, RELATED TO JUST INSERTED DATA
// OTHER THAN THIS CASE, USE QUERY METHOD
func WriterQuery(dest interface{}, query string, args ...interface{}) error {
	return cs.writer().query(dest, query, args...)
}

// NewTransaction creates a new database transaction
// This method guaratees that the transaction will be run on writer
func NewTransaction() Transaction {
	return cs.writer().newTransaction()
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

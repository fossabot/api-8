package database

import (
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
)

func newGormSQL(host string, db *gorm.DB) sql {
	return &gormSQL{
		host: host,
		db:   db,
	}
}

type gormSQL struct {
	db   *gorm.DB
	host string
}

func (g *gormSQL) query(dest interface{}, query string, args ...interface{}) error {
	query = injectCallerInfo(query)
	err := g.db.Raw(query, args...).Scan(dest).Error
	if err != nil {
		logrus.WithError(err).WithField("query", query).WithField("args", args).Errorln("a query returned error")
	}
	return err
}

func (g *gormSQL) exec(query string, args ...interface{}) error {
	query = injectCallerInfo(query)
	err := g.db.Exec(query, args...).Error
	if err != nil {
		logrus.WithError(err).WithField("host", g.host).WithField("query", query).WithField("args", args).Errorln("a query returned error")
	}
	return err
}

func (g *gormSQL) newTransaction() Transaction {
	return &gormTransaction{
		db:   g.db.Begin(),
		host: g.host,
	}
}

type gormTransaction struct {
	db       *gorm.DB
	host     string
	finished bool
}

func (g *gormTransaction) Query(dest interface{}, query string, args ...interface{}) error {
	if g.finished {
		return ErrTxFinished
	}
	query = injectCallerInfo(query)
	err := g.db.Raw(query, args...).Scan(dest).Error
	if err != nil {
		logrus.WithError(err).WithField("host", g.host).WithField("query", query).WithField("args", args).Errorln("a query returned error")
	}
	return err
}

func (g *gormTransaction) Exec(query string, args ...interface{}) error {
	if g.finished {
		return ErrTxFinished
	}
	query = injectCallerInfo(query)
	err := g.db.Exec(query, args...).Error
	if err != nil {
		logrus.WithError(err).WithField("host", g.host).WithField("query", query).WithField("args", args).Errorln("a query returned error")
	}
	return err
}

func (g *gormTransaction) Commit() error {
	if g.finished {
		return ErrTxFinished
	}
	g.finished = true
	err := g.db.Commit().Error
	if err != nil {
		logrus.WithError(err).WithField("host", g.host).Errorln("transaction commit failed")
	}
	return err
}

func (g *gormTransaction) Rollback() error {
	if g.finished {
		return ErrTxFinished
	}
	g.finished = true
	err := g.db.Rollback().Error
	if err != nil {
		logrus.WithError(err).WithField("host", g.host).Errorln("transaction rollback failed")
	}
	return err
}

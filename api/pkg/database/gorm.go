package database

import (
	"github.com/jinzhu/gorm"
	"github.com/rs/zerolog/log"
)

func newGormSQL(db *gorm.DB) sql {
	return &gormSQL{
		db: db,
	}
}

type gormSQL struct {
	db *gorm.DB
}

func (g *gormSQL) query(dest interface{}, query string, args ...interface{}) error {
	query = injectCallerInfo(query)
	err := g.db.Raw(query, args...).Scan(dest).Error
	if err != nil {
		log.Error().Err(err).Str("query", query).Interface("args", args)
	}
	return err
}

func (g *gormSQL) exec(query string, args ...interface{}) error {
	query = injectCallerInfo(query)
	err := g.db.Exec(query, args...).Error
	if err != nil {
		log.Error().Err(err).Str("query", query).Interface("args", args)
	}
	return err
}

func (g *gormSQL) newTransaction() Transaction {
	return &gormTransaction{
		db: g.db.Begin(),
	}
}

type gormTransaction struct {
	db       *gorm.DB
	finished bool
}

func (g *gormTransaction) Query(dest interface{}, query string, args ...interface{}) error {
	if g.finished {
		return ErrTxFinished
	}
	query = injectCallerInfo(query)
	err := g.db.Raw(query, args...).Scan(dest).Error
	if err != nil {
		log.Error().Err(err).Str("query", query).Interface("args", args)
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
		log.Error().Err(err).Str("query", query).Interface("args", args)
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
		log.Error().Err(err)
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
		log.Error().Err(err)
	}
	return err
}

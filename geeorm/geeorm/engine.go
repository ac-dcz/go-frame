package geeorm

import (
	"database/sql"
	"fmt"
	"geeorm/dialect"
	"geeorm/log"
	"geeorm/session"
)

type Engine struct {
	db      *sql.DB
	dialect dialect.Dialect
}

func NewEngine(driver, dsn string) (*Engine, error) {
	db, err := sql.Open(driver, dsn)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	if d, ok := dialect.GetDialect(driver); !ok {
		log.Errorf("Not Registry driver dialect %s", driver)
		return nil, fmt.Errorf("not registry driver dialect %s", driver)
	} else {
		return &Engine{db, d}, nil
	}
}

func (e *Engine) Close() error {
	err := e.db.Close()
	if err != nil {
		log.Error(err)
	}
	return err
}

func (e *Engine) NewSession() *session.Session {
	return session.New(e.db, e.dialect)
}

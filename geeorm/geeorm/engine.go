package geeorm

import (
	"database/sql"
	"geeorm/log"
	"geeorm/session"
)

type Engine struct {
	db *sql.DB
}

func NewEngine(driver, dsn string) (*Engine, error) {
	db, err := sql.Open(driver, dsn)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	return &Engine{db}, nil
}

func (e *Engine) Close() error {
	err := e.db.Close()
	if err != nil {
		log.Error(err)
	}
	return err
}

func (e *Engine) NewSession() *session.Session {
	return session.New(e.db)
}

package session

import (
	"database/sql"
	"geeorm/log"
	"strings"
)

type Session struct {
	db    *sql.DB
	sql   strings.Builder
	argvs []any
}

func New(db *sql.DB) *Session {
	return &Session{
		db:    db,
		sql:   strings.Builder{},
		argvs: nil,
	}
}

func (s *Session) Clear() {
	s.sql.Reset()
	s.argvs = nil
}

func (s *Session) DB() *sql.DB {
	return s.db
}

func (s *Session) Raw(sql string, argv ...any) *Session {
	s.sql.WriteString(sql)
	s.sql.WriteString(" ")
	s.argvs = append(s.argvs, argv...)
	return s
}

func (s *Session) Exec() (sql.Result, error) {
	defer s.Clear()
	return s.db.Exec(s.sql.String(), s.argvs...)
}

func (s *Session) QueryRow() *sql.Row {
	defer s.Clear()
	log.Info(s.sql.String(), s.argvs)
	return s.db.QueryRow(s.sql.String(), s.argvs...)
}

func (s *Session) QueryRows() (*sql.Rows, error) {
	defer s.Clear()
	log.Info(s.sql.String(), s.argvs)
	rows, err := s.db.Query(s.sql.String(), s.argvs...)
	if err != nil {
		log.Error(err)
	}
	return rows, err
}

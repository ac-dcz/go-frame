package session

import (
	"database/sql"
	"fmt"
	"geeorm/dialect"
	"geeorm/log"
	"geeorm/schema"
	"reflect"
	"strings"
)

type Session struct {
	db      *sql.DB
	sql     strings.Builder
	argvs   []any
	model   *schema.Schema
	dialect dialect.Dialect
}

func New(db *sql.DB, dialect dialect.Dialect) *Session {
	return &Session{
		db:      db,
		sql:     strings.Builder{},
		argvs:   nil,
		dialect: dialect,
	}
}

func (s *Session) Clear() {
	s.sql.Reset()
	s.argvs = nil
}

func (s *Session) DB() *sql.DB {
	return s.db
}

func (s *Session) SetModel(obj any) *Session {
	if s.model == nil || reflect.TypeOf(obj) != reflect.TypeOf(s.model.Model) {
		m, err := schema.NewSchema(obj)
		if err != nil {
			log.Error(err)
			return s
		}
		s.model = m
	}
	return s
}

func (s *Session) IsModel() bool {
	return s.model != nil
}

func (s *Session) CreateTable() (sql.Result, error) {
	if !s.IsModel() {
		log.Error("Not Set Model")
		return nil, fmt.Errorf("not set model")
	}
	var lines []string
	for i := range s.model.FieldNames {
		field := s.model.Fields[i]
		line := fmt.Sprintf("%s %s %s", field.Name, s.dialect.DataTypeOf(field.Typ), field.Tag)
		lines = append(lines, line)
	}
	cmd := fmt.Sprintf("create table %s(%s)", s.model.Name, strings.Join(lines, ","))

	return s.Raw(cmd).Exec()
}

func (s *Session) DropTable() (sql.Result, error) {
	if !s.IsModel() {
		log.Error("Not Set Model")
		return nil, fmt.Errorf("not set model")
	}
	cmd := fmt.Sprintf("drop table if exists %s", s.model.Name)
	return s.Raw(cmd).Exec()
}

func (s *Session) HasTable() error {
	if !s.IsModel() {
		log.Error("Not Set Model")
		return fmt.Errorf("not set model")
	}
	cmd := s.dialect.TableExistsSql(s.model.Name)
	row := s.Raw(cmd).QueryRow()
	if err := row.Err(); err != nil {
		return err
	} else {
		var name string
		err = row.Scan(&name)
		if err != nil {
			return err
		}
		return nil
	}
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

package main

import (
	"fmt"
	"geeorm"
	"geeorm/log"

	_ "github.com/go-sql-driver/mysql"
)

type DBConfig struct {
	User   string
	PassWD string
	Host   string
	DBname string
}

func (cfg *DBConfig) DSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s)/%s", cfg.User, cfg.PassWD, cfg.Host, cfg.DBname)
}

type User struct {
	Id   int    `geeorm:"primary key not null"`
	Name string `geeorm:"not null"`
	Age  int
}

func main() {
	cfg := &DBConfig{
		User:   "root",
		PassWD: "dcz.20001018",
		Host:   "127.0.0.1:3306",
		DBname: "study",
	}
	engine, _ := geeorm.NewEngine("mysql", cfg.DSN())
	defer engine.Close()
	s := engine.NewSession()
	s.SetModel(&User{})
	if err := s.HasTable(); err != nil {
		log.Error(err)
	} else {
		log.Infof("table %s is exists", "User")
	}

	res, err := s.DropTable()
	if err != nil {
		log.Error(err)
	} else {
		num, _ := res.RowsAffected()
		log.Infof("drop affect rows %d\n", num)
	}

	res, err = s.CreateTable()
	if err != nil {
		log.Error(err)
	} else {
		num, _ := res.RowsAffected()
		log.Infof("create affect rows %d\n", num)
	}
}

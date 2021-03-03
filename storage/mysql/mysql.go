package mysql

import (
	"fmt"
	"time"
	"trace"

	_ "github.com/go-sql-driver/mysql"

	"github.com/jinzhu/gorm"
)

var db *gorm.DB

type Mysql struct{}

func setDB() (err error) {
	if db != nil {
		if err = db.DB().Ping(); err == nil {
			return nil
		}
		db.DB().Close()
		db = nil
	}

	hostname := trace.Config.Mysql.Hostname
	port := trace.Config.Mysql.Port
	username := trace.Config.Mysql.Username
	password := trace.Config.Mysql.Password
	network := trace.Config.Mysql.Network
	database := trace.Config.Mysql.Database

	db_desc := fmt.Sprintf("%v:%v@%v(%v:%v)/%v",
		username, password, network, hostname, port, database)

	db, err = gorm.Open("mysql", db_desc)
	if err != nil {
		return err
	}

	err = db.DB().Ping()
	if err != nil {
		return err
	}

	db.DB().SetConnMaxLifetime(
		time.Duration(trace.Config.Mysql.ConnMaxLifeTime) * time.Second)
	db.DB().SetMaxIdleConns(trace.Config.Mysql.MaxIdleConns)
	db.DB().SetMaxOpenConns(trace.Config.Mysql.MaxOpenConns)

	return nil
}

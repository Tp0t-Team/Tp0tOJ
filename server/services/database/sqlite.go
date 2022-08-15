//go:build DatabaseSqlite

package database

import (
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"log"
	"os"
)

func InitDB(dsn string) {
	prefix, _ := os.Getwd()
	dbPath := prefix + "/resources/data.db"
	test, err := os.Lstat(dbPath)
	if os.IsNotExist(err) {
		_, err := os.Create(dbPath)
		if err != nil {
			log.Panicln(err, test)
			return
		}
	} else if err != nil {
		if err != nil {
			log.Panicln(err, test)
			return
		}
	}
	DataBase, err = gorm.Open(sqlite.Open("file://"+dbPath+"?cache=shared&_pragma=journal_mode(WAL)&_pragma=busy_timeout(10000)"), &gorm.Config{})
	if err != nil {
		log.Panicln("DB connect error", err.Error())
	}
	createTables()
	createDefaultUser()
}

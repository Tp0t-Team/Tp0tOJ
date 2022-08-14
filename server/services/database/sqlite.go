//go:build DatabaseSqlite

package database

import (
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"log"
	"os"
	"server/entity"
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
	DataBase, err = gorm.Open(sqlite.Open(dbPath+"?_pragma=busy_timeout%3d1000"), &gorm.Config{})
	if err != nil {
		log.Panicln("DB connect error", err.Error())
	}
	err = DataBase.AutoMigrate(&entity.Bulletin{}, &entity.Challenge{}, &entity.Replica{}, &entity.ReplicaAlloc{}, &entity.ResetToken{}, &entity.Submit{}, &entity.User{}, &entity.Behavior{}, &entity.GameEvent{})
	if err != nil {
		log.Panicln("DB connect error", err.Error())
		return
	}
	createDefaultUser()
}

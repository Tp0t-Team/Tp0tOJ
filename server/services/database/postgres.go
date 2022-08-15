//go:build DatabasePostgres

package database

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

func InitDB(dsn string) {
	var err error
	DataBase, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Panicln("DB connect error", err.Error())
	}
	createTables()
	createDefaultUser()
}

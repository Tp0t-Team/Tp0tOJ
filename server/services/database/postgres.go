//go:build DatabasePostgres

package database

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"server/entity"
)

func InitDB(dsn string) {
	var err error
	DataBase, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
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

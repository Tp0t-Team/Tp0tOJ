package utils

import (
	//"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log"
	"server/entity"
)

func init() {
	//dsn := "host=localhost user=gorm password=gorm dbname= port=9920 sslmode=disable TimeZone=Asia/Shanghai"
	db, err := gorm.Open(sqlite.Open("/Users/lordcasser/Workspace/Tp0tOJ/server/resources/test.db"), &gorm.Config{})
	if err != nil {
		log.Panicln("DB connect error", err.Error())
	}
	err = db.AutoMigrate(&entity.Bulletin{}, &entity.Challenge{}, &entity.Replica{}, &entity.ReplicaAlloc{}, &entity.ResetToken{}, &entity.Submit{}, &entity.User{})
	if err != nil {
		log.Panicln("DB connect error", err.Error())
		return
	}
}

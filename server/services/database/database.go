package database

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log"
	"os"
	"server/entity"
	"server/utils/configure"
	"time"
)

var DataBase *gorm.DB

//any database error log should be handle and log out inside resolvers, should not return to caller
func init() {
	needInit := false
	prefix, _ := os.Getwd()
	dbPath := prefix + "/resources/data.db"
	test, err := os.Lstat(dbPath)
	if os.IsNotExist(err) {
		_, err := os.Create(dbPath)
		needInit = true
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
	DataBase, err = gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		log.Panicln("DB connect error", err.Error())
	}
	err = DataBase.AutoMigrate(&entity.Bulletin{}, &entity.Challenge{}, &entity.Replica{}, &entity.ReplicaAlloc{}, &entity.ResetToken{}, &entity.Submit{}, &entity.User{})
	if err != nil {
		log.Panicln("DB connect error", err.Error())
		return
	}
	if needInit {
		defaultUser := entity.User{
			Name:     configure.Configure.Server.Username,
			Password: configure.Configure.Server.Password,
			State:    "normal",
			Mail:     configure.Configure.Server.Mail,
			JoinTime: time.Now(),
			Role:     "admin",
		}
		result := DataBase.Create(&defaultUser)
		if result.Error != nil {
			log.Panicln(result)
		}
	}
}

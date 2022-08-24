package database

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"io"
	"log"
	"server/entity"
	"server/utils/configure"
	"time"
)

var DataBase *gorm.DB

func passwordHash(password string) string {
	hash1 := sha256.New()
	_, err := io.WriteString(hash1, configure.Configure.Server.Salt+password)
	if err != nil {
		log.Panicln(err.Error())
	}
	hash2 := sha256.New()
	_, err = io.WriteString(hash2, configure.Configure.Server.Salt+fmt.Sprintf("%02x", hash1.Sum(nil)))
	if err != nil {
		log.Panicln(err.Error())
	}
	return fmt.Sprintf("%02x", hash2.Sum(nil))
}

//any database error log should be handle and log out inside resolvers, should not return to caller

func createTables() {
	err := DataBase.AutoMigrate(&entity.Bulletin{}, &entity.Challenge{}, &entity.Replica{}, &entity.ReplicaAlloc{}, &entity.ResetToken{}, &entity.Submit{}, &entity.User{}, &entity.Behavior{}, &entity.GameEvent{})
	if err != nil {
		log.Panicln("DB connect error", err.Error())
		return
	}
	if configure.Configure.Server.Debug.DBOpDetail {
		DataBase = DataBase.Debug()
	}
}

func createDefaultUser() {
	var users []entity.User
	result := DataBase.Find(&users)
	if result.Error != nil && !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		log.Panicln(result.Error)
	}
	if len(users) == 0 {
		defaultUser := entity.User{
			Name:     configure.Configure.Server.Username,
			Password: passwordHash(configure.Configure.Server.Password),
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

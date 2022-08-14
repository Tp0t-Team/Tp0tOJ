package database

import (
	"crypto/sha256"
	"fmt"
	"gorm.io/gorm"
	"io"
	"log"
	"os"
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

func createDefaultUser() {
	needInit := false
	prefix, _ := os.Getwd()
	lockPath := prefix + "/resources/dbInit.lock"
	test, err := os.Lstat(lockPath)
	if os.IsNotExist(err) {
		_, err := os.Create(lockPath)
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
	if needInit {
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

package entity

import (
	"crypto/md5"
	"fmt"
	"io"
	"log"
	"strings"
	"time"
)

type User struct {
	UserId    uint64    `gorm:"primaryKey"`
	Name      string    `gorm:"check: name <> ''"`
	Password  string    `gorm:"check: password <> ''"`
	State     string    `gorm:"check: state <> ''"` //normal,disabled
	Mail      string    `gorm:"check: mail <> ''"`
	JoinTime  time.Time `gorm:"not null"`
	Role      string    `gorm:"check: role <> ''"`
	CreatedAt time.Time `gorm:"autoCreateTime;not null;"`
	UpdatedAt time.Time `gorm:"autoUpdateTime;not null;"`
}

func (t User) MakeAvatarUrl() string {
	hash := md5.New()
	// TODO: user info leakage risk
	_, err := io.WriteString(hash, strings.ToLower(t.Mail))
	if err != nil {
		log.Fatalln("Can not calculate MD5 for making AvatarUrl")
		return "https://www.gravatar.com/avatar/"
	}
	result := fmt.Sprintf("https://www.gravatar.com/avatar/%x?d=404", hash.Sum(nil))
	return result
}

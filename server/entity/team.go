package entity

import (
	"crypto/md5"
	"fmt"
	"io"
	"log"
	"strings"
	"time"
)

type Team struct {
	TeamId      uint64    `gorm:"primaryKey"`
	Name        string    `gorm:"check: name <> ''"`
	Owner       uint64    `gorm:"not null"`
	Members     []uint64  `gorm:"not null"`
	Token       string    `gorm:"check: token <> ''"`
	State       string    `gorm:"check: state <> ''"` //normal,disabled
	CreatedTime time.Time `gorm:"not null"`
	CreatedAt   time.Time `gorm:"autoCreateTime;not null;"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime;not null;"`
}

func (t Team) MakeAvatarUrl() string {
	hash := md5.New()
	_, err := io.WriteString(hash, strings.ToLower(t.Name))
	if err != nil {
		log.Fatalln("Can not calculate MD5 for making AvatarUrl")
		return "https://cravatar.cn/avatar/"
	}
	result := fmt.Sprintf("https://cravatar.cn/avatar/%02x?d=retro", hash.Sum(nil))
	return result
}

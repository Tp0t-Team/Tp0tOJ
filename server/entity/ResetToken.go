package entity

import "time"

type ResetToken struct {
	TokenId   uint64    `gorm:"primaryKey"`
	Token     string    `gorm:"check: token <> ''"`
	User      User      `gorm:"foreignKey:User;not null;"`
	CreatedAt time.Time `gorm:"autoCreateTime;not null;"`
	UpdatedAt time.Time `gorm:"autoUpdateTime;not null;"`
}

package entity

import "time"

type ResetToken struct {
	TokenId   uint64    `gorm:"primaryKey"`
	Token     string    `gorm:"check: token <> ''"`
	UserId    uint64    `gorm:"not null"`
	User      User      `gorm:"references:UserId;"`
	CreatedAt time.Time `gorm:"autoCreateTime;not null;"`
	UpdatedAt time.Time `gorm:"autoUpdateTime;not null;"`
}

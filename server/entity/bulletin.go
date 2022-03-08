package entity

import "time"

type Bulletin struct {
	BulletinId  uint64 `gorm:"primaryKey"`
	Content     string `gorm:"check: content <> ''"`
	Topping     bool
	CreatedAt   time.Time `gorm:"autoCreateTime;not null;"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime;not null;"`
	Title       string    `gorm:"not null"`
	PublishTime time.Time `gorm:"not null"`
}

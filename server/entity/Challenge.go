package entity

import "time"

type Challenge struct {
	ChallengeId   uint64    `gorm:"primaryKey"`
	Configuration string    `gorm:"check: configuration <> ''"`
	FirstBlood    User      `gorm:"foreignKey:User"`
	SecondBlood   User      `gorm:"foreignKey:User"`
	ThirdBlood    User      `gorm:"foreignKey:User"`
	State         string    `gorm:"check: state <> ''"`
	CreatedAt     time.Time `gorm:"autoCreateTime;not null;"`
	UpdatedAt     time.Time `gorm:"autoUpdateTime;not null;"`
}

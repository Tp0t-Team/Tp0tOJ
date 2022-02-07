package entity

import "time"

type Challenge struct {
	ChallengeId   uint64 `gorm:"primaryKey"`
	Configuration string `gorm:"check: configuration <> ''"`
	FirstBloodId  uint64
	SecondBloodId uint64
	ThirdBloodId  uint64
	FirstBlood    User      `gorm:"foreignKey:FirstBloodId"`
	SecondBlood   User      `gorm:"foreignKey:SecondBloodId"`
	ThirdBlood    User      `gorm:"foreignKey:ThirdBloodId"`
	State         string    `gorm:"check: state <> ''"`
	CreatedAt     time.Time `gorm:"autoCreateTime;not null;"`
	UpdatedAt     time.Time `gorm:"autoUpdateTime;not null;"`
}

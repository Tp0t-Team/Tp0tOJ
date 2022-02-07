package entity

import "time"

type Submit struct {
	SubmitId    uint64    `gorm:"primaryKey"`
	UserId      uint64    `gorm:"not null"`
	User        User      `gorm:"foreignKey:UserId"`
	Mark        int64     `gorm:"not null"`
	SubmitTime  time.Time `gorm:"not null"`
	Flag        string    `gorm:"check: flag <> ''"`
	Correct     bool
	CreatedAt   time.Time `gorm:"autoCreateTime;not null;"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime;not null;"`
	ChallengeId uint64    `gorm:"not null"`
	Challenge   Challenge `gorm:"foreignKey:ChallengeId"`
}

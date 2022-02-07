package entity

import "time"

type Submit struct {
	SubmitId   uint64    `gorm:"primaryKey"`
	User       User      `gorm:"foreignKey:User;not null;"`
	Mark       int64     `gorm:"not null"`
	SubmitTime time.Time `gorm:"not null"`
	Flag       string    `gorm:"check: flag <> ''"`
	Correct    bool
	CreatedAt  time.Time `gorm:"autoCreateTime;not null;"`
	UpdatedAt  time.Time `gorm:"autoUpdateTime;not null;"`
	Challenge  Challenge `gorm:"foreignKey:Challenge;not null;"`
}

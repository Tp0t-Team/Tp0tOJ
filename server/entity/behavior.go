package entity

import "time"

type Behavior struct {
	Id          uint64    `gorm:"primaryKey"`
	CreatedAt   time.Time `gorm:"autoCreateTime;not null;"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime;not null;"`
	ActionTime  time.Time `gorm:"not null"`
	Action      uint64    `gorm:"not null"`
	ChallengeId uint64
	UserId      uint64 `gorm:"not null"`
	Content     string // Flag or IP
}

const (
	ActionAllocReplica = iota
	ActionSubmit
	ActionComplete
	ActionLogin
	ActionWatchDescription
	ActionUploadWP
	BehaviorMAX
)

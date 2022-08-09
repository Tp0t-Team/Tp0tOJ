package entity

import "time"

type GameEvent struct {
	EventId   uint64    `gorm:"primaryKey"`
	CreatedAt time.Time `gorm:"autoCreateTime;not null;"`
	UpdatedAt time.Time `gorm:"autoUpdateTime;not null;"`
	Time      time.Time `gorm:"not null"`
	Action    uint64    `gorm:"not null"`
}

const (
	PauseEvent = iota
	ResumeEvent
)

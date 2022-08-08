package entity

import "time"

type Replica struct {
	ReplicaId   uint64    `gorm:"primaryKey"`
	ChallengeId uint64    `gorm:"not null"`
	Challenge   Challenge `gorm:"references:ChallengeId;"`
	Singleton   bool
	Status      string    `gorm:"check: status <> ''"`
	Flag        string    `gorm:"check: flag <> ''"`
	FlagType    uint64    `gorm:"not null"`
	CreatedAt   time.Time `gorm:"autoCreateTime;not null;"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime;not null;"`
}

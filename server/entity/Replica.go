package entity

import "time"

type Replica struct {
	ReplicaId uint64    `gorm:"primaryKey"`
	Challenge Challenge `gorm:"foreignKey:Challenge;not null;"`
	Flag      string    `gorm:"check: flag <> ''"`
	CreatedAt time.Time `gorm:"autoCreateTime;not null;"`
	UpdatedAt time.Time `gorm:"autoUpdateTime;not null;"`
}

package entity

type ReplicaAlloc struct {
	ReplicaAllocId uint64  `gorm:"primaryKey"`
	UserId         uint64  `gorm:"not null"`
	ReplicaId      uint64  `gorm:"not null"`
	User           User    `gorm:"foreignKey:UserId"`
	Replica        Replica `gorm:"foreignKey:ReplicaId"`
}

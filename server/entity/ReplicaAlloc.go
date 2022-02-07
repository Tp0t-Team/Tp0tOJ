package entity

type ReplicaAlloc struct {
	ReplicaAllocId uint64  `gorm:"primaryKey"`
	User           User    `gorm:"foreignKey:User;not null;"`
	Replica        Replica `gorm:"foreignKey:Replica;not null;"`
}

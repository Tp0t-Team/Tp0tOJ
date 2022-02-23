package resolvers

import (
	"errors"
	"gorm.io/gorm"
	"server/entity"
)

func FindReplicaAllocByUserId(userId uint64) ([]entity.ReplicaAlloc, error) {
	var replicaAllocs []entity.ReplicaAlloc
	result := db.Where(map[string]interface{}{"UserId": userId}).Find(&replicaAllocs)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return []entity.ReplicaAlloc{}, nil
	} else if result.Error != nil {
		return nil, result.Error
	}
	return replicaAllocs, nil
}

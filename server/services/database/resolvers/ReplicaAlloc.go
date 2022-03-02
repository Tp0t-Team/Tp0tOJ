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

func AddReplicaAlloc(replicaId uint64, userId uint64) error {
	// TODO: maybe need some more check, to ensure user only have one replica for each challenge
	replicaAlloc := entity.ReplicaAlloc{
		ReplicaId: replicaId,
		UserId:    userId,
	}
	result := db.Create(&replicaAlloc)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func DeleteReplicaAllocByReplicaId(replicaId uint64, outsideTX *gorm.DB) error {
	if outsideTX == nil {
		outsideTX = db
	}
	err := outsideTX.Transaction(func(tx *gorm.DB) error {
		var replicaAllocs []entity.ReplicaAlloc
		getResult := tx.Where(map[string]interface{}{"ReplicaId": replicaId}).Find(&replicaAllocs)
		if getResult.Error != nil {
			return getResult.Error
		}
		for _, replicaAlloc := range replicaAllocs {
			db.Delete(&replicaAlloc)
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

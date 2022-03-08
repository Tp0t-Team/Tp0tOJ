package resolvers

import (
	"errors"
	"gorm.io/gorm"
	"log"
	"server/entity"
)

func FindReplicaAllocByUserId(userId uint64) []entity.ReplicaAlloc {
	var replicaAllocs []entity.ReplicaAlloc
	result := db.Preload("User").Preload("Replica").Preload("Replica.Challenge").Where(map[string]interface{}{"user_id": userId}).Find(&replicaAllocs)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return []entity.ReplicaAlloc{}
	} else if result.Error != nil {
		log.Println(result.Error)
		return nil
	}
	return replicaAllocs
}

func FindReplicaAllocByUserIdAndChallengeId(userId uint64, challengeId uint64, outsideTX *gorm.DB) (*entity.ReplicaAlloc, error) {
	if outsideTX == nil {
		outsideTX = db
	}
	var found []entity.ReplicaAlloc
	err := db.Transaction(func(tx *gorm.DB) error {
		var allocs []entity.ReplicaAlloc
		allocResult := tx.Preload("User").Preload("Replica").Preload("Replica.Challenge").Where(map[string]interface{}{"user_id": userId}).Find(&allocs)
		if errors.Is(allocResult.Error, gorm.ErrRecordNotFound) {
			return nil
		} else if allocResult.Error != nil {
			return allocResult.Error
		}
		for _, alloc := range allocs {
			if alloc.Replica.ChallengeId == challengeId {
				found = append(found, alloc)
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	if found == nil {
		return nil, nil
	}
	if len(found) != 1 {
		return nil, errors.New("repeat replica alloc")
	}
	return &found[0], nil
}
func FindReplicaAllocByUserIdAndReplicaId(userId uint64, replicaId uint64, outsideTX *gorm.DB) (*entity.ReplicaAlloc, error) {
	if outsideTX == nil {
		outsideTX = db
	}
	var found []entity.ReplicaAlloc
	err := db.Transaction(func(tx *gorm.DB) error {
		var allocs []entity.ReplicaAlloc
		allocResult := tx.Preload("User").Preload("Replica").Preload("Replica.Challenge").Where(map[string]interface{}{"user_id": userId}).Find(&allocs)
		if errors.Is(allocResult.Error, gorm.ErrRecordNotFound) {
			return nil
		} else if allocResult.Error != nil {
			return allocResult.Error
		}
		for _, alloc := range allocs {
			if alloc.Replica.ReplicaId == replicaId {
				found = append(found, alloc)
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	if found == nil {
		return nil, nil
	}
	if len(found) != 1 {
		return nil, errors.New("repeat replica alloc")
	}
	return &found[0], nil
}
func AddReplicaAlloc(replicaId uint64, userId uint64, outsideTX *gorm.DB) bool {
	if outsideTX == nil {
		outsideTX = db
	}
	found, err := FindReplicaAllocByUserIdAndReplicaId(userId, replicaId, outsideTX)
	if err != nil {
		log.Println(err)
		return false
	}
	if found != nil {
		log.Println("recreate a repeat replicaAlloc error")
		return false
	}
	replicaAlloc := entity.ReplicaAlloc{
		ReplicaId: replicaId,
		UserId:    userId,
	}
	result := outsideTX.Create(&replicaAlloc)
	if result.Error != nil {
		log.Println(result.Error)
		return false
	}
	return true

}

func DeleteReplicaAllocByReplicaId(replicaId uint64, outsideTX *gorm.DB) bool {
	if outsideTX == nil {
		outsideTX = db
	}
	err := outsideTX.Transaction(func(tx *gorm.DB) error {
		var replicaAllocs []entity.ReplicaAlloc
		getResult := tx.Preload("User").Preload("Replica").Preload("Replica.Challenge").Where(map[string]interface{}{"replica_id": replicaId}).Find(&replicaAllocs)
		if getResult.Error != nil {
			return getResult.Error
		}
		for _, replicaAlloc := range replicaAllocs {
			tx.Delete(&replicaAlloc)
		}
		return nil
	})
	if err != nil {
		log.Println(err)
		return false
	}
	return true
}

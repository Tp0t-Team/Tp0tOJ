package resolvers

import (
	"errors"
	"gorm.io/gorm"
	"log"
	"server/entity"
)

func FindReplicaAllocByUserId(userId uint64) []entity.ReplicaAlloc {
	var replicaAllocs []entity.ReplicaAlloc
	result := db.Where(map[string]interface{}{"UserId": userId}).Find(&replicaAllocs)
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
		allocResult := tx.Where(map[string]interface{}{"UserId": userId}).Find(&allocs)
		if allocResult.Error != nil {
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
		allocResult := tx.Where(map[string]interface{}{"UserId": userId}).Find(&allocs)
		if allocResult.Error != nil {
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
	if len(found) != 1 {
		return nil, errors.New("repeat replica alloc")
	}
	return &found[0], nil
}
func AddReplicaAlloc(replicaId uint64, userId uint64, outsideTX *gorm.DB) bool {
	if outsideTX == nil {
		outsideTX = db
	}
	_, err := FindReplicaAllocByUserIdAndReplicaId(userId, replicaId, outsideTX)
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		log.Println("recreate a repeat replicaAlloc error")
		return false
	} else if err != nil {
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
		getResult := tx.Where(map[string]interface{}{"ReplicaId": replicaId}).Find(&replicaAllocs)
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

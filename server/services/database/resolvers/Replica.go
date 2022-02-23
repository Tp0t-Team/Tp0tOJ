package resolvers

import (
	"errors"
	"gorm.io/gorm"
	"server/entity"
)

func FindReplicaByChallengeId(challengeId uint64) ([]entity.Replica, error) {
	var replicas []entity.Replica
	result := db.Where(map[string]interface{}{"ChallengeId": challengeId}).Find(&replicas)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return []entity.Replica{}, nil
	} else if result.Error != nil {
		return nil, result.Error
	}
	return replicas, nil
}

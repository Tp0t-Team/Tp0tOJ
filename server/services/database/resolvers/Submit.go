package resolvers

import (
	"encoding/json"
	"errors"
	"gorm.io/gorm"
	"server/entity"
	"server/services/types"
	"time"
)

func CheckSubmitCorrectByUserIdAndChallengeId(userId uint64, challengeId uint64) (bool, error) {
	result := db.Where(map[string]interface{}{"UserId": userId, "ChallengeId": challengeId, "Correct": true, "Available": true}).Find(&entity.Submit{})
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return false, nil
	} else if result.Error != nil {
		return false, result.Error
	}
	return true, nil
}

func FindSubmitCorrectByChallengeId(challengeId uint64) ([]entity.Submit, error) {
	var submits []entity.Submit
	result := db.Where(map[string]interface{}{"ChallengeId": challengeId, "Correct": true, "Available": true}).Find(&submits)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return []entity.Submit{}, nil
	} else if result.Error != nil {
		return nil, result.Error
	}
	return submits, nil
}

func FindSubmitCorrectByUserId(userId uint64) ([]entity.Submit, error) {
	var submits []entity.Submit
	result := db.Where(map[string]interface{}{"UserId": userId, "Correct": true, "Available": true}).Find(&submits)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return []entity.Submit{}, nil
	} else if result.Error != nil {
		return nil, result.Error
	}
	return submits, nil
}

func FindSubmitCorrectSorted() ([]entity.Submit, error) {
	var submits []entity.Submit
	result := db.Where(map[string]interface{}{"Correct": true, "Available": true}).Order("SubmitTime").Find(&submits)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return []entity.Submit{}, nil
	} else if result.Error != nil {
		return nil, result.Error
	}
	return submits, nil
}

func AddSubmit(userId uint64, challengeId uint64, flag string, submitTime time.Time) error {
	err := db.Transaction(func(tx *gorm.DB) error {
		var allocs []entity.ReplicaAlloc
		allocResult := tx.Where(map[string]interface{}{"UserId": userId}).Find(&allocs)
		if allocResult.Error != nil {
			return allocResult.Error
		}
		for _, alloc := range allocs {
			if alloc.Replica.ChallengeId != challengeId {
				continue
			}
			var submits []entity.Submit
			result := tx.Where(map[string]interface{}{"ChallengeId": challengeId, "Correct": true, "Available": true}).Find(&submits)
			if result.Error != nil {
				return result.Error
			}
			var config types.ChallengeConfig
			err := json.Unmarshal([]byte(alloc.Replica.Challenge.Configuration), &config)
			if err != nil {
				return err
			}
			if alloc.Replica.ChallengeId == challengeId {
				newSubmit := entity.Submit{
					UserId:      userId,
					ChallengeId: challengeId,
					SubmitTime:  submitTime,
					Mark:        int64(len(submits)),
					Flag:        flag,
					Correct:     alloc.Replica.Flag == flag,
					Available:   config.State == "enabled",
				}
				tx.Create(&newSubmit)
				return nil
			}
		}
		return errors.New("no alloc exists")
	})
	if err != nil {
		return err
	}
	return nil
}

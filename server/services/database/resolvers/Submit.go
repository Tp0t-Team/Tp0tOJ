package resolvers

import (
	"encoding/json"
	"errors"
	"gorm.io/gorm"
	"log"
	"server/entity"
	"server/services/types"
	"time"
)

func CheckSubmitCorrectByUserIdAndChallengeId(userId uint64, challengeId uint64) bool {
	result := db.Where(map[string]interface{}{"UserId": userId, "ChallengeId": challengeId, "Correct": true, "Available": true}).Find(&entity.Submit{})
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return false
	} else if result.Error != nil {
		log.Println(result.Error)
		return false
	}
	return true
}

func FindSubmitCorrectByChallengeId(challengeId uint64) []entity.Submit {
	var submits []entity.Submit
	result := db.Where(map[string]interface{}{"ChallengeId": challengeId, "Correct": true, "Available": true}).Find(&submits)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return []entity.Submit{}
	} else if result.Error != nil {
		log.Println(result.Error)
		return nil
	}
	return submits
}

func FindSubmitCorrectByUserId(userId uint64) []entity.Submit {
	var submits []entity.Submit
	result := db.Where(map[string]interface{}{"UserId": userId, "Correct": true, "Available": true}).Find(&submits)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return []entity.Submit{}
	} else if result.Error != nil {
		log.Println(result.Error)
		return nil
	}
	return submits
}

func FindSubmitCorrectSorted() []entity.Submit {
	var submits []entity.Submit
	result := db.Where(map[string]interface{}{"Correct": true, "Available": true}).Order("SubmitTime").Find(&submits)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return []entity.Submit{}
	} else if result.Error != nil {
		log.Println(result.Error)
		return nil
	}
	return submits
}

func AddSubmit(userId uint64, challengeId uint64, flag string, submitTime time.Time) bool {
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
		log.Println(err)
		return false
	}
	return true

}

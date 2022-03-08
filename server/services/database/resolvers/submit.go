package resolvers

import (
	"encoding/json"
	"errors"
	"gorm.io/gorm"
	"log"
	"server/entity"
	"server/services/types"
	"server/utils"
	"time"
)

func CheckSubmitCorrectByUserIdAndChallengeId(userId uint64, challengeId uint64) bool {
	result := db.Where(map[string]interface{}{"user_id": userId, "challenge_id": challengeId, "correct": true, "available": true}).First(&entity.Submit{})
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return false
	} else if result.Error != nil {
		log.Println(result.Error)
		return false
	}
	return true
}

func FindAllSubmitByChallengeId(challengeId uint64) []entity.Submit {
	var submits []entity.Submit
	result := db.Preload("User").Preload("Challenge").Where(map[string]interface{}{"challenge_id": challengeId}).Find(&submits)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return []entity.Submit{}
	} else if result.Error != nil {
		log.Println(result.Error)
		return nil
	}
	return submits
}

func FindSubmitCorrectByChallengeId(challengeId uint64) []entity.Submit {
	var submits []entity.Submit
	result := db.Preload("User").Preload("Challenge").Where(map[string]interface{}{"challenge_id": challengeId, "correct": true, "available": true}).Find(&submits)
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
	result := db.Preload("User").Preload("Challenge").Where(map[string]interface{}{"user_id": userId, "correct": true, "available": true}).Find(&submits)
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
	result := db.Preload("User").Preload("Challenge").Where(map[string]interface{}{"correct": true, "available": true}).Order("submit_time").Find(&submits)
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
		alloc, err := FindReplicaAllocByUserIdAndChallengeId(userId, challengeId, tx)
		if err != nil {
			return err
		}
		if alloc == nil {
			return errors.New("no alloc exists")
		}
		var submits []entity.Submit
		result := tx.Preload("Challenge").Where(map[string]interface{}{"challenge_id": challengeId, "correct": true, "available": true}).Find(&submits)
		if result.Error != nil {
			return result.Error
		}
		var config types.ChallengeConfig
		err = json.Unmarshal([]byte(alloc.Replica.Challenge.Configuration), &config)
		if err != nil {
			return err
		}
		newSubmit := entity.Submit{
			UserId:      userId,
			ChallengeId: challengeId,
			SubmitTime:  submitTime,
			Mark:        int64(len(submits)),
			Flag:        flag,
			Correct:     alloc.Replica.Flag == flag,
			Available:   alloc.Replica.Challenge.State == "enabled",
		}
		tx.Create(&newSubmit)
		if newSubmit.Correct {
			utils.Cache.Submit(newSubmit.UserId, newSubmit.ChallengeId, submitTime)
		}
		return nil
	})
	if err != nil {
		log.Println(err)
		return false
	}
	return true

}

func DeleteSubmitsByChallengeId(challengeId uint64, outsideTX *gorm.DB) bool {
	if outsideTX == nil {
		outsideTX = db
	}
	err := outsideTX.Transaction(func(tx *gorm.DB) error {
		var submits []entity.Submit
		submits = FindAllSubmitByChallengeId(challengeId)
		if submits == nil {
			return errors.New("ChallengeId not find during Delete submits")
		}
		for _, submit := range submits {
			tx.Delete(&submit)
		}
		return nil
	})
	if err != nil {
		log.Println(err)
		return false
	}
	return true
}

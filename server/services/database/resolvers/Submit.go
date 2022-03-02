package resolvers

import (
	"errors"
	"gorm.io/gorm"
	"log"
	"server/entity"
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

func AddSubmit(userId uint64, challengeId uint64, flag string) bool {
	err := db.Transaction(func(tx *gorm.DB) error {
		// TODO:
		return nil
	})
	if err != nil {
		log.Println(err)
		return false
	}
	return true

}

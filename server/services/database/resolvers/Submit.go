package resolvers

import (
	"errors"
	"gorm.io/gorm"
	"server/entity"
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

func AddSubmit(userId uint64, challengeId uint64, flag string) error {
	err := db.Transaction(func(tx *gorm.DB) error {
		// TODO:
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

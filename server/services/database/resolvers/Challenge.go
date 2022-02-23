package resolvers

import (
	"encoding/json"
	"errors"
	"gorm.io/gorm"
	"log"
	"server/entity"
	"server/services/types"
	"strconv"
)

func FindChallengeByState(state string) ([]entity.Challenge, error) {
	var challenges []entity.Challenge
	result := db.Where(map[string]interface{}{"State": state}).Find(&challenges)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	} else if result.Error != nil {
		return nil, result.Error
	}
	return challenges, nil
}
func FindChallengeById(id uint64) (*entity.Challenge, error) {
	var challenge entity.Challenge
	result := db.Where(map[string]interface{}{"ChallengeId": id}).First(&challenge)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	} else if result.Error != nil {
		return nil, result.Error
	}
	return &challenge, nil
}

func AddChallenge(input types.ChallengeMutateInput) error {
	err := db.Transaction(func(tx *gorm.DB) error {
		//wo don't allow the same name between two challenges
		checkResult := tx.Where(map[string]interface{}{"Name": input.Name}).First(&entity.Challenge{})
		if errors.Is(checkResult.Error, gorm.ErrRecordNotFound) {
			newChallengeConfig := types.ChallengeConfig{Name: input.Name, Category: input.Category, Score: input.Score, Flag: input.Flag, Description: input.Description, ExternalLink: input.ExternalLink, Hint: input.Hint}
			marshalConfig, err := json.Marshal(newChallengeConfig)
			if err != nil {
				return err
			}
			newChallenge := entity.Challenge{Configuration: string(marshalConfig), State: input.State}
			result := tx.Create(&newChallenge)
			if result.Error != nil {
				return result.Error
			}
			return nil
		} else if checkResult.Error != nil {
			return checkResult.Error
		} else {
			return errors.New("database item exists")
		}
		// TODO:create replicas and allocate to all users
	})
	if err != nil {
		return err
	}

	return nil
}
func UpdateChallenge(input types.ChallengeMutateInput) error {
	err := db.Transaction(func(tx *gorm.DB) error {
		var challenge entity.Challenge
		inputChallengeId, err := strconv.ParseUint(input.ChallengeId, 10, 64)
		if err != nil {
			log.Println("challengeId parse error", err)
		}
		challengeItem := tx.Where(map[string]interface{}{"ChallengeId": inputChallengeId}).First(&challenge)
		if challengeItem.Error != nil {
			return errors.New("find challenge by ChallengeId error:\n" + challengeItem.Error.Error())
		}
		newChallengeConfig := types.ChallengeConfig{Name: input.Name, Category: input.Category, Score: input.Score, Flag: input.Flag, Description: input.Description, ExternalLink: input.ExternalLink, Hint: input.Hint}
		marshalConfig, err := json.Marshal(newChallengeConfig)
		if err != nil {
			return err
		}
		challenge.Configuration = string(marshalConfig)
		if input.State != "" {
			challenge.State = input.State
		}
		//TODO: wo don't allow the same name between two challenges
		//checkResult := tx.Where(map[string]interface{}{"Name": input.Name}).Find(&entity.Challenge{})
		db.Save(&challenge)
		// TODO: update flag replicas
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func FindEnabledChallenges() ([]entity.Challenge, error) {
	var challenges []entity.Challenge
	result := db.Where(map[string]interface{}{"State": "enabled"}).Find(&challenges)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	} else if result.Error != nil {
		return nil, result.Error
	}
	return challenges, nil
}

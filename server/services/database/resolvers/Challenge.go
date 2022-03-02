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

func FindChallengeByState(state string) []entity.Challenge {
	var challenges []entity.Challenge
	result := db.Where(map[string]interface{}{"State": state}).Find(&challenges)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return []entity.Challenge{}
	} else if result.Error != nil {
		log.Println(result.Error)
		return nil
	}
	return challenges
}

func FindChallengeById(id uint64) (*entity.Challenge, error) {
	var challenge entity.Challenge
	result := db.Where(map[string]interface{}{"ChallengeId": id}).First(&challenge)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	} else if result.Error != nil {
		//log.Println(result.Error)
		return nil, result.Error
	}
	return &challenge, nil
}

func FindAllChallenges() []entity.Challenge {
	var challenges []entity.Challenge
	result := db.Where(map[string]interface{}{}).Find(&challenges)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return []entity.Challenge{}
	} else if result.Error != nil {
		log.Println(result.Error)
		return nil
	}
	return challenges
}

func AddChallenge(input types.ChallengeMutateInput) bool {
	err := db.Transaction(func(tx *gorm.DB) error {
		//wo don't allow the same name between two challenges
		checkResult := tx.Where(map[string]interface{}{"Name": input.Name}).First(&entity.Challenge{})
		if errors.Is(checkResult.Error, gorm.ErrRecordNotFound) {
			newChallengeConfig := types.ChallengeConfig{
				Name:     input.Name,
				Category: input.Category,
				Score: types.ScoreType{
					Dynamic:   input.Score.Dynamic,
					BaseScore: input.Score.BaseScore,
				},
				Flag: types.FlagType{
					Dynamic: input.Flag.Dynamic,
					Value:   input.Flag.Value,
				},
				Description:  input.Description,
				ExternalLink: input.ExternalLink,
				Hint:         input.Hint,
			}
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
			return errors.New("database item challenge already exists")
		}
		// TODO: create replicas and allocate to all users if challenge is a singleton & enabled
		return nil
	})
	if err != nil {
		log.Println(err)
		return false
	}

	return true
}

func UpdateChallenge(input types.ChallengeMutateInput) bool {
	err := db.Transaction(func(tx *gorm.DB) error {
		var challenge entity.Challenge
		inputChallengeId, err := strconv.ParseUint(input.ChallengeId, 10, 64)
		if err != nil {
			return errors.New("challengeId Parse Error:\n" + err.Error())

		}
		challengeItem := tx.Where(map[string]interface{}{"ChallengeId": inputChallengeId}).First(&challenge)
		if challengeItem.Error != nil {
			return errors.New("find challenge by ChallengeId error:\n" + challengeItem.Error.Error())
		}
		newChallengeConfig := types.ChallengeConfig{
			Name:     input.Name,
			Category: input.Category,
			Score: types.ScoreType{
				Dynamic:   input.Score.Dynamic,
				BaseScore: input.Score.BaseScore,
			},
			Flag: types.FlagType{
				Dynamic: input.Flag.Dynamic,
				Value:   input.Flag.Value,
			},
			Description:  input.Description,
			ExternalLink: input.ExternalLink,
			Hint:         input.Hint,
		}
		marshalConfig, err := json.Marshal(newChallengeConfig)
		if err != nil {
			return err
		}
		challenge.Configuration = string(marshalConfig)
		if input.State != "" {
			challenge.State = input.State
		}
		//TODO: we don't allow the same name between two challenges
		//checkResult := tx.Where(map[string]interface{}{"Name": input.Name}).Find(&entity.Challenge{})
		db.Save(&challenge)
		// TODO: update flag replicas
		// TODO: if change state "disabled", replica delete & set all submits unavailable
		// TODO: if change dockerfile, replica re-create
		// TODO: you can't change flag dynamic-able
		// TODO: you can't change score dynamic-able
		return nil
	})
	if err != nil {
		log.Println(err)
		return false
	}
	return true
}

func FindEnabledChallenges() []entity.Challenge {
	var challenges []entity.Challenge
	result := db.Where(map[string]interface{}{"State": "enabled"}).Find(&challenges)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return []entity.Challenge{}
	} else if result.Error != nil {
		log.Println(result.Error)
		return nil
	}
	return challenges
}

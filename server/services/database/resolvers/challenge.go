package resolvers

import (
	"encoding/json"
	"errors"
	"gorm.io/gorm"
	"log"
	"server/entity"
	"server/services/types"
	"server/utils"
	"strconv"
)

func FindChallengeByState(state string) []entity.Challenge {
	var challenges []entity.Challenge
	result := db.Where(map[string]interface{}{"state": state}).Find(&challenges)
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
	result := db.Where(map[string]interface{}{"challenge_id": id}).First(&challenge)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	} else if result.Error != nil {
		//log.Println(result.Error)
		return nil, result.Error
	}
	return &challenge, nil
}

func FindChallengeByIdInTX(id uint64, outsideTX *gorm.DB) (*entity.Challenge, error) {
	if outsideTX == nil {
		outsideTX = db
	}
	var challenge entity.Challenge
	result := outsideTX.Where(map[string]interface{}{"challenge_id": id}).First(&challenge)
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

func FindChallengeByName(name string) (*entity.Challenge, error) {
	var challenge entity.Challenge
	result := db.Where(map[string]interface{}{"name": name}).First(&challenge)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	} else if result.Error != nil {
		//log.Println(result.Error)
		return nil, result.Error
	}
	return &challenge, nil
}
func AddChallenge(input types.ChallengeMutateInput) bool {
	//var challenge *entity.Challenge = nil
	err := db.Transaction(func(tx *gorm.DB) error {
		// we don't allow the same name between two challenges
		challenge, err := FindChallengeByName(input.Name)
		if err != nil {
			return err
		}
		if challenge != nil {
			return errors.New("can not add challenge,cause same name already existed")
		}

		nodes := []types.NodeConfig{}
		if input.NodeConfig != nil {
			for _, node := range *input.NodeConfig {
				nodes = append(nodes, node.ToNodeConfig())
			}
		}
		newChallengeConfig := types.ChallengeConfig{
			Category: input.Category,
			Score: types.ScoreType{
				Dynamic:   input.Score.Dynamic,
				BaseScore: input.Score.BaseScore,
			},
			Flag: types.FlagType{
				Type:  input.Flag.Type,
				Value: input.Flag.Value,
			},
			Description:  input.Description,
			ExternalLink: input.ExternalLink,
			Singleton:    input.Singleton,
			NodeConfig:   nodes,
		}
		marshalConfig, err := json.Marshal(newChallengeConfig)
		if err != nil {
			return err
		}
		//newChallenge := entity.Challenge{Configuration: string(marshalConfig), State: input.State, Name: input.Name}
		newChallenge := entity.Challenge{
			Configuration: string(marshalConfig),
			State:         "disabled",
			Name:          input.Name,
			FirstBloodId:  nil,
			SecondBloodId: nil,
			ThirdBloodId:  nil,
		}
		result := tx.Create(&newChallenge)
		if result.Error != nil {
			return result.Error
		}

		return nil
	})
	if err != nil {
		log.Println(err)
		return false
	}

	return true
}

// UpdateChallenge we don't allow update challenge name/singleton/[isDynamic]Flag/[isDynamic]Score/status
func UpdateChallenge(input types.ChallengeMutateInput) bool { //TODO: maybe we should checkout if the value need to be update
	warmUp := func() {}
	err := db.Transaction(func(tx *gorm.DB) error {
		var challenge entity.Challenge
		inputChallengeId, err := strconv.ParseUint(input.ChallengeId, 10, 64)
		if err != nil {
			return errors.New("challengeId Parse Error:\n" + err.Error())
		}
		challengeItem := tx.Where(map[string]interface{}{"challenge_id": inputChallengeId}).First(&challenge)
		if challengeItem.Error != nil {
			return errors.New("find challenge by ChallengeId error:\n" + challengeItem.Error.Error())
		}

		if input.Name != challenge.Name {
			return errors.New("can not update challenge name")
		}

		var oldConfig types.ChallengeConfig
		//we don't allow user to change singleton
		err = json.Unmarshal([]byte(challenge.Configuration), &oldConfig)
		if err != nil {
			return err
		}

		//check that user change NodeConfig or not
		var nodes []types.NodeConfig
		//var nodeRefreshFlag = false
		if input.NodeConfig != nil && challenge.State != "enabled" {
			for _, node := range *input.NodeConfig {
				nodes = append(nodes, node.ToNodeConfig())
			}
			//nodeRefreshFlag = true
		} else {
			nodes = oldConfig.NodeConfig
		}

		newChallengeConfig := types.ChallengeConfig{
			Category: input.Category,
			Score: types.ScoreType{
				Dynamic:   oldConfig.Score.Dynamic,
				BaseScore: input.Score.BaseScore,
			},
			Flag: types.FlagType{
				Type:  oldConfig.Flag.Type,
				Value: input.Flag.Value,
			},
			Description:  input.Description,
			ExternalLink: input.ExternalLink,

			Singleton:  oldConfig.Singleton,
			NodeConfig: nodes,
		}
		marshalConfig, err := json.Marshal(newChallengeConfig)
		if err != nil {
			return err
		}
		challenge.Configuration = string(marshalConfig)
		//  if change score or state, warm up all rank
		if oldConfig.Score.BaseScore != input.Score.BaseScore {
			warmUp = func() {
				baseScore, err := strconv.ParseUint(newChallengeConfig.Score.BaseScore, 10, 64)
				if err != nil {
					log.Panicln(err)
				}
				utils.Cache.MutateChallenge(challenge.ChallengeId, challenge.State != "disabled", newChallengeConfig.Score.Dynamic, baseScore)
			}
		}

		if challenge.State == "enabled" && oldConfig.Flag.Value != input.Flag.Value {
			return errors.New("can't change flag for enabled challenge")
		}

		result := tx.Save(&challenge)
		if result.Error != nil {
			log.Println(result.Error)
			return result.Error
		}
		return nil
	})
	if err != nil {
		log.Println(err)
		return false
	}
	warmUp()
	return true
}

func EnableChallengeById(challengeId string) bool {
	id, err := strconv.ParseUint(challengeId, 10, 64)
	if err != nil {
		return false
	}
	challenge, err := FindChallengeById(id)
	if challenge == nil {
		if err != nil {
			log.Println(err)
		}
		log.Println("can't find challenge by challenge id", challengeId)
		return false
	}
	if challenge.State == "enabled" {
		return true
	}
	var config types.ChallengeConfig
	//we don't allow user to change singleton
	err = json.Unmarshal([]byte(challenge.Configuration), &config)
	if err != nil {
		log.Println(err)
		return false
	}
	err = db.Transaction(func(tx *gorm.DB) error {
		challenge.State = "enabled"
		tx.Save(&challenge)

		if config.Singleton {
			replicaId := new(uint64)
			replica := AddReplica(challenge.ChallengeId, tx, nil)
			if replica == nil {
				return errors.New("add replica error")
			}
			*replicaId = replica.ReplicaId
		}

		//set all submits available,TODO: but need some rollback method?
		submits := FindAllSubmitByChallengeId(challenge.ChallengeId)
		for _, submit := range submits {
			submit.Available = true
			tx.Save(&submit)
		}
		return nil
	})
	if err != nil {
		log.Println("challenge enable error: ", err)
		return false
	}
	baseScore, err := strconv.ParseUint(config.Score.BaseScore, 10, 64)
	if err != nil {
		log.Panicln(err)
	}
	utils.Cache.MutateChallenge(challenge.ChallengeId, true, config.Score.Dynamic, baseScore)

	return true
}

func DisableChallengeById(challengeId string) bool {
	id, err := strconv.ParseUint(challengeId, 10, 64)
	if err != nil {
		return false
	}
	challenge, err := FindChallengeById(id)
	if challenge == nil {
		if err != nil {
			log.Println(err)
		}
		log.Println("can't find challenge by challenge id", challengeId)
		return false
	}
	if challenge.State == "disabled" {
		return true
	}
	var config types.ChallengeConfig
	//we don't allow user to change singleton
	err = json.Unmarshal([]byte(challenge.Configuration), &config)
	if err != nil {
		log.Println(err)
		return false
	}
	err = db.Transaction(func(tx *gorm.DB) error {
		challenge.State = "disabled"
		tx.Save(&challenge)

		ok := DeleteReplicaByChallengeId(challenge.ChallengeId, tx, nil)

		if !ok {
			return errors.New("delete replica error")
		}
		//set all submits unavailable,TODO: but need some rollback method?
		submits := FindAllSubmitByChallengeId(challenge.ChallengeId)
		for _, submit := range submits {
			submit.Available = false
			tx.Save(&submit)
		}
		return nil
	})
	if err != nil {
		log.Println("challenge disable error: ", err)
		return false
	}
	baseScore, err := strconv.ParseUint(config.Score.BaseScore, 10, 64)
	if err != nil {
		log.Panicln(err)
	}
	utils.Cache.MutateChallenge(challenge.ChallengeId, false, config.Score.Dynamic, baseScore)

	return true
}

func FindEnabledChallenges() []entity.Challenge {
	var challenges []entity.Challenge
	result := db.Where(map[string]interface{}{"state": "enabled"}).Find(&challenges)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return []entity.Challenge{}
	} else if result.Error != nil {
		log.Println(result.Error)
		return nil
	}
	return challenges
}
func CheckEnabledChallengesByImage(imageName string) bool {
	var challenges []entity.Challenge
	result := db.Where(map[string]interface{}{"state": "enabled"}).Find(&challenges)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return true
	} else if result.Error != nil {
		log.Println(result.Error)
		return false
	}
	for _, challenge := range challenges {
		var config types.ChallengeConfig
		err := json.Unmarshal([]byte(challenge.Configuration), &config)
		if err != nil {
			log.Println(err)
			return false
		}
		for _, nodeConfig := range config.NodeConfig {
			if nodeConfig.Image == imageName {
				return false
			}
		}
	}
	return true
}

func DeleteChallenge(challengeId string) bool {
	//var challenge entity.Challenge
	id, err := strconv.ParseUint(challengeId, 10, 64)
	if err != nil {
		return false
	}
	challenge, err := FindChallengeById(id)
	if challenge == nil {
		if err != nil {
			log.Println(err)
		}
		log.Println("can't find challenge by challenge id", challengeId)
		return false
	}
	var config types.ChallengeConfig

	err = json.Unmarshal([]byte(challenge.Configuration), &config)
	if err != nil {
		log.Println(err)
		return false
	}

	ok := DeleteReplicaByChallengeId(id, nil, func(status bool) {
		var err error
		ok := DeleteSubmitsByChallengeId(id, nil)
		if !ok {
			err = errors.New("delete submits by challengeId error")
		}
		if err != nil {
			log.Println("challenge remove error: ", err)
			return
		}
		//Delete Challenge
		db.Delete(&challenge)
		baseScore, err := strconv.ParseUint(config.Score.BaseScore, 10, 64)
		if err != nil {
			log.Panicln(err)
		}
		utils.Cache.MutateChallenge(id, false, config.Score.Dynamic, baseScore)

	})
	if !ok {
		err = errors.New("delete replica by challengeId error")
	}
	return true

}

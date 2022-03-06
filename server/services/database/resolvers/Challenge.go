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

func FindChallengeByName(name string) (*entity.Challenge, error) {
	var challenge entity.Challenge
	result := db.Where(map[string]interface{}{"Name": name}).First(&challenge)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	} else if result.Error != nil {
		//log.Println(result.Error)
		return nil, result.Error
	}
	return &challenge, nil
}
func AddChallenge(input types.ChallengeMutateInput) bool {
	var challenge *entity.Challenge = nil
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
		for _, node := range input.NodeConfig {
			nodes = append(nodes, node.ToNodeConfig())
		}
		newChallengeConfig := types.ChallengeConfig{
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
			Singleton:    input.Singleton,
			NodeConfig:   nodes,
		}
		marshalConfig, err := json.Marshal(newChallengeConfig)
		if err != nil {
			return err
		}
		newChallenge := entity.Challenge{Configuration: string(marshalConfig), State: input.State, Name: input.Name}
		result := tx.Create(&newChallenge)
		if result.Error != nil {
			return result.Error
		}
		if input.State == "enabled" && input.Singleton {
			replica := AddReplica(newChallenge.ChallengeId, tx)
			if replica == nil {
				return errors.New("create replica failed")
			}
			//ok := EnableReplica(replica.ReplicaId, tx)
			//if !ok {
			//	return errors.New("enabled replica failed")
			//}
			users := FindAllUser()
			if users == nil {
				return errors.New("find all users failed")
			}
			for _, user := range users {
				ok := AddReplicaAlloc(replica.ReplicaId, user.UserId, tx)
				if !ok {
					return errors.New("add replicaAlloc error")
				}
			}
		}
		if input.State == "enabled" {
			challenge = &newChallenge
		}
		return nil
	})
	if err != nil {
		log.Println(err)
		return false
	}
	if challenge != nil {
		originScore, err := strconv.ParseUint(input.Score.BaseScore, 10, 64)
		if err != nil {
			log.Println(err)
			return false
		}
		utils.Cache.AddChallenge(challenge.ChallengeId, originScore)
	}
	return true
}

// UpdateChallenge we don't allow update challenge name/singleton/[isDynamic]Flag/[isDynamic]Score
func UpdateChallenge(input types.ChallengeMutateInput) bool { //TODO: maybe we should checkout if the value need to be update
	needWarmUp := false
	err := db.Transaction(func(tx *gorm.DB) error {
		if input.Name != "" {
			return errors.New("can not update challenge name")
		}

		var challenge entity.Challenge
		inputChallengeId, err := strconv.ParseUint(input.ChallengeId, 10, 64)
		if err != nil {
			return errors.New("challengeId Parse Error:\n" + err.Error())

		}
		challengeItem := tx.Where(map[string]interface{}{"ChallengeId": inputChallengeId}).First(&challenge)
		if challengeItem.Error != nil {
			return errors.New("find challenge by ChallengeId error:\n" + challengeItem.Error.Error())
		}

		var oldConfig types.ChallengeConfig
		//we don't allow user to change singleton
		err = json.Unmarshal([]byte(challenge.Configuration), &oldConfig)
		if err != nil {
			return err
		}

		//check that user change NodeConfig or not
		var nodes []types.NodeConfig
		var nodeRefreshFlag = false
		if input.NodeConfig != nil {
			for _, node := range input.NodeConfig {
				nodes = append(nodes, node.ToNodeConfig())
			}
			nodeRefreshFlag = true
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
				Dynamic: oldConfig.Flag.Dynamic,
				Value:   input.Flag.Value,
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

		//checkResult := tx.Where(map[string]interface{}{"Name": input.Name}).Find(&entity.Challenge{})
		if (challenge.State == "enabled" || input.State == "enabled") && oldConfig.Flag.Value != input.Flag.Value {
			return errors.New("can't change flag for enabled challenge")
		}

		// if change state "disabled", replica delete & set all submits unavailable
		if challenge.State == "enabled" && input.State == "disabled" {
			ok := DeleteReplicaByChallengeId(challenge.ChallengeId, tx)
			if !ok {
				return errors.New("delete replica error")
			}
			//set all submits unavailable,TODO: but need some rollback method?
			submits := FindAllSubmitByChallengeId(challenge.ChallengeId)
			for _, submit := range submits {
				submit.Available = false
				tx.Save(&submit)
			}

		}

		// if change state "enabled" or change node-config, create replicas & alloc to users (only for singleton),set all submits available
		if challenge.State == "disabled" && input.State == "enabled" {
			if oldConfig.Singleton {
				replica := AddReplica(challenge.ChallengeId, tx)
				if replica == nil {
					return errors.New("add replica error")
				}
				users := FindAllUser()
				if users == nil {
					return errors.New("find users error")
				}
				for _, user := range users {
					ok := AddReplicaAlloc(replica.ReplicaId, user.UserId, tx)
					if !ok {
						return errors.New("add replica alloc error")
					}
				}
			}

			//set all submits available,TODO: but need some rollback method?
			submits := FindAllSubmitByChallengeId(challenge.ChallengeId)
			for _, submit := range submits {
				submit.Available = true
				tx.Save(&submit)
			}

		}

		// if change dockerfile, replica re-create
		if nodeRefreshFlag {
			if challenge.State == "enabled" {
				ok := DeleteReplicaByChallengeId(challenge.ChallengeId, tx)
				if !ok {
					return errors.New("delete replica error")
				}
			}
			if input.State == "enabled" {
				if oldConfig.Singleton {
					replica := AddReplica(challenge.ChallengeId, tx)
					if replica == nil {
						return errors.New("add replica error")
					}
					users := FindAllUser()
					if users == nil {
						return errors.New("find users error")
					}
					for _, user := range users {
						ok := AddReplicaAlloc(replica.ReplicaId, user.UserId, tx)
						if !ok {
							return errors.New("add replica alloc error")
						}
					}
				}
			}
		}
		//  if change score or state, warm up all rank
		if (challenge.State != input.State) || (oldConfig.Score.BaseScore != input.Score.BaseScore) {
			needWarmUp = true
		}

		challenge.Configuration = string(marshalConfig)
		challenge.State = input.State
		db.Save(&challenge)
		return nil
	})
	if err != nil {
		log.Println(err)
		return false
	}
	err = utils.Cache.WarmUp()
	if err != nil {
		log.Println("warm up error:\n" + err.Error())
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

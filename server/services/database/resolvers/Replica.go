package resolvers

import (
	"encoding/json"
	"errors"
	"gorm.io/gorm"
	"server/entity"
	"server/services/types"
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

func AddReplica(challengeId uint64) error {
	err := db.Transaction(func(tx *gorm.DB) error {
		var challenge entity.Challenge
		challengeResult := tx.Where(map[string]interface{}{"ChallengeId": challengeId}).First(&challenge)
		if challengeResult.Error != nil {
			return challengeResult.Error
		}
		if challenge.State == "disabled" {
			return errors.New("unable to create replica of disabled challenge")
		}
		var config types.ChallengeConfig
		err := json.Unmarshal([]byte(challenge.Configuration), &config)
		if err != nil {
			return err
		}
		var flag string
		if config.Flag.Dynamic {
			// TODO: generate dynamic flag
		} else {
			flag = config.Flag.Value
		}
		newReplica := entity.Replica{
			ChallengeId: challengeId,
			Status:      "disabled",
			Flag:        flag,
		}
		result := tx.Create(&newReplica)
		if result.Error != nil {
			return result.Error
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func EnableReplica(replicaId uint64, outsideTX *gorm.DB) error {
	if outsideTX == nil {
		outsideTX = db
	}
	err := outsideTX.Transaction(func(tx *gorm.DB) error {
		var replica entity.Replica
		getResult := tx.Where(map[string]interface{}{"ReplicaId": replicaId}).First(&replica)
		if getResult.Error != nil {
			return getResult.Error
		}
		replica.Status = "enabled"
		db.Save(&replica)
		// TODO: K8sPodAlloc
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func DisableReplica(replicaId uint64, outsideTX *gorm.DB) error {
	if outsideTX == nil {
		outsideTX = db
	}
	err := outsideTX.Transaction(func(tx *gorm.DB) error {
		var replica entity.Replica
		getResult := tx.Where(map[string]interface{}{"ReplicaId": replicaId}).First(&replica)
		if getResult.Error != nil {
			return getResult.Error
		}
		replica.Status = "disabled"
		// TODO: K8sPodDestroy
		db.Save(&replica)
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func DeleteReplicaByChallangeId(challengeId uint64, outsideTX *gorm.DB) error {
	if outsideTX == nil {
		outsideTX = db
	}
	err := outsideTX.Transaction(func(tx *gorm.DB) error {
		var replicas []entity.Replica
		getResult := tx.Where(map[string]interface{}{"ChallengeId": challengeId}).Find(&replicas)
		if getResult.Error != nil {
			return getResult.Error
		}
		for _, replica := range replicas {
			// TODO: delete replicaAlloc
			err := DeleteReplicaAllocByReplicaId(replica.ReplicaId, tx)
			if err != nil {
				return err
			}
			err = DisableReplica(replica.ReplicaId, tx)
			if err != nil {
				return err
			}
			db.Delete(&replica)
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}
